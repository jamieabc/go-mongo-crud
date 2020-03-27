package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jamieabc/go-mongo-crud/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strconv"
	"time"
)

const (
	errExitCode = -1
	dbTimeout   = 5 * time.Second
)

var confFile string

func init() {
	flag.StringVar(&confFile, "c", "", "config file path")
	flag.Parse()
}

func main() {
	if confFile == "" {
		help()
		os.Exit(errExitCode)
	}

	info, err := config.Parse(confFile)
	if nil != err {
		fmt.Println("parse log with error: ", err)
		os.Exit(errExitCode)
	}

	fmt.Println("server: ", info.Server)

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + info.Server.IP + ":" + strconv.Itoa(info.Server.Port)))
	if nil != err {
		fmt.Println("new mongo client with error: ", err)
		os.Exit(errExitCode)
	}

	ctx, _ := context.WithTimeout(context.Background(), dbTimeout)
	err = client.Connect(ctx)
	if nil != err {
		fmt.Println("connect mongo db with error: ", err)
		os.Exit(errExitCode)
	}

	collection := client.Database(info.Server.Database).Collection("places")
	cur, err := collection.Find(ctx, bson.D{})
	if nil != err {
		fmt.Println("find with error: ", err)
		os.Exit(errExitCode)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var record bson.M
		err := cur.Decode(&record)
		if nil != err {
			fmt.Println("decode record with error: ", err)
			continue
		}
		fmt.Println("record: ", record)
	}
}

func help() {
	fmt.Println("Usage: go-mongo-crud -c config_file_path")
}
