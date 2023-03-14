package repositories

import (
	"books-api/dtos"
	e "books-api/utils/errors"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	json "github.com/json-iterator/go"
)

type RepositoryMemcached struct {
	Client *memcache.Client
}

func NewMemcached(host string, port int) *RepositoryMemcached {
	client := memcache.New(fmt.Sprintf("%s:%d", host, port))
	fmt.Println("[Memcached] Initialized connection")
	return &RepositoryMemcached{
		Client: client,
	}
}

func (repo *RepositoryMemcached) Get(id string) (dtos.BookDTO, e.ApiError) {
	item, err := repo.Client.Get(id)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return dtos.BookDTO{}, e.NewNotFoundApiError(fmt.Sprintf("book %s not found", id))
		}
		return dtos.BookDTO{}, e.NewInternalServerApiError(fmt.Sprintf("error getting book %s", id), err)
	}

	var bookDTO dtos.BookDTO
	if err := json.Unmarshal(item.Value, &bookDTO); err != nil {
		return dtos.BookDTO{}, e.NewInternalServerApiError(fmt.Sprintf("error getting book %s", id), err)
	}

	return bookDTO, nil
}

func (repo *RepositoryMemcached) Insert(book dtos.BookDTO) (dtos.BookDTO, e.ApiError) {
	bytes, err := json.Marshal(book)
	if err != nil {
		return dtos.BookDTO{}, e.NewBadRequestApiError(err.Error())
	}

	if err := repo.Client.Set(&memcache.Item{
		Key:   book.Id,
		Value: bytes,
	}); err != nil {
		return dtos.BookDTO{}, e.NewInternalServerApiError(fmt.Sprintf("error inserting book %s", book.Id), err)
	}

	return book, nil
}

func (repo *RepositoryMemcached) Update(book dtos.BookDTO) (dtos.BookDTO, e.ApiError) {
	bytes, err := json.Marshal(book)
	if err != nil {
		return dtos.BookDTO{}, e.NewBadRequestApiError(fmt.Sprintf("invalid book %s: %v", book.Id, err))
	}

	if err := repo.Client.Set(&memcache.Item{
		Key:   book.Id,
		Value: bytes,
	}); err != nil {
		return dtos.BookDTO{}, e.NewInternalServerApiError(fmt.Sprintf("error updating book %s", book.Id), err)
	}

	return book, nil
}

func (repo *RepositoryMemcached) Delete(id string) e.ApiError {
	err := repo.Client.Delete(id)
	if err != nil {
		return e.NewInternalServerApiError(fmt.Sprintf("error deleting book %s", id), err)
	}
	return nil
}
