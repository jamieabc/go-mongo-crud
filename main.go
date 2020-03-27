package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/spf13/viper"
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

type serverInfo struct {
	ip       string
	port     int
	database string
}

type config struct {
	server serverInfo
}

func init() {
	flag.StringVar(&confFile, "c", "", "config file path")
	flag.Parse()
}

func main() {
	if confFile == "" {
		help()
		os.Exit(errExitCode)
	}

	info, err := parseLog()
	if nil != err {
		fmt.Println("parse log with error: ", err)
		os.Exit(errExitCode)
	}

	fmt.Println("server: ", info.server)

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + info.server.ip + ":" + strconv.Itoa(info.server.port)))
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

	collection := client.Database(info.server.database).Collection("places")
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

func parseLog() (config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(confFile)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if nil != err {
		return config{}, err
	}

	return config{
		serverInfo{
			ip:       viper.GetString("server.ip"),
			port:     viper.GetInt("server.port"),
			database: viper.GetString("server.database"),
		},
	}, nil
}
