package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aquasecurity/table"
)

type Todo struct {
	Title       string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time
	DueDate     *time.Time
}

type Todos []Todo

func (todos *Todos) add(title string, dueDate *time.Time) {
	todo := Todo{
		Title:       title,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
		DueDate:     dueDate,
	}
	*todos = append(*todos, todo)
}

func (todos *Todos) validateIndex(index int) error {
	if index < 0 || index >= len(*todos) {
		err := errors.New("invalid index")
		fmt.Printf("Error %v", err)
		return err
	}
	return nil
}

func (todos *Todos) delete(index int) error {
	t := *todos
	err := t.validateIndex(index)
	if err != nil {
		return err
	}
	*todos = append(t[:index], t[index+1:]...)

	return nil
}

func (todos *Todos) toggle(index int) error {
	t := *todos

	if err := t.validateIndex(index); err != nil {
		return err
	}

	isCompleted := t[index].Completed

	if !isCompleted {
		completionTime := time.Now()
		t[index].CompletedAt = &completionTime
	}

	t[index].Completed = !isCompleted

	return nil
}

func (todos *Todos) edit(index int) error {
	t := *todos

	if err := t.validateIndex(index); err != nil {
		return err
	}

	fmt.Println("What would you like to edit?")
	fmt.Println("1. Title")
	fmt.Println("2. Due Date")

	var choice string
	fmt.Print("Enter your choice (1-2): ")
	fmt.Scanln(&choice)

	switch choice {
	case "1":
		fmt.Print("Enter new title: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			t[index].Title = scanner.Text()
		}
		fmt.Println("Title updated successfully")
	case "2":
		fmt.Print("Enter new due date (DD.MM.YYYY) or leave empty to remove: ")
		scanner := bufio.NewScanner(os.Stdin)
		var dateStr string
		if scanner.Scan() {
			dateStr = scanner.Text()
		}

		if dateStr == "" {
			t[index].DueDate = nil
			fmt.Println("Due date removed")
		} else {
			parsedTime, err := time.Parse("02.01.2006", dateStr)
			if err != nil {
				fmt.Println("Error: Invalid date format. Please use DD.MM.YYYY format")
				os.Exit(1)
			}
			t[index].DueDate = &parsedTime
			fmt.Println("Due date updated successfully")
		}
	default:
		fmt.Println("Invalid choice. No changes made.")
		os.Exit(1)
	}

	return nil
}

func (todos *Todos) print() {
	table := table.New(os.Stdout)
	table.SetRowLines(false)
	table.SetHeaders("#", "Title", "Completed", "CreatedAt", "CompletedAt", "DueDate")
	for index, t := range *todos {
		completed := "❌"
		completedAt := ""

		if t.Completed {
			completed = "✅"
			if t.CompletedAt != nil {
				completedAt = t.CompletedAt.Format(time.RFC1123)
			}
		}

		dueDate := ""
		if t.DueDate != nil {
			dueDate = t.DueDate.Format(time.RFC1123)
		}

		table.AddRow(strconv.Itoa(index), t.Title, completed, t.CreatedAt.Format(time.RFC1123), completedAt, dueDate)
	}
	table.Render()
}

func (todos *Todos) clean(deleteAll bool) {

	if deleteAll {
		*todos = Todos{}
	} else {
		activeTodos := Todos{}
		for _, t := range *todos {
			if !t.Completed {
				activeTodos = append(activeTodos, t)
			}
		}
		*todos = activeTodos
	}
}
