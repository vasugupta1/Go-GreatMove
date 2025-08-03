package services

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

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
