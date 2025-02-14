package get_all

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ItemsI interface {
	GetAll() *GetAllI
}

func (o *object) Items(collection string) ItemsI {
	return &APIItem{
		collection: collection,
		config: &innerConfig{
			MongoDb:    o.mongoDb,
			PostgresDb: o.postgresDb,
			DB_TYPE:    o.config.DB_TYPE,
		},
	}
}

func (i *APIItem) GetAll() *GetAllI {
	return &GetAllI{
		collection: i.collection,
		config:     i.config,
	}
}

func (g *GetAllI) Filter(filter map[string]interface{}) *GetAllI {
	g.filter = filter
	return g
}

func (g *GetAllI) Sort(sort map[string]interface{}) *GetAllI {
	g.sort = sort
	return g
}

func (g *GetAllI) Limit(limit int64) *GetAllI {
	g.limit = limit
	return g
}

func (g *GetAllI) Skip(skip int64) *GetAllI {
	g.skip = skip
	return g
}

func (g *GetAllI) Pipeline(pipeline []map[string]any) *GetAllI {
	if g.config.DB_TYPE == "mongo" {
		g.pipeline = pipeline
	}
	return g
}

// Count returns the number of records matching the filter
func (g *GetAllI) Count() (int64, error) {
	if g.collection == "" || g.config == nil {
		return 0, fmt.Errorf("collection and config must be set")
	}

	if g.config.DB_TYPE == "mongo" {
		return g.countMongo()
	} else if g.config.DB_TYPE == "postgres" {
		return g.countPostgres()
	}

	return 0, fmt.Errorf("unsupported DB_TYPE")
}

func (g *GetAllI) countMongo() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Hour)
	defer cancel()

	collection := g.config.MongoDb.Collection(g.collection)

	filter := bson.M{}
	if g.filter != nil {
		filter = g.filter
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %v", err)
	}

	return count, nil
}

func (g *GetAllI) countPostgres() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Hour)
	defer cancel()

	// Base query
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", g.collection)

	var values []interface{}
	if g.filter != nil {
		conditions, vals := buildPostgresConditions(g.filter)
		query += fmt.Sprintf(" WHERE %s", conditions)
		values = vals
	}

	// Execute the query
	var count int64
	err := g.config.PostgresDb.QueryRow(ctx, query, values...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count rows: %v", err)
	}

	return count, nil
}

// Exec executes the query and returns the results
func (g *GetAllI) Exec() ([]map[string]interface{}, error) {
	if g.collection == "" || g.config == nil {
		return nil, fmt.Errorf("collection and config must be set")
	}

	if g.config.DB_TYPE == "mongo" {
		return g.execMongo()
	} else if g.config.DB_TYPE == "postgres" {
		return g.execPostgres()
	}

	return nil, fmt.Errorf("unsupported DB_TYPE")
}

// func (g *GetAllI) execMongo() ([]map[string]interface{}, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	collection := g.config.MongoDb.Collection(g.collection)
// 	opts := options.Find()

// 	if g.sort != nil {
// 		opts.SetSort(g.sort)
// 	}
// 	if g.limit > 0 {
// 		opts.SetLimit(g.limit)
// 	}
// 	if g.skip > 0 {
// 		opts.SetSkip(g.skip)
// 	}

// 	filter := bson.M{}
// 	if g.filter != nil {
// 		filter = g.filter
// 	}

// 	cursor, err := collection.Find(ctx, filter, opts)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to execute query: %v", err)
// 	}
// 	defer cursor.Close(ctx)

// 	var results []map[string]interface{}
// 	for cursor.Next(ctx) {
// 		var doc map[string]interface{}
// 		if err := cursor.Decode(&doc); err != nil {
// 			return nil, fmt.Errorf("failed to decode document: %v", err)
// 		}
// 		results = append(results, doc)
// 	}

// 	if err := cursor.Err(); err != nil {
// 		return nil, fmt.Errorf("cursor error: %v", err)
// 	}

// 	return results, nil
// }

func (g *GetAllI) execMongo() ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Hour)
	defer cancel()

	collection := g.config.MongoDb.Collection(g.collection)
	var results []map[string]any

	matchedPipeline := mongo.Pipeline{}
	modifiedPipeline := mongo.Pipeline{}
	if g.pipeline != nil { // If pipeline exists, modify it dynamically

		// Convert pipeline (map[string]any) to mongo.Pipeline
		for _, stage := range g.pipeline {

			for key, value := range stage {

				if key == "$match" {

					matchStage := bson.D{}
					matchMap, ok := value.(map[string]any)

					if ok {

						for field, val := range matchMap {

							switch v := val.(type) {
							case []string: // Convert to {"field": {"$in": [...]}}
								matchStage = append(matchStage, bson.E{field, bson.D{{"$in", v}}})
							default:
								matchStage = append(matchStage, bson.E{field, v})
							}
						}

					}

					matchedPipeline = append(matchedPipeline, bson.D{{"$match", matchStage}})

				} else {

					modifiedPipeline = append(modifiedPipeline, bson.D{{key, value}})

				}
			}
		}
	}

	// If a filter exists, prepend a $match stage
	if g.filter != nil {
		matchedPipeline = append(matchedPipeline, bson.D{{"$match", g.filter}})
	}

	// If sort exists, append a $sort stage
	if g.sort != nil {
		modifiedPipeline = append(modifiedPipeline, bson.D{{"$sort", g.sort}})
	}

	// If skip exists, append a $skip stage
	if g.skip > 0 {
		modifiedPipeline = append(modifiedPipeline, bson.D{{"$skip", g.skip}})
	}

	// If limit exists, append a $limit stage
	if g.limit > 0 {
		modifiedPipeline = append(modifiedPipeline, bson.D{{"$limit", g.limit}})
	}

	modifiedPipeline = append(matchedPipeline, modifiedPipeline...)

	cursor, err := collection.Aggregate(ctx, modifiedPipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to execute aggregation: %v", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var doc map[string]any
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode document: %v", err)
		}
		results = append(results, doc)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}
	return results, nil
}

