package todo

import "errors"

type TodoList struct {
	Id          int    `json:"id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type UsersList struct {
	Id     int
	UserId int
	ListId int
}

type TodoItem struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type ListsItem struct {
	Id      int
	ListsId int
	ItemId  int
}

type UpdateTodoListInput struct { // потрібні вказівники, щоб в разі відсутнності поля був nil. це потрібно для валідації
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (l UpdateTodoListInput) Validate() error {
	if l.Title == nil && l.Description == nil {
		return errors.New("update struct has no new values")
	}
	return nil
}

type UpdateTodoItemInput struct { // потрібні вказівники, щоб в разі відсутнності поля був nil. це потрібно для валідації
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

func (r UpdateTodoItemInput) Validate() error {
	if r.Title == nil && r.Description == nil && r.Done == nil {
		return errors.New("update struct has no new values")
	}
	return nil
}
