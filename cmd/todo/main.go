package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"strings"

	"github.com/EternalBytes/todolist/db"
	"github.com/EternalBytes/todolist/service"
)

func main() {
	add := flag.Bool("add", false, "add a new todo")
	complete := flag.Int("complete", 0, "set todo as completed")
	del := flag.Int("delete", 0, "delete a todo")
	delAll := flag.Bool("delete-all", false, "delete all todos and zero de database")
	list := flag.Bool("list", false, "list all todos")
	flag.Parse()

	var ctx = context.Background()
	dbconn, err := service.GetDB()
	if err != nil {
		log.Fatalln(err)
	}

	store := db.NewStore(dbconn)

	switch {
	case *add:
		task, _ := parseInput(os.Stdin, flag.Args()...)
		store.Add(ctx, task)
	case *complete > 0:
		store.Complete(ctx, *complete)
	case *del > 0:
		store.Delete(ctx, *del)
	case *delAll:
		store.DelAll(ctx)
	case *list:
		store.List(ctx)
	default:
		log.Fatalln("invalid command or value")
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
