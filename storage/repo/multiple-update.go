package repo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type MultipleUpdateI interface {
	Exec(context.Context, []mongo.WriteModel) error
}
