package database

import (
	"context"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoConnectionPrefix = "mongodb://"
	mongoLogPrefix        = "mongo"
	mongoTimeout          = 5 * time.Second
)

type mongoDatabase struct {
	client   *mongo.Client
	database string
}

func (m mongoDatabase) Find(collection string) (Cursor, error) {
	logrus.WithField("prefix", mongoLogPrefix).Debugf("connect to database %s, collection %s", m.database, collection)
	c := m.client.Database(m.database).Collection(collection)
	ctx, _ := context.WithTimeout(context.Background(), mongoTimeout)
	return c.Find(ctx, bson.D{})
}

//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//defer cancel()
//
//_, err := c.Collection(CitizenReportCollection).InsertOne(ctx, *data)

func (m *mongoDatabase) InsertOne(collection string, document interface{}) error {
	c := m.client.Database(m.database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), mongoTimeout)
	defer cancel()

	result, err := c.InsertOne(ctx, document)
	if nil != err {
		return err
	}

	logrus.WithField("prefix", mongoLogPrefix).Infof("insert %v result: %v", document, result)
	return nil
}

// NewMongo - new mongo database
func NewMongo(info Info) (Database, error) {
	opts := options.Client().ApplyURI(mongoConnectionString(info))
	client, err := mongo.NewClient(opts)
	if nil != err {
		logrus.WithField("category", mongoLogPrefix).Errorf("new mongo client with error: ", err)
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), mongoTimeout)
	err = client.Connect(ctx)
	if nil != err {
		logrus.WithField("category", mongoLogPrefix).Errorf("connect mongo db with error: ", err)
		return nil, err
	}

	logrus.WithField("prefix", mongoLogPrefix).Infof("connect to mongo database %v", opts.Hosts)

	return &mongoDatabase{client: client, database: info.Database}, nil
}

func mongoConnectionString(info Info) string {
	if info.User == "" {
		return mongoConnectionPrefix + info.IP + ":" + strconv.Itoa(info.Port)
	}
	return mongoConnectionPrefix + info.User + ":" + info.Password + "@" +
		info.IP + ":" + strconv.Itoa(info.Port)
}
