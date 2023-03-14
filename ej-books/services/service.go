package services

import (
	"books-api/dtos"
	e "books-api/utils/errors"
)

type Service interface {
	Get(id string) (dtos.BookDTO, e.ApiError)
	Insert(book dtos.BookDTO) (dtos.BookDTO, e.ApiError)
}
