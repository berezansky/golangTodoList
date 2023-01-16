package repository

import (
	"todolist/consts"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(consts.User) (int, error)
	GetUser(username, password string) (consts.User, error)
}

type TodoList interface {
	Create(userId int, list consts.TodoList) (int, error)
	GetAll(userId int) ([]consts.TodoList, error)
	GetById(userId int, id int) (consts.TodoList, error)
	Delete(userId, listId int) error
	Update(userId int, id int, input consts.UpdateListInput) error
}

type TodoItem interface {
	Create(listId int, input consts.TodoItem) (int, error)
	GetAll(userId, listId int) ([]consts.TodoItem, error)
	GetById(userId, itemId int) (consts.TodoItem, error)
	Update(userId, itemId int, input consts.UpdateItemInput) error
	Delete(userId, itemId int) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
