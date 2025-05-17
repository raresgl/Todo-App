package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type CmdFlags struct {
	Add      string
	DueDate  string
	Delete   int
	Edit     int
	Toggle   int
	List     bool
	Clean    bool
	Priority Priority
}

func NewCmdFlags() *CmdFlags {
	cf := CmdFlags{}

	flag.StringVar(&cf.Add, "add", "", "Add a new todo specify title and optionally --due for due date")
	flag.StringVar(&cf.DueDate, "due", "", "Specify due date for new todo (DD.MM.YYYY)")
	flag.Var(&cf.Priority, "prio", "Add a priority to the to do item. 1-High, 2-Medium, 3-Low")
	flag.IntVar(&cf.Edit, "edit", -1, "Edit a to do by index. Can edit the name or the due date")
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
		var dueDate *time.Time
		if cf.DueDate != "" {
			parsedTime, err := time.Parse("02.01.2006", cf.DueDate)
			if err != nil {
				fmt.Println("Error, invalid format for --due. Please use name:due date (DD.LL.YYYY)")
				os.Exit(1)
			}
			dueDate = &parsedTime
		}
		prio := PriorityMedium
		if cf.Priority.isValid() {
			prio = cf.Priority
		}
		todos.add(cf.Add, dueDate, prio)
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
