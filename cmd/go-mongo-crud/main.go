package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

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
	TypeString  string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}

type recordStruct struct {
	Name     string         `bson:"name"`
	Location locationStruct `bson:"location"`
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

	cur, err := m.Find("places")
	if nil != err {
		logrus.WithField("prefix", "main").Panicf("find with error: %s", err)
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		logrus.WithField("prefix", logPrefix).Debug("finding database records")
		var record recordStruct
		err := cur.Decode(&record)
		if nil != err {
			logrus.WithField("prefix", logPrefix).Errorf("decode record with error: ", err)
			continue
		}
		logrus.WithField("prefix", logPrefix).Debugf("record: %v", record)
	}
}

func help() {
	fmt.Println("Usage: go-mongo-crud -c config_file_path")
}
