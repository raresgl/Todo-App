A simple command-line Todo application written in Go.

## Features

- Add todos with optional due dates
- List todos (with sorting by due date soon)
- Mark todos as complete
- Edit todo titles and due dates
- Clean/remove todos

## Usage

```
# Add a todo
./todo -add=\"Buy groceries\"
./todo -add=\"Pay bills:15.05.2025\"

# List todos
./todo -list

# Complete a todo
./todo -complete 1

# Edit a todo
./todo -edit 2
```

## Installation

```
go install github.com/yourusername/todo-cli@latest
```
