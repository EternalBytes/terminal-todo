package main

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"strings"

	"github.com/EternalBytes/todolist"
	"github.com/EternalBytes/todolist/service"
)

func init() {
	var err error
	todolist.Db, err = service.GetDB()
	check(err)
}

func main() {
	add := flag.Bool("add", false, "add a new todo")
	complete := flag.Int("complete", 0, "set todo as completed")
	del := flag.Int("delete", 0, "delete a todo")
	delAll := flag.Bool("delete-all", false, "delete all todos and zero de database")
	list := flag.Bool("list", false, "list all todos")
	flag.Parse()

	todo := new(todolist.Todo)

	switch {
	case *add:
		task, err := parseInput(os.Stdin, flag.Args()...)
		check(err)
		err = todo.Add(task)
		check(err)
	case *complete > 0:
		todo.Complete(*complete)
	case *del > 0 || *delAll:
		todo.Delete(*del)
	case *list:
		todo.Print()
	default:
		check(errors.New("invalid command or value "))
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

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
