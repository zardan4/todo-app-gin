package repository

import (
	"database/sql"
	"fmt"

	"github.com/zardan4/todo-app-gin"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db}
}

func (r *AuthPostgres) CreateUser(user todo.User) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) VALUES ($1, $2, $3) RETURNING id", usersTable)
	// request
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)

	if err := row.Scan(&id); err != nil { // записуємо айді створеного користувача
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (todo.User, error) {
	var user todo.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	// request
	row := r.db.QueryRow(query, username, password)

	if err := row.Scan(&user.Id, &user.Name, &user.Username, &user.Password); err != nil { // айді користувача
		return todo.User{}, err
	}

	return user, nil
}
