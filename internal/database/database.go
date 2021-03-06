package database

import "context"

// Database - database operations
type Database interface {
	Find(string) (Cursor, error)
	InsertOne(string, interface{}) error
	DeleteOne(string, interface{}) error
	UpdateOne(string, interface{}, interface{}) error
}

// Cursor - operation for database record cursor
type Cursor interface {
	Next(context.Context) bool
	Close(context.Context) error
	Decode(interface{}) error
}

// Info - database connection info
type Info struct {
	IP       string
	Port     int
	User     string
	Password string
	Database string
}
