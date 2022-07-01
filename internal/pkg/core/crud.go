package core

import "context"

type CrudRepository[T any] interface {
	Save(ctx context.Context, entity T) (T, error)
	SaveAll(ctx context.Context, entities []T) ([]T, error)

	Count() int
	ExistsById(id int) bool

	FindAll(ctx context.Context) ([]T, error)
	FindAllById(ctx context.Context, ids []int) ([]T, error)
	FindById(ctx context.Context, id int) (T, error)

	Delete(entity T) error
	DeleteAll(entities ...T) error
	DeleteById(id int) error
	DeleteAllById(ids []int) error
}