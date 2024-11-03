package todolist

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"slices"
	"time"

	"github.com/alexeyco/simpletable"
)

type item struct {
	Task        string    `json:"task"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"createdat"`
	CompletedAt time.Time `json:"completedat"`
}

type Todos []item

func (t *Todos) Add(task string) {
	todo := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*t = append(*t, todo)
}

func (t *Todos) Complete(ind int) error {
	td := *t
	if ind <= 0 || ind > len(*t) {
		return errors.New("invalid index")
	}

	td[ind-1].CompletedAt = time.Now()
	td[ind-1].Done = true

	return nil
}

func (t *Todos) Delete(ind int) error {
	if ind <= 0 || ind > len(*t) {
		return errors.New("invalid index")
	}

	*t = slices.Delete(*t, ind-1, ind)

	return nil
}

func (t *Todos) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if len(data) == 0 {
		return err
	}

	err = json.Unmarshal(data, t)
	if err != nil {
		return err
	}
	return nil
}

func (t *Todos) Store(filename string) error {
	bytes, err := json.Marshal(t)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, bytes, fs.ModePerm)
	if err != nil {
		return err
	}
	return err
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

	for ind, value := range *t {
		ind++
		task := blue(value.Task)
		if value.Done {
			task = green(fmt.Sprintf("\u2705 %s", value.Task))
		}
		*cells = append(*cells, []*simpletable.Cell{
			{Text: fmt.Sprintf("%d", ind)},
			{Text: task},
			{Text: fmt.Sprintf("%t", value.Done)},
			{Text: value.CreatedAt.Format(time.RFC822)},
			{Text: value.CompletedAt.Format(time.RFC822)},
		})
	}

	table.Body = &simpletable.Body{Cells: *cells}

	undone := t.countUndone()
	undTxt := red(fmt.Sprintf("\u26A0 %s", "You have "+fmt.Sprint(undone)+" task to do"))
	if undone == 0 {
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

func (t *Todos) countUndone() int {
	var done int
	for _, v := range *t {
		if !v.Done {
			done++
		}
	}
	return done
}
