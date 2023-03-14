package services

import (
	"books-api/dtos"
	"books-api/services/repositories"
	e "books-api/utils/errors"
	"fmt"
	"net/http"
)

type ServiceImpl struct {
	localCache repositories.Repository
	distCache  repositories.Repository
	db         repositories.Repository
}

func NewServiceImpl(
	localCache repositories.Repository,
	distCache repositories.Repository,
	db repositories.Repository,
) *ServiceImpl {
	return &ServiceImpl{
		localCache: localCache,
		distCache:  distCache,
		db:         db,
	}
}

func (serv *ServiceImpl) Get(id string) (dtos.BookDTO, e.ApiError) {
	var book dtos.BookDTO
	var apiErr e.ApiError
	var source string

	// try to find it in localCache
	book, apiErr = serv.localCache.Get(id)
	if apiErr != nil {
		if apiErr.Status() != http.StatusNotFound {
			return dtos.BookDTO{}, apiErr
		}
		// try to find it in distCache
		book, apiErr = serv.distCache.Get(id)
		if apiErr != nil {
			if apiErr.Status() != http.StatusNotFound {
				return dtos.BookDTO{}, apiErr
			}
			// try to find it in db
			book, apiErr = serv.db.Get(id)
			if apiErr != nil {
				if apiErr.Status() != http.StatusNotFound {
					return dtos.BookDTO{}, apiErr
				} else {
					fmt.Println(fmt.Sprintf("Not found book %s in any datasource", id))
					apiErr = e.NewNotFoundApiError(fmt.Sprintf("book %s not found", id))
					return dtos.BookDTO{}, apiErr
				}
			} else {
				source = "db"
				defer func() {
					if _, apiErr := serv.distCache.Insert(book); apiErr != nil {
						fmt.Println(fmt.Sprintf("Error trying to save book in distCache %v", apiErr))
					}
					if _, apiErr := serv.localCache.Insert(book); apiErr != nil {
						fmt.Println(fmt.Sprintf("Error trying to save book in localCache %v", apiErr))
					}
				}()
			}
		} else {
			source = "distCache"
			defer func() {
				if _, apiErr := serv.localCache.Insert(book); apiErr != nil {
					fmt.Println(fmt.Sprintf("Error trying to save book in localCache %v", apiErr))
				}
			}()
		}
	} else {
		source = "localCache"
	}

	fmt.Println(fmt.Sprintf("Obtained book from %s!", source))
	return book, nil
}

func (serv *ServiceImpl) Insert(book dtos.BookDTO) (dtos.BookDTO, e.ApiError) {
	result, apiErr := serv.db.Insert(book)
	if apiErr != nil {
		fmt.Println(fmt.Sprintf("Error inserting book in db: %v", apiErr))
		return dtos.BookDTO{}, apiErr
	}
	fmt.Println(fmt.Sprintf("Inserted book in db: %v", result))

	_, apiErr = serv.distCache.Insert(result)
	if apiErr != nil {
		fmt.Println(fmt.Sprintf("Error inserting book in distCache: %v", apiErr))
		return result, nil
	}
	fmt.Println(fmt.Sprintf("Inserted book in distCache: %v", result))

	_, apiErr = serv.localCache.Insert(result)
	if apiErr != nil {
		fmt.Println(fmt.Sprintf("Error inserting book in localCache: %v", apiErr))
		return result, nil
	}
	fmt.Println(fmt.Sprintf("Inserted book in localCache: %v", result))

	return result, nil
}
