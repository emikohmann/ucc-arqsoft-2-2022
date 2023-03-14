package repositories

import (
	"books-api/dtos"
	model "books-api/models"
	e "books-api/utils/errors"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RepositoryMongoDB struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection string
}

func NewMongoDB(host string, port int, collection string) *RepositoryMongoDB {
	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", host, port)))
	if err != nil {
		panic(fmt.Sprintf("Error initializing MongoDB: %v", err))
	}

	names, err := client.ListDatabaseNames(context.TODO(), bson.M{})
	if err != nil {
		panic(fmt.Sprintf("Error initializing MongoDB: %v", err))
	}

	fmt.Println("[MongoDB] Initialized connection")
	fmt.Println(fmt.Sprintf("[MongoDB] Available databases: %s", names))

	return &RepositoryMongoDB{
		Client:     client,
		Database:   client.Database("books"),
		Collection: collection,
	}
}

func (repo *RepositoryMongoDB) Get(id string) (dtos.BookDTO, e.ApiError) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return dtos.BookDTO{}, e.NewBadRequestApiError(fmt.Sprintf("error getting book %s invalid id", id))
	}
	result := repo.Database.Collection(repo.Collection).FindOne(context.TODO(), bson.M{
		"_id": objectID,
	})
	if result.Err() == mongo.ErrNoDocuments {
		return dtos.BookDTO{}, e.NewNotFoundApiError(fmt.Sprintf("book %s not found", id))
	}
	var book model.Book
	if err := result.Decode(&book); err != nil {
		return dtos.BookDTO{}, e.NewInternalServerApiError(fmt.Sprintf("error getting book %s", id), err)
	}
	return dtos.BookDTO{
		Id:   id,
		Name: book.Name,
	}, nil
}

func (repo *RepositoryMongoDB) Insert(book dtos.BookDTO) (dtos.BookDTO, e.ApiError) {
	result, err := repo.Database.Collection(repo.Collection).InsertOne(context.TODO(), model.Book{
		Name: book.Name,
	})
	if err != nil {
		return dtos.BookDTO{}, e.NewInternalServerApiError(fmt.Sprintf("error inserting book %s", book.Id), err)
	}
	book.Id = fmt.Sprintf(result.InsertedID.(primitive.ObjectID).Hex())
	return book, nil
}

func (repo *RepositoryMongoDB) Update(book dtos.BookDTO) (dtos.BookDTO, e.ApiError) {
	_, err := repo.Database.Collection(repo.Collection).UpdateByID(context.TODO(), fmt.Sprintf("%v", book.Id), model.Book{
		Name: book.Name,
	})
	if err != nil {
		return dtos.BookDTO{}, e.NewInternalServerApiError(fmt.Sprintf("error inserting book %s", book.Id), err)
	}
	return book, nil
}

func (repo *RepositoryMongoDB) Delete(id string) e.ApiError {
	_, err := repo.Database.Collection(repo.Collection).DeleteOne(context.TODO(), bson.M{"_id": fmt.Sprintf("%s", id)})
	if err != nil {
		return e.NewInternalServerApiError(fmt.Sprintf("error deleting book %s", id), err)
	}
	return nil
}
