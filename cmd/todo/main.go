package main

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"os"
	"strings"

	"github.com/EternalBytes/todolist"
)

const (
	todoFile = ".todos.json" // name of the data file
)

func main() {
	add := flag.Bool("add", false, "add a new todo")
	complete := flag.Int("complete", 0, "set todo as completed")
	del := flag.Int("delete", 0, "delete a todo")
	list := flag.Bool("list", false, "list all todos")
	flag.Parse()

	todos := new(todolist.Todos)

	if err := todos.Load(todoFile); err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	switch {
	case *add:
		task, err := parseInput(os.Stdin, flag.Args()...)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
		todos.Add(task)
		if err := todos.Store(todoFile); err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
	case *complete > 0:
		if err := todos.Complete(*complete); err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
		if err := todos.Store(todoFile); err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
	case *del > 0:
		if err := todos.Delete(*del); err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
		if err := todos.Store(todoFile); err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
	case *list:
		todos.Print()
	default:
		os.Stderr.WriteString("invalid command or value")
		os.Exit(1)
	}
}

func parseInput(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}

	text := scanner.Text()

	if len(text) == 0 {
		return "", errors.New("empty todo isn't allowed")
	}

	return text, nil
}
