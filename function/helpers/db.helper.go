package helpers

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBConnection struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var once sync.Once
var dbError error

func (dbc *DBConnection) GetInstanceDB() (*DBConnection, error) {
	once.Do(func() {
		if dbc.Db == nil {
			uri := os.Getenv("MONGODB_URI")
			dbName := os.Getenv("MONGO_DB_NAME")
			// https://www.mongodb.com/docs/drivers/go/current/fundamentals/logging/
			//loggerOptions := options.Logger().SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			//client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetLoggerOptions(loggerOptions))
			client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
			if err != nil {
				dbError = err
				log.Printf("error to connect with mongo: %s", err.Error())
			}
			dbc.Client = client
			dbc.Db = dbc.Client.Database(dbName)
		}
	})
	return dbc, dbError
}

func (dbc *DBConnection) CloseConnection() error {
	if dbc.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := dbc.Client.Disconnect(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
