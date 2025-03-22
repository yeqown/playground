package mongodb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBTestSuite struct {
	suite.Suite
	client     *mongo.Client
	collection *mongo.Collection
	ctx        context.Context
}

func (s *MongoDBTestSuite) SetupSuite() {
	s.ctx = context.Background()
	client, err := mongo.Connect(s.ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	s.Require().NoError(err)
	s.client = client
	s.collection = client.Database("test").Collection("users")

	// 测试前清空集合
	_, err = s.collection.DeleteMany(s.ctx, bson.M{})
	s.Require().NoError(err)
}

func (s *MongoDBTestSuite) TearDownSuite() {
	s.Require().NoError(s.client.Disconnect(s.ctx))
}

func (s *MongoDBTestSuite) TestConnection() {
	err := s.client.Ping(s.ctx, nil)
	s.Require().NoError(err)
}

func (s *MongoDBTestSuite) TestInsertDocument() {
	user := User{
		Name:      "Test User",
		Age:       25,
		CreatedAt: time.Now(),
	}

	// 测试插入操作
	result, err := s.collection.InsertOne(s.ctx, user)
	s.Require().NoError(err)
	s.Require().NotNil(result.InsertedID)

	// 验证插入的文档
	var foundUser User
	err = s.collection.FindOne(s.ctx, bson.M{"name": "Test User"}).Decode(&foundUser)
	s.Require().NoError(err)
	s.Equal(user.Name, foundUser.Name)
	s.Equal(user.Age, foundUser.Age)
}

func (s *MongoDBTestSuite) TestInsertBsonRaw() {

	user := User{
		Name:      "raw",
		Age:       25,
		CreatedAt: time.Now(),
	}
	raw, err := bson.Marshal(user)
	s.Require().NoError(err)

	// 测试插入操作
	result, err := s.collection.InsertOne(s.ctx, raw)
	s.Require().NoError(err)
	s.Require().NotNil(result.InsertedID)

	// 验证插入的文档
	var foundUser User
	err = s.collection.FindOne(s.ctx, bson.M{"name": "raw"}).Decode(&foundUser)
	s.Require().NoError(err)
	s.Equal(user.Name, foundUser.Name)
	s.Equal(user.Age, foundUser.Age)
}

func (s *MongoDBTestSuite) TestInsertManyBsonRaw() {

	user := User{
		Name:      "raw",
		Age:       25,
		CreatedAt: time.Now(),
	}
	raw, err := bson.Marshal(user)
	s.Require().NoError(err)

	// 测试插入操作
	result, err := s.collection.InsertMany(s.ctx, []interface{}{raw, raw})
	s.Require().NoError(err)
	s.Require().NotNil(result.InsertedIDs)

	// 验证插入的文档
	count, err := s.collection.CountDocuments(s.ctx, bson.M{"name": "raw"})
	s.Require().NoError(err)
	s.GreaterOrEqual(2, int(count))
}

type Address struct {
	Street     string `bson:"street"`
	City       string `bson:"city"`
	Country    string `bson:"country"`
	PostalCode string `bson:"postal_code"`
}

type Contact struct {
	Email     string   `bson:"email"`
	Phone     string   `bson:"phone"`
	Emergency []string `bson:"emergency_contacts"`
}

type User struct {
	Name        string                 `bson:"name"`
	Age         int                    `bson:"age"`
	CreatedAt   time.Time              `bson:"created_at"`
	Address     Address                `bson:"address"`
	Contact     Contact                `bson:"contact"`
	Tags        []string               `bson:"tags"`
	Preferences map[string]interface{} `bson:"preferences"`
	LastLogins  []time.Time            `bson:"last_logins"`
	IsActive    bool                   `bson:"is_active"`
}

// Insert 100000 documents using struct directly, total time: 704.342917ms, avg: 70.434291ms per batch
func (s *MongoDBTestSuite) TestInsertManyBench_marshal() {
	const (
		batchSize = 10000
		rounds    = 10 // 总共插入 10w 条记录
	)

	users := make([]interface{}, batchSize)
	for i := 0; i < batchSize; i++ {
		users[i] = User{
			Name:      fmt.Sprintf("user_%d", i),
			Age:       25 + (i % 50),
			CreatedAt: time.Now(),
			Address: Address{
				Street:     fmt.Sprintf("Street %d", i),
				City:       "Beijing",
				Country:    "China",
				PostalCode: fmt.Sprintf("1000%d", i%99),
			},
			Contact: Contact{
				Email:     fmt.Sprintf("user_%d@example.com", i),
				Phone:     fmt.Sprintf("138%08d", i),
				Emergency: []string{"110", "119", "120"},
			},
			Tags: []string{"vip", "new", "active"},
			Preferences: map[string]interface{}{
				"theme":         "dark",
				"language":      "zh-CN",
				"timezone":      "Asia/Shanghai",
				"notifications": true,
			},
			LastLogins: []time.Time{
				time.Now().Add(-24 * time.Hour),
				time.Now().Add(-48 * time.Hour),
			},
			IsActive: true,
		}
	}

	start := time.Now()
	for i := 0; i < rounds; i++ {
		result, err := s.collection.InsertMany(s.ctx, users)
		s.Require().NoError(err)
		s.Require().Len(result.InsertedIDs, batchSize)
	}
	duration := time.Since(start)

	s.T().Logf("Insert %d documents using struct directly, total time: %v, avg: %v per batch",
		batchSize*rounds, duration, duration/time.Duration(rounds))
}

// Insert 100000 documents using bson.Raw, total time: 392.836334ms, avg: 39.283633ms per batch
func (s *MongoDBTestSuite) TestInsertManyBench_bsonRaw() {
	const (
		batchSize = 10000
		rounds    = 10 // 总共插入 10w 条记录
	)

	rawDocs := make([]interface{}, batchSize)
	for i := 0; i < batchSize; i++ {
		user := User{
			Name:      fmt.Sprintf("user_%d", i),
			Age:       25 + (i % 50),
			CreatedAt: time.Now(),
			Address: Address{
				Street:     fmt.Sprintf("Street %d", i),
				City:       "Beijing",
				Country:    "China",
				PostalCode: fmt.Sprintf("1000%d", i%99),
			},
			Contact: Contact{
				Email:     fmt.Sprintf("user_%d@example.com", i),
				Phone:     fmt.Sprintf("138%08d", i),
				Emergency: []string{"110", "119", "120"},
			},
			Tags: []string{"vip", "new", "active"},
			Preferences: map[string]interface{}{
				"theme":         "dark",
				"language":      "zh-CN",
				"timezone":      "Asia/Shanghai",
				"notifications": true,
			},
			LastLogins: []time.Time{
				time.Now().Add(-24 * time.Hour),
				time.Now().Add(-48 * time.Hour),
			},
			IsActive: true,
		}
		raw, err := bson.Marshal(user)
		s.Require().NoError(err)
		rawDocs[i] = raw
	}

	start := time.Now()
	for i := 0; i < rounds; i++ {
		result, err := s.collection.InsertMany(s.ctx, rawDocs)
		s.Require().NoError(err)
		s.Require().Len(result.InsertedIDs, batchSize)
	}
	duration := time.Since(start)

	s.T().Logf("Insert %d documents using bson.Raw, total time: %v, avg: %v per batch",
		batchSize*rounds, duration, duration/time.Duration(rounds))
}

func TestMongoDBSuite(t *testing.T) {
	suite.Run(t, new(MongoDBTestSuite))
}
