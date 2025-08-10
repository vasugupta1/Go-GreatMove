package services

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoRepository[T any] interface {
	Create(item T) (T, error)
	FindByID(id string) (T, error)
	FindByFilter(filter interface{}) ([]T, error)
}

type GenericMongoRepository[T any] struct {
	collection *mongo.Collection
}

func ConstructMongoRepository[T any](database *mongo.Database, collectionName string) *GenericMongoRepository[T] {
	return &GenericMongoRepository[T]{
		collection: database.Collection(collectionName),
	}
}

func (r *GenericMongoRepository[T]) Create(item T) (T, error) {
	_, err := r.collection.InsertOne(context.Background(), item)
	if err != nil {
		return item, err
	}
	return item, nil
}

func (r *GenericMongoRepository[T]) FindByID(id string) (T, error) {
	var item T
	err := r.collection.FindOne(context.Background(), map[string]string{"_id": id}).Decode(&item)
	if err != nil {
		return item, err
	}
	return item, nil
}

func (r *GenericMongoRepository[T]) FindByFilter(filter interface{}) ([]T, error) {
	var result []T
	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}

	return result, nil
}
