package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aquasecurity/table"
)

type Priority int

const (
	PriorityHigh   Priority = 1
	PriorityMedium Priority = 2
	PriorityLow    Priority = 3
)

type Todo struct {
	Title       string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time
	DueDate     *time.Time
	Priority    Priority
}

type Todos []Todo

type Value interface {
	String() string
	Set(string) error
}

func (p Priority) isValid() bool {
	return p == PriorityHigh || p == PriorityMedium || p == PriorityLow
}

func (p Priority) String() string {
	switch p {
	case PriorityHigh:
		return "High"
	case PriorityMedium:
		return "Medium"
	case PriorityLow:
		return "Low"
	default:
		return "Unknown"
	}
}

func (p *Priority) Set(s string) error {
	val, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("priority must be a number")
	}
	tmp := Priority(val)
	if !tmp.isValid() {
		return fmt.Errorf("invalid priority: %d (valid: 1=High, 2=Medium, 3=Low)", val)
	}
	*p = tmp
	return nil
}

func (todos *Todos) add(title string, dueDate *time.Time, priority Priority) {
	todo := Todo{
		Title:       title,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
		DueDate:     dueDate,
		Priority:    priority,
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
	fmt.Println("3. Priority")

	var choice string
	fmt.Print("Enter your choice (1-3): ")
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
	case "3":
		fmt.Print("Enter the new desired priority: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			input := strings.TrimSpace(scanner.Text())
			parsed, err := strconv.Atoi(input)
			p := Priority(parsed)
			if err != nil || !p.isValid() {
				fmt.Println("Error: Invalid priority. Please enter 1 (High), 2 (Medium), or 3 (Low):")
				os.Exit(1)
			}
			t[index].Priority = p
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
	table.SetHeaders("#", "Title", "Completed", "CreatedAt", "CompletedAt", "DueDate", "Priority")
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

		table.AddRow(strconv.Itoa(index), t.Title, completed, t.CreatedAt.Format(time.RFC1123), completedAt, dueDate, t.Priority.String())
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
