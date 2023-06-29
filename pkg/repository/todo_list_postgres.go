package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/zardan4/todo-app-gin"
)

type TodoListPostgres struct {
	db *sql.DB
}

func NewTodoListPostgres(db *sql.DB) *TodoListPostgres {
	return &TodoListPostgres{db}
}

func (l *TodoListPostgres) Create(userId int, list todo.TodoList) (int, error) {
	tx, err := l.db.Begin() // транзакція, бо потрібно вставити в таблиці users_lists і todo_lists
	if err != nil {
		return 0, err
	}

	// частина транзакції з todo_lists
	query := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row := tx.QueryRow(query, list.Title, list.Description)

	var listId int
	err = row.Scan(&listId) // записуємо айді створеного списку
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// частина транзакції з users_lists
	query = fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2) RETURNING id", usersListsTable)
	_, err = tx.Exec(query, userId, listId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return listId, tx.Commit()
}

// отримуємо всі лісти
func (l *TodoListPostgres) GetAll(userId int) ([]todo.TodoList, error) {
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1",
		todoListsTable, usersListsTable) // хз, не шарю за sql, але має брати всі списки, які створені користувачем userId
	rows, err := l.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []todo.TodoList
	for rows.Next() {
		var list todo.TodoList
		if err := rows.Scan(&list.Id, &list.Title, &list.Description); err != nil {
			return res, err
		}
		res = append(res, list)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, err
}

func (l *TodoListPostgres) GetById(userId, listId int) (todo.TodoList, error) {
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2",
		todoListsTable, usersListsTable) // хз, не шарю за sql, але має брати всі списки, які створені користувачем userId і є N-ним лістом
	row := l.db.QueryRow(query, userId, listId)

	var res todo.TodoList
	if err := row.Scan(&res.Id, &res.Title, &res.Description); err != nil {
		return res, err
	}

	return res, nil
}

func (l *TodoListPostgres) Update(userId, listId int, list todo.UpdateTodoListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if list.Title != nil { // перевіряємо на наявність такого аргументу
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId)) // додаємо аргумент в строфу, яка буде вставлятися в query
		args = append(args, *list.Title)                               // всі значення до setValues
		argId++                                                        // збільшуємо аргумент, щоб потім могли створтювати нові плейсххолдери
	}

	if list.Description != nil { // перевіряємо на наявність такого аргументу
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId)) // додаємо аргумент в строфу, яка буде вставлятися в query
		args = append(args, *list.Description)                               // всі значення до setValues
		argId++                                                              // збільшуємо аргумент, щоб потім могли створтювати нові плейсххолдери
	}

	// title=$1
	// description=$1
	// title=$1, description=$2
	setQuery := strings.Join(setValues, ", ") // остаточно формуємо строфу, яку будемо вставляти в query

	// апдейтимо таблицю лістів, задаємо в неї всі значення(setQuery), а потім перевіряємо на все, що було і в минулих запитах, АЛЕ тут ми це робимо за допомогою інкрементування argId
	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d",
		todoListsTable, setQuery, usersListsTable, argId, argId+1)

	args = append(args, listId, userId) // додаємо аргументи id

	_, err := l.db.Exec(query, args...)

	return err
}

func (l *TodoListPostgres) Delete(userId, listId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id=ul.list_id AND ul.user_id=$1 AND ul.list_id=$2", todoListsTable, usersListsTable)
	_, err := l.db.Exec(query, userId, listId)

	return err
}
