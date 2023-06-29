package service

import (
	"github.com/zardan4/todo-app-gin"
	"github.com/zardan4/todo-app-gin/pkg/repository"
)

type TodoListService struct { // створємо сервіс, який залежить від бд, тому додаємо в поле залежностей repository(взаємодія з бд)
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

// в сервісі немає ніякої додаткової логіки, тому просто посилаємо запит на створення в бд
func (s *TodoListService) Create(userId int, list todo.TodoList) (int, error) {
	return s.repo.Create(userId, list)
}

func (s *TodoListService) GetAll(userId int) ([]todo.TodoList, error) {
	return s.repo.GetAll(userId)
}

func (s *TodoListService) GetById(userId, listId int) (todo.TodoList, error) {
	return s.repo.GetById(userId, listId)
}

func (s *TodoListService) Update(userId, listId int, list todo.UpdateTodoListInput) error {
	if err := list.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, listId, list)
}

func (s *TodoListService) Delete(userId, listId int) error {
	return s.repo.Delete(userId, listId)
}
