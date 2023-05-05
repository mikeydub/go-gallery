package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.24

import (
	"context"

	"github.com/mikeydub/go-gallery/graphql/generated"
	"github.com/mikeydub/go-gallery/graphql/model"
	"github.com/mikeydub/go-gallery/service/persist"
)

// FindFeedEventByDbid is the resolver for the findFeedEventByDbid field.
func (r *entityResolver) FindFeedEventByDbid(ctx context.Context, dbid persist.DBID) (*model.FeedEvent, error) {
	return resolveFeedEventByEventID(ctx, dbid)
}

// Entity returns generated.EntityResolver implementation.
func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
