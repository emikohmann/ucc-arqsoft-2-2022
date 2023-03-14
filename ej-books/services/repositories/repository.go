package repositories

import (
	"books-api/dtos"
	"books-api/utils/errors"
)

type Repository interface {
	Get(id string) (dtos.BookDTO, errors.ApiError)
	Insert(book dtos.BookDTO) (dtos.BookDTO, errors.ApiError)
	Update(book dtos.BookDTO) (dtos.BookDTO, errors.ApiError)
	Delete(id string) errors.ApiError
}