func (g *GetAllI) execPostgres() ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Hour)
	defer cancel()

	query := fmt.Sprintf("SELECT * FROM %s", g.collection)

	var values []interface{}
	if g.filter != nil {
		conditions, vals := buildPostgresConditions(g.filter)
		query += fmt.Sprintf(" WHERE %s", conditions)
		values = vals
	}

	if g.sort != nil {
		query += " ORDER BY " + buildPostgresSort(g.sort)
	}
	if g.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", g.limit)
	}
	if g.skip > 0 {
		query += fmt.Sprintf(" OFFSET %d", g.skip)
	}

	rows, err := g.config.PostgresDb.Query(ctx, query, values...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	columns := rows.FieldDescriptions()

	var results []map[string]interface{}
	for rows.Next() {
		row := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(columns))
		for i := range row {
			scanArgs[i] = &row[i]
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		rowMap := make(map[string]interface{})
		for i, col := range columns {
			rowMap[string(col.Name)] = row[i]
		}
		results = append(results, rowMap)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows error: %v", rows.Err())
	}

	return results, nil
}

func buildPostgresConditions(filter map[string]interface{}) (string, []interface{}) {
	conditions := ""
	values := []interface{}{}
	index := 1
	for key, value := range filter {
		if conditions != "" {
			conditions += " AND "
		}
		conditions += fmt.Sprintf("%s = $%d", key, index)
		values = append(values, value)
		index++
	}
	return conditions, values
}

func buildPostgresSort(sort map[string]interface{}) string {
	sortStr := ""
	for key, value := range sort {
		if sortStr != "" {
			sortStr += ", "
		}
		direction := "ASC"
		if value == -1 {
			direction = "DESC"
		}
		sortStr += fmt.Sprintf("%s %s", key, direction)
	}
	return sortStr
}

// package get_all

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// type ItemsI interface {
// 	GetAll() *GetAllI
// }

// func (o *object) Items(collection string) ItemsI {
// 	return &APIItem{
// 		collection: collection,
// 		config: &innerConfig{
// 			MongoDb:    o.mongoDb,
// 			PostgresDb: o.postgresDb,
// 			DB_TYPE:    o.config.DB_TYPE,
// 		},
// 	}
// }

// func (i *APIItem) GetAll() *GetAllI {
// 	return &GetAllI{
// 		collection: i.collection,
// 		config:     i.config,
// 	}
// }

// func (g *GetAllI) Filter(filter map[string]interface{}) *GetAllI {
// 	g.filter = filter
// 	return g
// }

// func (g *GetAllI) Sort(sort map[string]interface{}) *GetAllI {
// 	g.sort = sort
// 	return g
// }

// func (g *GetAllI) Limit(limit int64) *GetAllI {
// 	g.limit = limit
// 	return g
// }

// func (g *GetAllI) Skip(skip int64) *GetAllI {
// 	g.skip = skip
// 	return g
// }

// // Count returns the number of documents matching the filter
// func (g *GetAllI) Count() (int64, error) {
// 	if g.collection == "" || g.config == nil {
// 		return 0, fmt.Errorf("collection and config must be set")
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	collection := g.config.MongoDb.Collection(g.collection)

// 	filter := bson.M{}
// 	if g.filter != nil {
// 		filter = g.filter
// 	}

// 	count, err := collection.CountDocuments(ctx, filter)
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to count documents: %v", err)
// 	}

// 	return count, nil
// }

// // Exec executes the query and returns the results
// func (g *GetAllI) Exec() ([]map[string]interface{}, error) {
// 	// Ensure the collection and config are set
// 	if g.collection == "" || g.config == nil {
// 		return nil, fmt.Errorf("collection and config must be set")
// 	}

// 	// Create a context with timeout
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// Get the MongoDB collection
// 	collection := g.config.MongoDb.Collection(g.collection)

// 	// Build options for the query
// 	opts := options.Find()
// 	if g.sort != nil {
// 		opts.SetSort(g.sort)
// 	}
// 	if g.limit > 0 {
// 		opts.SetLimit(g.limit)
// 	}
// 	if g.skip > 0 {
// 		opts.SetSkip(g.skip)
// 	}

// 	// Convert filter map to BSON
// 	filter := bson.M{}
// 	if g.filter != nil {
// 		filter = g.filter
// 	}

// 	// Execute the query
// 	cursor, err := collection.Find(ctx, filter, opts)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to execute query: %v", err)
// 	}
// 	defer cursor.Close(ctx)

// 	// Parse the results
// 	var results []map[string]interface{}
// 	for cursor.Next(ctx) {
// 		var doc map[string]interface{}
// 		if err := cursor.Decode(&doc); err != nil {
// 			return nil, fmt.Errorf("failed to decode document: %v", err)
// 		}

// 		results = append(results, doc)
// 	}

// 	// Return any cursor error
// 	if err := cursor.Err(); err != nil {
// 		return nil, fmt.Errorf("cursor error: %v", err)
// 	}

// 	return results, nil
// }
