package database

import (
	"context"
	"fmt"
	"go-graphql-mongodb-api/graph/model"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	client *mongo.Client
}

func Connect(url string) *DB {
	// generates a new client to connect to the deployment
	client, err := mongo.NewClient(options.Client().ApplyURI(url))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	// test if the connection to the mongo.Client was created successfully
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	// if the connection is successful, return the DB struct
	return &DB{
		client: client,
	}
}

// function for inserting a new Kolo into database
func (db *DB) InsertKolo(kolo model.NewKolo) *model.Kolo {
	koloColl := db.client.Database("graphql-mongodb-api-db").Collection("kolo")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	inserg, err := koloColl.InsertOne(ctx, bson.D{{Key: "serijska_stevilka", Value: kolo.SerijskaStevilka}, {Key: "mnenje", Value: bson.A{}}})

	if err != nil {
		log.Fatal(err)
	}

	insertedID := inserg.InsertedID.(primitive.ObjectID).Hex()
	returnKolo := model.Kolo{ID: insertedID, SerijskaStevilka: kolo.SerijskaStevilka, Mnenje: make([]*string, 500)}

	return &returnKolo
}

// function for finding a kolo in database
func (db *DB) FindKolo(id string) *model.Kolo {
	ObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}

	koloColl := db.client.Database("graphql-mongodb-api-db").Collection("kolo")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := koloColl.FindOne(ctx, bson.M{"_id": ObjectID})

	kolo := model.Kolo{ID: id}

	result.Decode(&kolo)

	return &kolo
}

// function for finding all Kolo in database
func (db *DB) FindAllKolo() []*model.Kolo {
	koloColl := db.client.Database("graphql-mongodb-api-db").Collection("kolo")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := koloColl.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var kolesa []*model.Kolo
	for cursor.Next(ctx) {
		sus, err := cursor.Current.Elements()
		fmt.Println(sus)
		if err != nil {
			log.Fatal(err)
		}

		kolo := model.Kolo{ID: sus[0].Value().StringValue(), SerijskaStevilka: sus[1].Value().StringValue(), Mnenje: make([]*string, 500)}

		kolesa = append(kolesa, &kolo)
	}

	return kolesa
}
