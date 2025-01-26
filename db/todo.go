package db

import (
	"context"
	"time"
)

const addTodo = `
INSERT INTO todos(
	Task, 
	Done, 
	CreatedAt, 
	CompletedAt) VALUES (
		$1,$2,$3,$4)`

type AddTodoParams struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

func (q *Queries) AddTodo(ctx context.Context, args AddTodoParams) (rowsAffected int64, err error) {
	result, err := q.db.ExecContext(ctx, addTodo, args.Task, args.Done, args.CreatedAt, args.CompletedAt)
	if err != nil {
		return
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}

	return
}

const completeTodo = `UPDATE todos SET Done = $1, CompletedAt = $2 WHERE Ind = $3`

func (q *Queries) CompleteTodo(ctx context.Context, index int) (rowsAffected int64, err error) {
	result, err := q.db.ExecContext(ctx, completeTodo, true, time.Now(), index)
	if err != nil {
		return
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}

	return
}

const deleteTodo = `DELETE FROM todos WHERE Ind = $1`

func (q *Queries) DeleteTodo(ctx context.Context, index int) (rowsAffected int64, err error) {
	result, err := q.db.ExecContext(ctx, deleteTodo, index)
	if err != nil {
		return
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}

	return
}

const deleteAll = `DELETE FROM todos;DELETE FROM sqlite_sequence`

func (q *Queries) DeleteAll(ctx context.Context) (err error) {
	_, err = q.db.ExecContext(ctx, deleteAll)
	if err != nil {
		return
	}

	return
}

const listTodos = `SELECT * FROM todos ORDER BY Ind`

func (q *Queries) ListTodos(ctx context.Context) ([]Todo, error) {
	rows, err := q.db.QueryContext(ctx, listTodos)
	if err != nil {
		return nil, err
	}

	todos := []Todo{}

	for rows.Next() {
		var td Todo
		if err := rows.Scan(
			&td.Ind,
			&td.Task,
			&td.Done,
			&td.CreatedAt,
			&td.CompletedAt,
		); err != nil {
			return nil, err
		}
		todos = append(todos, td)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}
