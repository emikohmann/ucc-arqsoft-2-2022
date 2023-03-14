package repositories

import (
	"books-api/dtos"
	e "books-api/utils/errors"
	"fmt"
	"github.com/karlseguin/ccache/v2"
	"time"
)

type RepositoryCCache struct {
	Client     *ccache.Cache
	DefaultTTL time.Duration
}

func NewCCache(maxSize int64, itemsToPrune uint32, defaultTTL time.Duration) *RepositoryCCache {
	client := ccache.New(ccache.Configure().MaxSize(maxSize).ItemsToPrune(itemsToPrune))
	fmt.Println("[CCache] Initialized")
	return &RepositoryCCache{
		Client:     client,
		DefaultTTL: defaultTTL,
	}
}

func (repo *RepositoryCCache) Get(id string) (dtos.BookDTO, e.ApiError) {
	item := repo.Client.Get(id)
	if item == nil {
		return dtos.BookDTO{}, e.NewNotFoundApiError(fmt.Sprintf("book %s not found", id))
	}
	if item.Expired() {
		return dtos.BookDTO{}, e.NewNotFoundApiError(fmt.Sprintf("book %s not found", id))
	}
	return item.Value().(dtos.BookDTO), nil
}

func (repo *RepositoryCCache) Insert(book dtos.BookDTO) (dtos.BookDTO, e.ApiError) {
	repo.Client.Set(book.Id, book, repo.DefaultTTL)
	return book, nil
}

func (repo *RepositoryCCache) Update(book dtos.BookDTO) (dtos.BookDTO, e.ApiError) {
	repo.Client.Set(book.Id, book, repo.DefaultTTL)
	return book, nil
}

func (repo *RepositoryCCache) Delete(id string) e.ApiError {
	repo.Client.Delete(id)
	return nil
}
