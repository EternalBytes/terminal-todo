package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/alexeyco/simpletable"
)

const (
	ColorDefault = "\x1b[39m"
	ColorGreen   = "\x1b[32m"
	ColorBlue    = "\x1b[94m"
	ColorRed     = "\x1b[91m"
)

type Store struct {
	*Queries
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
	}
}

func (s *Store) Add(ctx context.Context, task string) {
	args := AddTodoParams{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	r, err := s.AddTodo(ctx, args)
	check(err)

	if r > 0 {
		fmt.Println(ColorGreen + "Task added" + ColorDefault)
	}
}

func (s *Store) Complete(ctx context.Context, index int) {
	r, err := s.CompleteTodo(ctx, index)
	check(err)

	if r > 0 {
		fmt.Println(ColorGreen + "Task Completed" + ColorDefault)
		return
	}
	fmt.Println(ColorRed + "No Task to Complete" + ColorDefault)
}

func (s *Store) Delete(ctx context.Context, index int) {
	r, err := s.DeleteTodo(ctx, index)
	check(err)

	if r > 0 {
		fmt.Println(ColorRed + "Task Deleted" + ColorDefault)
		return
	}
	fmt.Println(ColorRed + "No Task to Delete" + ColorDefault)
}

func (s *Store) DelAll(ctx context.Context) {
	err := s.DeleteAll(ctx)
	check(err)
	fmt.Println(ColorRed + "All Tasks Was Deleted" + ColorDefault)
}

func (s *Store) List(ctx context.Context) {
	rowschan := make(chan []Todo)
	// creates a goroutine to get data in parallel
	go func(rs chan []Todo) {
		rows, err := s.ListTodos(ctx)
		check(err)
		rs <- rows
	}(rowschan)

	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Done?"},
			{Align: simpletable.AlignRight, Text: "CreatedAt"},
			{Align: simpletable.AlignRight, Text: "CompletedAt"},
		},
	}

	cells := new([][]*simpletable.Cell)

	var countUndone, countDone int
	for _, v := range <-rowschan {
		task := fmt.Sprint(ColorBlue + v.Task + ColorDefault)
		if v.Done {
			task = fmt.Sprint(ColorGreen + "\u2705 " + v.Task + ColorDefault)
			// COUNT DONE
			countDone++
		} else {
			// COUNT UNDONE
			countUndone++
		}

		*cells = append(*cells, []*simpletable.Cell{
			{Text: fmt.Sprintf("%d", v.Ind)},
			{Text: task},
			{Text: fmt.Sprintf("%t", v.Done)},
			{Text: v.CreatedAt.Format(time.RFC822)},
			{Text: v.CompletedAt.Format(time.RFC822)},
		})
	}

	table.Body = &simpletable.Body{Cells: *cells}

	undTxt := fmt.Sprint(ColorRed + "\u26A0 You have " + fmt.Sprint(countUndone) + " tasks to do" + ColorDefault)
	if countUndone == 0 {
		if countDone == 0 {
			undTxt = fmt.Sprint(ColorRed + "No task was found." + ColorDefault)
			goto footer // goto is very useful for a programmer as me @EternalBytes
		}
		undTxt = fmt.Sprint(ColorBlue + "🎉 You've done all tasks" + ColorDefault)
	}

footer:
	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 5, Text: undTxt},
	}}

	table.SetStyle(simpletable.StyleUnicode)

	table.Println()
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
