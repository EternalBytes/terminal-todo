package todolist

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/EternalBytes/todolist/service"
	"github.com/alexeyco/simpletable"
)

type item struct {
	Ind         int
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type Todos []item

func (t *Todos) Add(task string) error {
	db, err := service.GetDB()
	defer close(db)
	if err != nil {
		return err
	}
	result, err := db.Exec("INSERT INTO todos(Task, Done, CreatedAt, CompletedAt) VALUES(?,?,?,?)",
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

func (t *Todos) Complete(ind int) {
	db, err := service.GetDB()
	defer close(db)
	check(err)
	rowsAf, err := db.Exec("UPDATE todos SET Done=?, CompletedAt=? WHERE Ind=?", true, time.Now(), ind)
	check(err)
	rows, err := rowsAf.RowsAffected()
	check(err)
	if rows > 0 {
		fmt.Println(ColorGreen + "Task Completed" + ColorDefault)
	}
}

func (t *Todos) Delete(ind int) {
	db, err := service.GetDB()
	defer close(db)
	check(err)
	rowsAf, err := db.Exec("DELETE FROM todos WHERE Ind=?", ind)
	check(err)
	rows, err := rowsAf.RowsAffected()
	check(err)
	if rows > 0 {
		fmt.Println(ColorRed + "Task Deleted" + ColorDefault)
	}
}

func (t *Todos) Print() {
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

	db, err := service.GetDB()
	defer close(db)
	check(err)

	rows, err := db.Query("SELECT * FROM todos")
	defer func() {
		err := rows.Close()
		check(err)
	}()
	check(err)
	var countUndone int
	var it item
	for rows.Next() {
		rows.Scan(&it.Ind, &it.Task, &it.Done, &it.CreatedAt, &it.CompletedAt)

		task := blue(it.Task)
		if it.Done {
			task = green(fmt.Sprintf("\u2705 %s", it.Task))
		}
		*cells = append(*cells, []*simpletable.Cell{
			{Text: fmt.Sprintf("%d", it.Ind)},
			{Text: task},
			{Text: fmt.Sprintf("%t", it.Done)},
			{Text: it.CreatedAt.Format(time.RFC822)},
			{Text: it.CompletedAt.Format(time.RFC822)},
		})
		/// COUNT UNDONE
		if !it.Done {
			countUndone++
		}
	}

	table.Body = &simpletable.Body{Cells: *cells}

	undTxt := red(fmt.Sprintf("\u26A0 %s", "You have "+fmt.Sprint(countUndone)+" task to do"))
	if countUndone == 0 {
		undTxt = blue(fmt.Sprintf("🎉 %s", "You've done all tasks"))
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

func red(s string) string {
	return fmt.Sprintf("%s%s%s", ColorRed, s, ColorDefault)
}

func green(s string) string {
	return fmt.Sprintf("%s%s%s", ColorGreen, s, ColorDefault)
}

func blue(s string) string {
	return fmt.Sprintf("%s%s%s", ColorBlue, s, ColorDefault)
}

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
