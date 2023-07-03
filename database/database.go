package database

import (
	"context"
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

const (
	DATABASE     = "projekt"
	KOLESA       = "kolesa"
	POSTAJALISCA = "postajalisca"
	DOGS         = "dogs"
)

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

	log.Default().Println("Successfully connected to database!")

	// if the connection is successful, return the DB struct
	return &DB{
		client: client,
	}
}

// function for inserting a new Kolo into database
func (db *DB) InsertKolo(kolo model.NewKolo) *model.Kolo {
	koloColl := db.client.Database(DATABASE).Collection(KOLESA)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	inserg, err := koloColl.InsertOne(ctx, bson.D{{Key: "serijska_stevilka", Value: kolo.SerijskaStevilka}, {Key: "mnenje", Value: bson.A{}}})

	if err != nil {
		log.Default().Println("Napaka pri vstavljanju")
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

	koloColl := db.client.Database(DATABASE).Collection(KOLESA)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := koloColl.FindOne(ctx, bson.M{"_id": ObjectID})

	kolo := model.Kolo{}

	err = result.Decode(&kolo)

	log.Default().Println(kolo.ID)
	log.Default().Println(kolo.Mnenje)
	log.Default().Println(kolo.SerijskaStevilka)

	return &kolo
}

// function for finding all Kolo in database
func (db *DB) FindAllKolo() []*model.Kolo {
	koloColl := db.client.Database(DATABASE).Collection(KOLESA)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := koloColl.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var kolesa []*model.Kolo
	for cursor.Next(ctx) {
		var kolo *model.Kolo
		err := cursor.Decode(&kolo)

		if err != nil {
			log.Default().Println("Napaka pri kurzorju")
			log.Fatal(err)
		}

		kolesa = append(kolesa, kolo)
	}

	return kolesa
}

func (db *DB) InsertDog(input *model.NewDog) *model.Dog {
	collection := db.client.Database(DATABASE).Collection(KOLESA)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, input)

	if err != nil {
		log.Fatal(err)
	}

	return &model.Dog{
		ID:        res.InsertedID.(primitive.ObjectID).Hex(),
		Name:      input.Name,
		IsGoodBoi: input.IsGoodBoi,
	}
}

func (db *DB) FindDog(id string) *model.Dog {
	ObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}

	collection := db.client.Database(DATABASE).Collection(DOGS)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res := collection.FindOne(ctx, bson.M{"_id": ObjectID})

	dog := model.Dog{}

	res.Decode(&dog)

	return &dog
}

func (db *DB) FindAllDog() []*model.Dog {
	collection := db.client.Database(DATABASE).Collection(DOGS)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var dogs []*model.Dog
	for cursor.Next(ctx) {
		var dog *model.Dog
		err := cursor.Decode(dog)
		if err != nil {
			log.Fatal(err)
		}
		dogs = append(dogs, dog)
	}

	return dogs

}
