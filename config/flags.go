package config

import (
	"flag"
	"strings"
	"os"
	"fmt"
	"strconv"
	"go_cli.princeaman.net/internal/models"
)

type Flags struct {
	Add    string
	Delete int
	Edit   string
	List   bool
	Search string
	// DSN string
}

func Commands() *Flags {
	command := Flags{}

	flag.StringVar(&command.Add, "add", "", "Add new task in the list")
	flag.IntVar(&command.Delete, "del", -1, "Delete the following task")
	flag.StringVar(&command.Edit, "edit", "", "Edit the task")
	flag.BoolVar(&command.List, "list", false, "Show all tasks")
	flag.StringVar(&command.Search, "search", "", "Search the content from web")
	// flag.Parse()

	return &command
}

func (cf *Flags) Execute(todos *models.Todos) {
	switch {
	case cf.List:
		todos.Print()
	case cf.Add != "":
		err := todos.Add(cf.Add)
		if err != nil {
			fmt.Println("Error adding task:", err)
		}
	case cf.Edit != "":
		parts := strings.SplitN(cf.Edit, ":", 2)
		if len(parts) != 2 {
			fmt.Println("Error: invalid format for edit. Please use id:new_title")
			os.Exit(1)
		}
		index, err := strconv.Atoi(parts[0])

		if err != nil {
			fmt.Println("Error: invalid index for edit")
			os.Exit(1)
		}
		todos.Edit(index, parts[1])
	case cf.Delete != -1:
		err := todos.Delete(cf.Delete)
		if err != nil {
			fmt.Println("Error deleting task:", err)
		}
	default:
		fmt.Println("Invalid command")
	}
}
