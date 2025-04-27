package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type CmdFlags struct {
	Add    string
	Delete int
	Edit   int
	Toggle int
	List   bool
	Clean  bool
}

func NewCmdFlags() *CmdFlags {
	cf := CmdFlags{}

	flag.StringVar(&cf.Add, "add", "", "Add a new todo specify title")
	flag.IntVar(&cf.Edit, "edit", -1, "Edit a to do by index")
	flag.IntVar(&cf.Delete, "delete", -1, "Specify a todo by index to delete")
	flag.IntVar(&cf.Toggle, "toggle", -1, "Specify a todo by index to toggle")
	flag.BoolVar(&cf.List, "list", false, "List all todos")
	flag.BoolVar(&cf.Clean, "clean", false, "Clear table but choose if you want to keep active todos (active:true)")
	flag.Parse()

	return &cf
}

func (cf *CmdFlags) Execute(todos *Todos) {
	switch {
	case cf.List:
		todos.print()
	case cf.Add != "":
		parts := strings.SplitN(cf.Add, ":", 2)
		if len(parts) == 2 {
			dateStr := strings.TrimSpace(parts[1])
			var dueDate *time.Time
			if dateStr != "" {
				parsedTime, err := time.Parse("02.01.2006", dateStr)
				if err != nil {
					fmt.Println("Error, invalid format for add. Please use name:due date (DD.LL.YYYY)")
					os.Exit(1)
				}
				dueDate = &parsedTime
			}
			todos.add(parts[0], dueDate)
		} else if len(parts) == 1 {
			todos.add(cf.Add, nil)
		} else {
			fmt.Println("Error, invalid format for add. Please use name:due date (DD.LL.YYYY)")
			os.Exit(1)
		}
	case cf.Edit != -1:
		todos.edit(cf.Edit)
	case cf.Toggle != -1:
		todos.toggle(cf.Toggle)
	case cf.Delete != -1:
		todos.delete(cf.Delete)
	case cf.Clean:
		todos.clean(cf.Clean)
	default:
		fmt.Println("Invalid command")

	}
}
