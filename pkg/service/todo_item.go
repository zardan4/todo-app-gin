package service

import (
	"github.com/zardan4/todo-app-gin"
	"github.com/zardan4/todo-app-gin/pkg/repository"
)

type TodoItemService struct { // створємо сервіс, який залежить від бд, тому додаємо в поле залежностей repository(взаємодія з бд)
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo}
}

func (r *TodoItemService) Create(userId, listId int, item todo.TodoItem) (int, error) {
	// якщо такого списку не існує/список не належить користувачу
	_, err := r.listRepo.GetById(userId, listId)
	if err != nil {
		return 0, err
	}

	return r.repo.Create(listId, item)
}

func (r *TodoItemService) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	// якщо такого списку не існує/список не належить користувачу
	_, err := r.listRepo.GetById(userId, listId)
	if err != nil {
		return nil, err
	}

	return r.repo.GetAll(userId, listId)
}

func (r *TodoItemService) GetById(userId, itemId int) (todo.TodoItem, error) {
	return r.repo.GetById(userId, itemId)
}

func (r *TodoItemService) Update(userId, itemId int, item todo.UpdateTodoItemInput) error {
	if err := item.Validate(); err != nil {
		return err
	}
	return r.repo.Update(userId, itemId, item)
}

func (r *TodoItemService) Delete(userId, itemId int) error {
	return r.repo.Delete(userId, itemId)
}
