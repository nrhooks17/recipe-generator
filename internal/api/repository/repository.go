package repository

type Repository[T any] interface {
	Insert() error
	Get(id int) (T, error)
	GetAll() ([]T, error)
	Update(item T) error
	Delete(item T) error
}
