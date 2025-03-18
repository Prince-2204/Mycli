package models

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aquasecurity/table"
)

// Todo represents a single task in the todo list
type Todo struct {
	ID          int
	Title       string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}

// Todos interacts with the database for task operations
type Todos struct {
	DB *DataModel
}

// Add inserts a new task into the database
func (todos *Todos) Add(title string) error {
	stmt := `INSERT INTO todo (title, completed) VALUES (?, false)`
	_, err := todos.DB.DB.Exec(stmt, title)
	if err != nil {
		fmt.Println("Error adding task:", err)
		return err
	}
	fmt.Println("Task added successfully:", title)
	return nil
}

// ValidateIndex checks if a task with the given ID exists
func (todos *Todos) ValidateIndex(id int) error {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM todo WHERE id = ?)`
	err := todos.DB.DB.QueryRow(query, id).Scan(&exists)
	if err != nil {
		fmt.Println("Error validating task ID:", err)
		return err
	}
	if !exists {
		fmt.Println("Invalid ID:", id)
		return fmt.Errorf("task ID %d does not exist", id)
	}
	return nil
}

// Delete removes a task from the database
func (todos *Todos) Delete(id int) error {
	err := todos.ValidateIndex(id)
	if err != nil {
		return err
	}

	stmt := `DELETE FROM todo WHERE id = ?`
	_, err = todos.DB.DB.Exec(stmt, id)
	if err != nil {
		fmt.Println("Error deleting task:", err)
		return err
	}

	fmt.Println("Task deleted successfully:", id)
	return nil
}

// Edit updates the title of a task
func (todos *Todos) Edit(id int, title string) error {
	err := todos.ValidateIndex(id)
	if err != nil {
		return err
	}

	stmt := `UPDATE todo SET title = ? WHERE id = ?`
	_, err = todos.DB.DB.Exec(stmt, title, id)
	if err != nil {
		fmt.Println("Error updating task:", err)
		return err
	}

	fmt.Println("Task updated successfully:", id)
	return nil
}

// MarkCompleted sets a task as completed
func (todos *Todos) MarkCompleted(id int) error {
	err := todos.ValidateIndex(id)
	if err != nil {
		return err
	}

	stmt := `UPDATE todo SET completed = true, completedat = ? WHERE id = ?`
	_, err = todos.DB.DB.Exec(stmt, time.Now(), id)
	if err != nil {
		fmt.Println("Error marking task as completed:", err)
		return err
	}

	fmt.Println("Task marked as completed:", id)
	return nil
}

// Print displays all tasks in a formatted table
func (todos *Todos) Print() {
	query := `SELECT id, title, completed, createdat, completedat FROM todo ORDER BY createdat DESC`
	rows, err := todos.DB.DB.Query(query)
	if err != nil {
		fmt.Println("Error fetching tasks:", err)
		return
	}
	defer rows.Close()

	table := table.New(os.Stdout)
	table.SetRowLines(false)
	table.SetHeaders("#", "Title", "Completed", "Created At", "Completed At")

	for rows.Next() {
		var t Todo
		var completedAt sql.NullTime

		err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt, &completedAt)
		if err != nil {
			fmt.Println("Error reading row:", err)
			continue
		}

		if completedAt.Valid {
			t.CompletedAt = &completedAt.Time
		}

		completed := "❌"
		completedAtStr := ""

		if t.Completed {
			completed = "✅"
			if t.CompletedAt != nil {
				completedAtStr = t.CompletedAt.Format(time.RFC1123)
			}
		}

		table.AddRow(strconv.Itoa(t.ID), t.Title, completed, t.CreatedAt.Format(time.RFC1123), completedAtStr)
	}

	table.Render()
}
