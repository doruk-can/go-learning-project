package mongo

import (
	"context"
	"fmt"
	"time"

	"foover/internal/config"
	"foover/internal/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Store defines behaviors of the MongoDB store
type Store interface {
	Close() error
	CreateSession(ctx context.Context) (string, error)
	SaveVote(ctx context.Context, vote models.Vote) error
	GetVotesBySessionID(ctx context.Context, sessionID string) ([]models.Vote, error)
	GetAggregatedProductScores(ctx context.Context) ([]models.ProductScore, error)
	SaveProducts(ctx context.Context, products []models.Product) error
	IsValidProductID(ctx context.Context, productID string) (bool, error)
	SessionExists(ctx context.Context, sessionID string) (bool, error)
}

// store represents the MongoDB store
type store struct {
	client           *mongo.Client
	db               *mongo.Database
	sessionTimeout   time.Duration
	voteTimeout      time.Duration
	aggregateTimeout time.Duration
}

// NewStore creates and returns a new MongoDB store
func NewStore(cfg config.Mongo) (Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.URI)
	clientOptions.SetMinPoolSize(cfg.MinPoolSize)
	clientOptions.SetMaxPoolSize(cfg.MaxPoolSize)
	clientOptions.SetConnectTimeout(cfg.ConnectTimeout)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("connecting failed, %s", err.Error())
	}

	// Ping the database to ensure a connection is established
	pingCtx, pingCancel := context.WithTimeout(context.Background(), cfg.PingTimeout)
	defer pingCancel()

	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		return nil, err
	}

	db := client.Database(cfg.Database)

	s := &store{
		client:           client,
		db:               db,
		sessionTimeout:   cfg.ReadTimeout,
		voteTimeout:      cfg.WriteTimeout,
		aggregateTimeout: cfg.ReadTimeout,
	}

	if err := s.ensureIndexes(); err != nil {
		return nil, err
	}

	return s, nil
}

// Close disconnects the MongoDB client
func (s *store) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.client.Disconnect(ctx)
}

// CreateSession generates and stores a new unique session ID
func (s *store) CreateSession(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.sessionTimeout)
	defer cancel()

	sessionID := uuid.New().String()
	session := models.Session{
		SessionID: sessionID,
		CreatedAt: time.Now(),
	}

	_, err := s.db.Collection("sessions").InsertOne(ctx, session)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

// SessionExists checks if a session ID exists
func (s *store) SessionExists(ctx context.Context, sessionID string) (bool, error) {
	collection := s.db.Collection("sessions")

	count, err := collection.CountDocuments(ctx, bson.M{"session_id": sessionID})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// SaveVote stores or updates a vote for a given session ID and product ID
func (s *store) SaveVote(ctx context.Context, vote models.Vote) error {
	ctx, cancel := context.WithTimeout(ctx, s.voteTimeout)
	defer cancel()

	filter := bson.M{
		"session_id": vote.SessionID,
		"product_id": vote.ProductID,
	}

	update := bson.M{
		"$set": bson.M{
			"score":      vote.Score,
			"updated_at": time.Now(),
		},
	}

	options := options.Update().SetUpsert(true)

	_, err := s.db.Collection("votes").UpdateOne(ctx, filter, update, options)
	return err
}

// GetVotesBySessionID retrieves all votes associated with a session ID
func (s *store) GetVotesBySessionID(ctx context.Context, sessionID string) ([]models.Vote, error) {
	ctx, cancel := context.WithTimeout(ctx, s.sessionTimeout)
	defer cancel()

	filter := bson.M{
		"session_id": sessionID,
	}

	cursor, err := s.db.Collection("votes").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var votes []models.Vote
	if err := cursor.All(ctx, &votes); err != nil {
		return nil, err
	}

	return votes, nil
}

// GetAggregatedProductScores retrieves aggregated average scores for all products
func (s *store) GetAggregatedProductScores(ctx context.Context) ([]models.ProductScore, error) {
	ctx, cancel := context.WithTimeout(ctx, s.aggregateTimeout)
	defer cancel()

	pipeline := mongo.Pipeline{
		{
			{"$group", bson.D{
				{"_id", "$product_id"},
				{"avg_score", bson.D{{"$avg", "$score"}}},
				{"vote_count", bson.D{{"$sum", 1}}},
			}},
		},
	}

	cursor, err := s.db.Collection("votes").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.ProductScore
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// SaveProducts stores a list of products
func (s *store) SaveProducts(ctx context.Context, products []models.Product) error {
	collection := s.db.Collection("products")

	var docs []interface{}
	for _, product := range products {
		docs = append(docs, product)
	}

	// Clear existing products
	if _, err := collection.DeleteMany(ctx, bson.M{}); err != nil {
		return err
	}

	// Insert new products
	if len(docs) > 0 {
		if _, err := collection.InsertMany(ctx, docs); err != nil {
			return err
		}
	}

	return nil
}

// IsValidProductID checks if a product ID exists
func (s *store) IsValidProductID(ctx context.Context, productID string) (bool, error) {
	collection := s.db.Collection("products")

	count, err := collection.CountDocuments(ctx, bson.M{"product_id": productID})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ensureIndexes creates indexes on the MongoDB collections
func (s *store) ensureIndexes() error {
	ctx := context.Background()

	// Ensure indexes on the sessions collection
	sessionsCollection := s.db.Collection("sessions")
	_, err := sessionsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"session_id": 1,
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create index on sessions collection: %v", err)
	}

	// Ensure indexes on the votes collection
	votesCollection := s.db.Collection("votes")
	_, err = votesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "session_id", Value: 1},
			{Key: "product_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create index on votes collection: %v", err)
	}

	// Ensure indexes on the products collection
	productsCollection := s.db.Collection("products")
	_, err = productsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"product_id": 1,
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create index on products collection: %v", err)
	}

	return nil
}
