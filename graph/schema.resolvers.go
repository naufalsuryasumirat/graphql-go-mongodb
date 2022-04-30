package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/naufalsuryasumirat/graphql-go-mongodb/database"
	"github.com/naufalsuryasumirat/graphql-go-mongodb/graph/generated"
	"github.com/naufalsuryasumirat/graphql-go-mongodb/graph/model"
)

var db = database.Connect()

func (r *mutationResolver) CreateAuthor(ctx context.Context, input *model.AuthorInput) (*model.Author, error) {
	return db.AddAuthor(input), nil
}

func (r *mutationResolver) CreateBook(ctx context.Context, input *model.BookInput, idAuthor *string) (*model.Book, error) {
	return db.AddBook(input, idAuthor), nil
}

func (r *queryResolver) Book(ctx context.Context, id string) (*model.Book, error) {
	return db.FindBookByID(id), nil
}

func (r *queryResolver) Books(ctx context.Context) ([]*model.Book, error) {
	return db.AllBooks(), nil
}

func (r *queryResolver) BooksByAuthor(ctx context.Context, idAuthor string) ([]*model.Book, error) {
	return db.AllBooksByAuthor(idAuthor), nil
}

func (r *queryResolver) Author(ctx context.Context, id string) (*model.Author, error) {
	return db.FindAuthorByID(id), nil
}

func (r *queryResolver) Authors(ctx context.Context) ([]*model.Author, error) {
	return db.AllAuthors(), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
