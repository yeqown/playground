package mongodb_test

import (
	"context"
	"encoding/json"
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

var jsonTemplate = `{
	"name": "jhon",
	"age": 10,
	"created_at": "2023-08-21T10:00:00Z",
	"address": {
		"street": "Street xxx",
		"city": "Beijing",
		"country": "China",
		"postal_code": "100000"
	},
	"contact": {
		"email": "user_1000@example.com",
		"phone": "13812308123",
		"emergency_contacts": ["110", "119", "120"]
	},
	"tags": ["vip", "new", "active"],
	"preferences": {
		"theme": "dark",
		"language": "zh-CN",
		"timezone": "Asia/Shanghai",
		"notifications": true
	},
	"last_logins": ["2023-08-20T10:00:00Z", "2023-08-19T10:00:00Z"],
	"is_active": true
}`

// Insert 100000 documents using struct directly, total time: 1.721432459s, avg: 172.143245ms per batch
func (s *MongoDBTestSuite) TestInsertManyBench_marshal() {
	// 准备 JSON 模板

	const (
		batchSize = 10000
		rounds    = 10 // 总共插入 10w 条记录
	)

	// 根据 batchSize 从 jsonTemplate 从解析出 map[string]interface{}
	unmarshalJson := func(batchSize int) []interface{} {
		users := make([]interface{}, batchSize)
		for i := 0; i < batchSize; i++ {
			u := make(map[string]interface{}, 16)
			err := json.Unmarshal([]byte(jsonTemplate), &u)
			s.Require().NoError(err)
			users[i] = u
		}

		return users
	}

	start := time.Now()
	for i := 0; i < rounds; i++ {
		users := unmarshalJson(batchSize)
		result, err := s.collection.InsertMany(s.ctx, users)
		s.Require().NoError(err)
		s.Require().Len(result.InsertedIDs, batchSize)
	}
	duration := time.Since(start)

	s.T().Logf("Insert %d documents using struct directly, total time: %v, avg: %v per batch",
		batchSize*rounds, duration, duration/time.Duration(rounds))
}

// Insert 100000 documents using bson.Raw, total time: 453.779208ms, avg: 45.37792ms per batch
func (s *MongoDBTestSuite) TestInsertManyBench_bsonRaw() {
	const (
		batchSize = 10000
		rounds    = 10 // 总共插入 10w 条记录
	)

	u := make(map[string]interface{}, 16)
	err := json.Unmarshal([]byte(jsonTemplate), &u)
	s.Require().NoError(err)

	raw, err := bson.Marshal(u)
	s.Require().NoError(err)

	rawDocs := make([]interface{}, batchSize)
	for i := 0; i < batchSize; i++ {
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
