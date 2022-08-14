package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MONGODB_USER = "test"
	MONGODB_PASS = "test"
	MONGODB_IP   = "localhost"
	MONGODB_DB   = "masscan_go"
)

type WhoisDetails struct {
	organisation map[string]string
	autnum       map[string]string
	persons      []map[string]string
}

type DBClient struct {
	client *mongo.Client
	ctx    context.Context
	cancel context.CancelFunc
}

var DB DBClient

// open connection in package init
func init() {
	ctx, cancel := context.WithCancel(context.Background())
	fmt.Printf("%#v", cancel)
	fmt.Printf("mongodb://%s:%s@%s:27017/%s\n",
		MONGODB_USER, MONGODB_PASS, MONGODB_IP, MONGODB_DB)
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:27017/%s",
			MONGODB_USER, MONGODB_PASS, MONGODB_IP, MONGODB_DB)))
	if err != nil {
		panic(err)
	}
	DB = DBClient{
		client: client,
		ctx:    ctx,
		cancel: cancel,
	}
}
