package inner

import (
	"context"
	"fmt"

	"github.com/SaidakbarPardaboyev/get-all-from-ucode/pkg"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MultipleUpdate struct {
	Collection string
	Config     *pkg.InnerConfig
}

func NewMultipleUpdate(collection string, config *pkg.InnerConfig) *MultipleUpdate {
	return &MultipleUpdate{
		Collection: collection,
		Config:     config,
	}
}

func (m *MultipleUpdate) Exec(ctx context.Context, updates []mongo.WriteModel) error {
	collection := m.Config.MongoDb.Collection(m.Collection)

	bulkOptions := options.BulkWrite().SetOrdered(false) // Allow parallel execution
	_, err := collection.BulkWrite(ctx, updates, bulkOptions)
	if err != nil {
		return fmt.Errorf("bulk update failed at batch %v", err)
	}

	return nil
}
