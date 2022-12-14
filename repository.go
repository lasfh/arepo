package arepo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type abstractRepo[T any] struct {
	collection *mongo.Collection
}

func NewAbstractRepository[T any](collection *mongo.Collection) AbstractRepository[T] {
	return &abstractRepo[T]{
		collection: collection,
	}
}

func (a *abstractRepo[T]) GetByID(ctx context.Context, id ID, opts ...*options.FindOneOptions) (*T, error) {
	return FindOneByID[T](ctx, a.collection, id, opts...)
}

func (a *abstractRepo[T]) FindOne(ctx context.Context, filter any, opts ...*options.FindOneOptions) (*T, error) {
	return FindOne[T](ctx, a.collection, filter, opts...)
}

func (a *abstractRepo[T]) FindOneAndUpdate(ctx context.Context, filter, update any, opts ...*options.FindOneAndUpdateOptions) (*T, error) {
	return FindOneAndUpdate[T](ctx, a.collection, filter, update, opts...)
}

func (a *abstractRepo[T]) Find(ctx context.Context, filter any, opts ...*options.FindOptions) ([]*T, error) {
	return Find[T](ctx, a.collection, filter, opts...)
}

func (a *abstractRepo[T]) InsertOne(ctx context.Context, document *T, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return a.collection.InsertOne(ctx, document, opts...)
}

func (a *abstractRepo[T]) InsertMany(ctx context.Context, documents []*T, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	docs := make([]any, 0, len(documents))

	for _, doc := range documents {
		docs = append(docs, doc)
	}

	return a.collection.InsertMany(ctx, docs, opts...)
}

func (a *abstractRepo[T]) UpdateOneByID(ctx context.Context, id ID, update any) error {
	return UpdateOneByID(ctx, a.collection, id, update)
}

func (a *abstractRepo[T]) UpdateOne(ctx context.Context, filter, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return UpdateOne(ctx, a.collection, filter, update, opts...)
}

func (a *abstractRepo[T]) DeleteOneByID(ctx context.Context, id ID, opts ...*options.DeleteOptions) error {
	return DeleteOneByID(ctx, a.collection, id, opts...)
}

func (a *abstractRepo[T]) DeleteOne(ctx context.Context, filter any, opts ...*options.DeleteOptions) error {
	return DeleteOne(ctx, a.collection, filter, opts...)
}

func (a *abstractRepo[T]) Select(fields ...string) SelectAndOmitFields[T] {
	setProjection := bson.D{}

	for i := range fields {
		setProjection = append(setProjection, bson.E{Key: fields[i], Value: 1})
	}

	return &selectAndOmitFields[T]{
		collection:    a.collection,
		setProjection: setProjection,
	}
}

func (a *abstractRepo[T]) Omit(fields ...string) SelectAndOmitFields[T] {
	setProjection := bson.D{}

	for i := range fields {
		setProjection = append(setProjection, bson.E{Key: fields[i], Value: 0})
	}

	return &selectAndOmitFields[T]{
		collection:    a.collection,
		setProjection: setProjection,
	}
}

type selectAndOmitFields[T any] struct {
	collection    *mongo.Collection
	setProjection bson.D
}

func (s *selectAndOmitFields[T]) Select(fields ...string) SelectAndOmitFields[T] {
	for i := range fields {
		s.setProjection = append(s.setProjection, bson.E{Key: fields[i], Value: 1})
	}

	return s
}

func (s *selectAndOmitFields[T]) Omit(fields ...string) SelectAndOmitFields[T] {
	for i := range fields {
		s.setProjection = append(s.setProjection, bson.E{Key: fields[i], Value: 0})
	}

	return s
}

func (s *selectAndOmitFields[T]) GetByID(ctx context.Context, id ID) (*T, error) {
	if len(s.setProjection) == 0 {
		return nil, ErrNotSelectOmitFields
	}

	return FindOneByID[T](ctx, s.collection, id, options.FindOne().SetProjection(s.setProjection))
}

func (s *selectAndOmitFields[T]) FindOne(ctx context.Context, filter any) (*T, error) {
	return FindOne[T](ctx, s.collection, filter, options.FindOne().SetProjection(s.setProjection))
}

func (s *selectAndOmitFields[T]) FindOneAndUpdate(ctx context.Context, filter, update any) (*T, error) {
	return FindOneAndUpdate[T](ctx, s.collection, filter, update, options.FindOneAndUpdate().SetProjection(s.setProjection))
}

func (s *selectAndOmitFields[T]) Find(ctx context.Context, filter any) ([]*T, error) {
	return Find[T](ctx, s.collection, filter, options.Find().SetProjection(s.setProjection))
}
