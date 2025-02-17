package repo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type ProductI interface {
	MultipleUpdate(context.Context, []mongo.WriteModel) error
}
