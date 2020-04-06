package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/bxcodec/faker/v3"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/jamieabc/go-mongo-crud/internal/config"
	"github.com/jamieabc/go-mongo-crud/internal/database"
)

const (
	errExitCode = -1
	logPrefix   = "main"
)

var confFile string

func init() {
	flag.StringVar(&confFile, "c", "", "config file path")
	flag.Parse()

	// Log as JSON instead of the default ASCII formatter.
	//logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log. the warning severity or above.
	logrus.SetLevel(logrus.DebugLevel)
}

type locationStruct struct {
	TypeString  string    `bson:"type" faker:"-"`
	Coordinates []float64 `bson:"coordinates" faker:"-"`
}

type recordStruct struct {
	Name     string         `bson:"name" faker:"word"`
	Location locationStruct `bson:"location" faker:"-"`
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

	m, err := database.NewMongo(database.Info{
		IP:       info.Server.IP,
		Port:     info.Server.Port,
		User:     "",
		Password: "",
		Database: info.Server.Database,
	})
	if nil != err {
		logrus.WithField("prefix", "main").Panicf("create mongo with error: %s", err)
	}

	// insert
	r := recordStruct{}
	_ = faker.FakeData(&r)
	r.Location = locationStruct{
		TypeString:  "Point",
		Coordinates: []float64{-1, -1},
	}

	err = m.InsertOne("places", &r)
	if nil != err {
		logrus.WithField("prefix", logPrefix).Panicf("insert record with error: %s", err)
	}

	// search
	cur, err := m.Find("places")
	if nil != err {
		logrus.WithField("prefix", "main").Panicf("find with error: %s", err)
	}
	defer cur.Close(context.Background())

	var toDelete string

	for cur.Next(context.Background()) {
		var record recordStruct
		err := cur.Decode(&record)
		if toDelete == "" {
			toDelete = record.Name
		}

		if nil != err {
			logrus.WithField("prefix", logPrefix).Errorf("decode record with error: ", err)
			continue
		}
		logrus.WithField("prefix", logPrefix).Debugf("record: %v", record)
	}

	// update
	err = m.UpdateOne(
		"places",
		bson.M{"name": bson.M{"$eq": toDelete}},
		bson.M{"$set": bson.M{"name": "toDelete"}},
	)

	// delete
	err = m.DeleteOne("places", bson.M{"name": "toDelete"})
	if nil != err {
		logrus.WithField("prefix", logPrefix).Panicf("delete record %s with error: %s", toDelete, err)
	}

}

func help() {
	fmt.Println("Usage: go-mongo-crud -c config_file_path")
}
