package service

import (
	"todolist/consts"
	"todolist/pkg/repository"
)

type Authorization interface {
	CreateUser(consts.User) (int, error)
	GenerateToken(username string, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	Create(userId int, list consts.TodoList) (int, error)
	GetAll(userId int) ([]consts.TodoList, error)
	GetById(userId int, id int) (consts.TodoList, error)
	Delete(userId int, id int) error
	Update(userId int, id int, input consts.UpdateListInput) error
}

type TodoItem interface {
	Create(userId, listId int, input consts.TodoItem) (int, error)
	GetAll(userId, listId int) ([]consts.TodoItem, error)
	GetById(userId, itemId int) (consts.TodoItem, error)
	Update(userId, itemId int, input consts.UpdateItemInput) error
	Delete(userId, itemId int) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		TodoList:      NewTodoListService(repo.TodoList),
		TodoItem:      NewTodoItemService(repo.TodoItem, repo.TodoList),
	}
}
