package todolist

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/alexeyco/simpletable"
)

var Db *sql.DB

type Todo struct {
	Ind         int
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

func (t *Todo) Add(task string) error {
	defer close(Db)
	result, err := Db.Exec("INSERT INTO todos(Task, Done, CreatedAt, CompletedAt) VALUES(?,?,?,?)",
		task,
		false,
		time.Now(),
		time.Time{})
	if err != nil {
		return err
	}

	rowsAf, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAf > 0 {
		fmt.Println(ColorGreen + "Task added" + ColorDefault)
	}
	return nil
}

func (t *Todo) Complete(ind int) {
	defer close(Db)
	rowsAf, err := Db.Exec("UPDATE todos SET Done=?, CompletedAt=? WHERE Ind=?", true, time.Now(), ind)
	check(err)
	rows, err := rowsAf.RowsAffected()
	check(err)
	if rows > 0 {
		fmt.Println(ColorGreen + "Task Completed" + ColorDefault)
	}
}

func (t *Todo) Delete(index int) {
	defer close(Db)
	var query string = "DELETE FROM todos WHERE Ind=?"
	if index == 0 {
		query = "DELETE FROM todos;DELETE FROM sqlite_sequence"
	}

	rowsAf, err := Db.Exec(query, index)
	check(err)
	rows, err := rowsAf.RowsAffected()
	check(err)
	if rows > 0 {
		if index == 0 {
			fmt.Println(ColorRed + "All Tasks Was Deleted" + ColorDefault)
			return
		}
		fmt.Println(ColorRed + "Task Deleted" + ColorDefault)
	}
}

func (t *Todo) Print() {
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

	defer close(Db)
	rows, err := Db.Query("SELECT * FROM todos")
	defer func() {
		err := rows.Close()
		check(err)
	}()
	check(err)
	var countUndone int
	item := *t
	for rows.Next() {
		rows.Scan(&item.Ind, &item.Task, &item.Done, &item.CreatedAt, &item.CompletedAt)

		task := fmt.Sprint(ColorBlue + item.Task + ColorDefault)
		if t.Done {
			task = fmt.Sprint(ColorGreen + "\u2705 " + item.Task + ColorDefault)
		} else {
			// COUNT UNDONE
			countUndone++
		}
		*cells = append(*cells, []*simpletable.Cell{
			{Text: fmt.Sprintf("%d", item.Ind)},
			{Text: task},
			{Text: fmt.Sprintf("%t", item.Done)},
			{Text: item.CreatedAt.Format(time.RFC822)},
			{Text: item.CompletedAt.Format(time.RFC822)},
		})
	}

	table.Body = &simpletable.Body{Cells: *cells}

	undTxt := fmt.Sprint(ColorRed + "\u26A0 You have " + fmt.Sprint(countUndone) + " task to do" + ColorDefault)
	if countUndone == 0 {
		undTxt = fmt.Sprint(ColorBlue + "🎉 You've done all tasks" + ColorDefault)
	}
	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 5, Text: undTxt},
	}}

	table.SetStyle(simpletable.StyleUnicode)

	table.Println()
}

const (
	ColorDefault = "\x1b[39m"
	ColorGreen   = "\x1b[32m"
	ColorBlue    = "\x1b[94m"
	ColorRed     = "\x1b[91m"
)

func close(db *sql.DB) {
	err := db.Close()
	if err != nil {
		panic(err)
	}
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
