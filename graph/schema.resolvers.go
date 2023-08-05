package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.34

import (
	"context"
	"go-graphql-mongodb-api/database"
	"go-graphql-mongodb-api/graph/model"
)

// CreateKolo is the resolver for the createKolo field.
func (r *mutationResolver) CreateKolo(ctx context.Context, input model.NewKolo) (*model.Kolo, error) {
	return db.InsertKolo(input), nil
}

// UpdateKolo is the resolver for the updateKolo field.
func (r *mutationResolver) UpdateKolo(ctx context.Context, input model.UpdateKolo) (*model.Kolo, error) {
	return db.UpdateKolo(input), nil
}

// DeleteKolo is the resolver for the deleteKolo field.
func (r *mutationResolver) DeleteKolo(ctx context.Context, input string) (string, error) {
	return db.DeleteKolo(input), nil
}

// CreatePostajalisce is the resolver for the createPostajalisce field.
func (r *mutationResolver) CreatePostajalisce(ctx context.Context, input model.NewPostajalisce) (*model.Postajalisce, error) {
	return db.InsertPostajalisce(input), nil
}

// UpdatePostajalisce is the resolver for the updatePostajalisce field.
func (r *mutationResolver) UpdatePostajalisce(ctx context.Context, input model.UpdatePostajalisce) (*model.Postajalisce, error) {
	return db.UpdatePostajalisce(input), nil
}

// DeletePostajalisce is the resolver for the deletePostajalisce field.
func (r *mutationResolver) DeletePostajalisce(ctx context.Context, input string) (string, error) {
	return db.DeletePostajalisce(input), nil
}

// IzposojaKolesa is the resolver for the izposojaKolesa field.
func (r *mutationResolver) IzposojaKolesa(ctx context.Context, input model.IzposojaKolesa) (*model.Izposoja, error) {
	return db.BorrowKolo(input), nil
}

// VraciloKolesa is the resolver for the VraciloKolesa field.
func (r *mutationResolver) VraciloKolesa(ctx context.Context, input model.VraciloKolesa) (*model.Izposoja, error) {
	return db.ReturnKolo(input), nil
}

// DeleteIzposoja is the resolver for the deleteIzposoja field.
func (r *mutationResolver) DeleteIzposoja(ctx context.Context, input string) (string, error) {
	return db.DeleteIzposoja(input), nil
}

// Kolo is the resolver for the kolo field.
func (r *queryResolver) Kolo(ctx context.Context, id string) (*model.Kolo, error) {
	return db.FindKolo(id), nil
}

// Kolesa is the resolver for the kolesa field.
func (r *queryResolver) Kolesa(ctx context.Context) ([]*model.Kolo, error) {
	return db.FindAllKolo(), nil
}

// Postajalisce is the resolver for the postajalisce field.
func (r *queryResolver) Postajalisce(ctx context.Context, id string) (*model.Postajalisce, error) {
	return db.FindPostajalisce(id), nil
}

// Postajalisca is the resolver for the postajalisca field.
func (r *queryResolver) Postajalisca(ctx context.Context) ([]*model.Postajalisce, error) {
	return db.FindAllPostajalisce(), nil
}

// Izposoja is the resolver for the izposoja field.
func (r *queryResolver) Izposoja(ctx context.Context, id string) (*model.Izposoja, error) {
	return db.FindIzposoja(id), nil
}

// Izposoje is the resolver for the izposoje field.
func (r *queryResolver) Izposoje(ctx context.Context) ([]*model.Izposoja, error) {
	return db.FindAllIzposoja(), nil
}

// NearestPostajalisce is the resolver for the nearestPostajalisce field.
func (r *queryResolver) NearestPostajalisce(ctx context.Context, latitude float64, longitude float64, stPostaj int) ([]*model.Postajalisce, error) {
	return db.FindNearestPostajalisce(latitude, longitude, stPostaj), nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
var db = database.Connect("mongodb+srv://admin:admin@cluster0.nikbntq.mongodb.net/?retryWrites=true&w=majority")
