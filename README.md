# Versioned migrations for MongoDB for Golang 1.20

[![Version](https://img.shields.io/badge/version-v1.0.0-blue.svg)](https://github.com/erodriguezg/go-mongodb-migrate/commits/tag/1.0.0) 
[![Go version](https://img.shields.io/badge/go-v1.16-blue.svg)](https://golang.org/doc/devel/release.html#go1.20) 

Fork from: [github.com/xakep666/go-mongodb-migrate](https://github.com/xakep666/go-mongodb-migrate)

This package allows to perform versioned migrations on your MongoDB using [mongo-go-driver](https://github.com/mongodb/mongo-go-driver).
Inspired by [go-pg migrations](https://github.com/go-pg/migrations).

Table of Contents
=================

* [Prerequisites](#prerequisites)
* [Installation](#installation)
* [Usage](#usage)
  * [Use case \#1\. Migrations in files\.](#use-case-1-migrations-in-files)
(#use-case-2-migrations-in-application-code)
* [How it works?](#how-it-works)
* [License](#license)

## Prerequisites
* Golang >= 1.20 

## Installation
```bash
go get -v -u github.com/erodriguezg/go-mongodb-migrate
```

## Usage
### Use case #1. Migrations in files.

* Create a package with migration files.
File name should be like `<version>_<description>.go`.

`1_add-my-index.go`

```go
package migrations

import (
	_ "embed"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	migrate "github.com/erodriguezg/go-mongodb-migrate"
)

//go:embed 1_add-my-index.go
var migration001 string

func init() {
	migrate.Register(&migration001, func(db *mongo.Database) error {
		opt := options.Index().SetName("my-index")
		keys := bson.D{{"my-key", 1}}
		model := mongo.IndexModel{Keys: keys, Options: opt}
		_, err := db.Collection("my-coll").Indexes().CreateOne(context.TODO(), model)
		if err != nil {
			return err
		}

		return nil
	}, func(db *mongo.Database) error {
		_, err := db.Collection("my-coll").Indexes().DropOne(context.TODO(), "my-index")
		if err != nil {
			return err
		}
		return nil
	})
}
```

* Import it in your application.
```go
import (
    ...
    migrate "github.com/erodriguezg/go-mongodb-migrate"
    _ "path/to/migrations_package" // database migrations
    ...
)
```

* Run migrations.
```go
func MongoConnect(host, user, password, database string) (*mongo.Database, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:27017", user, password, host)
	opt := options.Client().ApplyURI(uri)
	client, err := mongo.NewClient(opt)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	db = client.Database(database)
	migrate.SetDatabase(db)
	migrate.SetEnabled(true) // enable o disable the migrations process
	migrate.SetAutoRepair(true) // enable o disable the auto repair process
	if err := migrate.Up(migrate.AllAvailable); err != nil {
		return nil, err
	}
	return db, nil
}
```

## How it works?
This package creates a special collection (by default it`s name is "migrations") for versioning.
In this collection stored documents like
```json
{
    "_id": "<mongodb-generated id>",
    "version": 1,
    "description": "add my-index",
    "timestamp": "<when applied>",
	"hash" : "<the calculate hash of the input file>"
}
```
Current database version determined as version from latest inserted document.

You can change collection name using `SetMigrationsCollection` methods.
Remember that if you want to use custom collection name you need to set it before running migrations.

## What is new?

Taking the fork as the base of the project, new functionality is added:
	
	- Added Drone pipeline
	- Added Makefile
	- Refactor of the code in pkg
	- File hash calculation
	- Auto repair logic
	- Enable / Disable migrations logic  

	