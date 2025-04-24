package main

import (
	"context"
	"flag"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	insertFlag = false
)

func main() {
	flag.BoolVar(&insertFlag, "insert", false, "insert flag")
	flag.Parse()

	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("test").Collection("coll20250424")

	if insertFlag {
		insert(coll)
	}

	// 查询数据
	query(coll)
}

func insert(coll *mongo.Collection) {

	// 插入数据
	// {
	// 	"resource_id": "1234567890",
	// 	"resource_type": "image",
	// 	"ownership_tags": [
	// 		{
	// 			"Key": "Owner",
	// 			"Value": "John Doe"
	// 		},
	// 		{
	// 			"Key": "Department",
	// 			"Value": "Engineering"
	// 		},
	// 		{
	// 			"Key": "Project",
	// 			"Value": "ABC123"
	// 		}
	// 	]
	// }

	doc := map[string]interface{}{
		"resource_id":   "1234567890",
		"resource_type": "image",
		"ownership_tags": []map[string]string{
			{
				"Key":   "Owner",
				"Value": "John Doe",
			},
			{
				"Key":   "Department",
				"Value": "Engineering",
			},
			{
				"Key":   "Project",
				"Value": "ABC123",
			},
		},
	}

	_, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		panic(err)
	}
}

func query(coll *mongo.Collection) {
	filter := bson.M{
		"ownership_tags": bson.M{
			"$all": bson.A{
				bson.M{
					"$elemMatch": bson.M{
						"Key":   "Owner",
						"Value": "John Doe",
					},
				},
				bson.M{
					"$elemMatch": bson.M{
						"Key":   "Department",
						"Value": "Engineering",
					},
				},
			},
		},
	}

	cur, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	defer cur.Close(context.TODO())

	var docs []bson.M
	if err := cur.All(context.TODO(), &docs); err != nil {
		panic(err)
	}

	// count filtered docs
	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	fmt.Println("count: ", count)
	fmt.Println("docs: ", len(docs))
}
