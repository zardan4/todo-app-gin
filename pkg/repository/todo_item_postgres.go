package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/zardan4/todo-app-gin"
)

type TodoItemPostgres struct {
	db *sql.DB
}

func NewTodoItemPostgres(db *sql.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db}
}

func (l *TodoItemPostgres) Create(listId int, item todo.TodoItem) (int, error) {
	tx, err := l.db.Begin() // транзакція, бо потрібно вставити в таблиці users_lists і todo_lists
	if err != nil {
		return 0, err
	}

	// частина транзакції з todo_items
	query := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoItemsTable)
	row := tx.QueryRow(query, item.Title, item.Description)

	var itemId int
	err = row.Scan(&itemId) // записуємо айді створеного списку
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// частина транзакції з users_lists
	query = fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2) RETURNING id", listsItemsTable)
	_, err = tx.Exec(query, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (l *TodoItemPostgres) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li on li.item_id = ti.id 
	INNER JOIN %s ul on ul.list_id = li.list_id WHERE li.list_id = $1 AND ul.user_id = $2`, todoItemsTable, listsItemsTable, usersListsTable)
	rows, err := l.db.Query(query, listId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []todo.TodoItem
	for rows.Next() {
		var item todo.TodoItem
		if err := rows.Scan(&item.Id, &item.Title, &item.Description, &item.Done); err != nil {
			return res, err
		}
		res = append(res, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, err
}

func (l *TodoItemPostgres) GetById(userId, itemId int) (todo.TodoItem, error) {
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li on li.item_id = ti.id	INNER JOIN %s ul on ul.list_id = li.list_id WHERE ti.id = $1 AND ul.user_id = $2`, todoItemsTable, listsItemsTable, usersListsTable)

	row := l.db.QueryRow(query, itemId, userId)

	var res todo.TodoItem
	if err := row.Scan(&res.Id, &res.Title, &res.Description, &res.Done); err != nil {
		return res, err
	}

	return res, nil
}

func (l *TodoItemPostgres) Update(userId, itemId int, item todo.UpdateTodoItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if item.Title != nil { // перевіряємо на наявність такого аргументу
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId)) // додаємо аргумент в строфу, яка буде вставлятися в query
		args = append(args, *item.Title)                               // всі значення до setValues
		argId++                                                        // збільшуємо аргумент, щоб потім могли створтювати нові плейсххолдери
	}

	if item.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *item.Description)
		argId++
	}

	if item.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *item.Done)
		argId++
	}

	// те саме, що і в todo_list_postgres.go
	setQuery := strings.Join(setValues, ", ") // остаточно формуємо строфу, яку будемо вставляти в query

	query := fmt.Sprintf(`UPDATE %s ti SET %s FROM %s li, %s ul WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d`,
		todoItemsTable, setQuery, listsItemsTable, usersListsTable, argId, argId+1)

	args = append(args, userId, itemId) // додаємо аргументи id

	_, err := l.db.Exec(query, args...)

	return err
}

func (l *TodoItemPostgres) Delete(userId, itemId int) error {
	query := fmt.Sprintf(`DELETE FROM %s ti USING %s li, %s ul WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND ti.id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)

	_, err := l.db.Exec(query, userId, itemId)
	return err
}
