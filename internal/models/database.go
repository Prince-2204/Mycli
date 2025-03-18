package models

import (
	"database/sql"
	"log"
	"time"
	"fmt"
)

// DataModel holds the database connection
type DataModel struct {
	DB *sql.DB
}

// AddTask inserts a new task into the todo table
func (m *DataModel) AddTask(title string) error {
	query := "INSERT INTO todo (title, completed) VALUES (?, false)"
	_, err := m.DB.Exec(query, title)
	if err != nil {
		log.Println("Error inserting task:", err)
		return err
	}
	log.Println("Task added successfully:", title)
	return nil
}

// ListTasks retrieves all tasks from the todo table
func (m *DataModel) ListTasks() ([]string, error) {
	query := "SELECT id, title, completed FROM todo"
	rows, err := m.DB.Query(query)
	if err != nil {
		log.Println("Error retrieving tasks:", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []string
	for rows.Next() {
		var id int
		var title string
		var completed bool
		if err := rows.Scan(&id, &title, &completed); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		status := "Pending"
		if completed {
			status = "Completed"
		}
		tasks = append(tasks, fmt.Sprintf("[%d] %s - %s", id, title, status))
	}
	return tasks, nil
}

// MarkTaskCompleted updates a task as completed
func (m *DataModel) MarkTaskCompleted(id int) error {
	query := "UPDATE todo SET completed = true, completedat = ? WHERE id = ?"
	_, err := m.DB.Exec(query, time.Now(), id)
	if err != nil {
		log.Println("Error marking task as completed:", err)
		return err
	}
	log.Println("Task marked as completed:", id)
	return nil
}

// DeleteTask removes a task from the database
func (m *DataModel) DeleteTask(id int) error {
	query := "DELETE FROM todo WHERE id = ?"
	_, err := m.DB.Exec(query, id)
	if err != nil {
		log.Println("Error deleting task:", err)
		return err
	}
	log.Println("Task deleted successfully:", id)
	return nil
}
