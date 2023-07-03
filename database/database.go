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
		log.Default().Println("ERROR -> couldn't insert Kolo into database")
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

	if err != nil {
		log.Fatal(err)
	}

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

// function for inserting a Postajalisce into database
func (db *DB) InsertPostajalisce(postajalisce model.NewPostajalisce) *model.Postajalisce {
	postajalisceColl := db.client.Database(DATABASE).Collection(POSTAJALISCA)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	inserg, err := postajalisceColl.InsertOne(ctx, bson.D{{Key: "ime", Value: postajalisce.Ime}, {Key: "naslov", Value: postajalisce.Naslov},
		{Key: "latitude", Value: postajalisce.Latitude}, {Key: "longitude", Value: postajalisce.Longitude}, {Key: "kolesaArray", Value: bson.A{}}})

	if err != nil {
		log.Default().Println("ERROR -> couldn't insert Postajalisce into database.")
		log.Fatal(err)
	}

	insertedID := inserg.InsertedID.(primitive.ObjectID).Hex()
	returnPostajalisce := model.Postajalisce{ID: insertedID, Ime: postajalisce.Ime, Naslov: postajalisce.Naslov, Latitude: postajalisce.Latitude, Longitude: postajalisce.Longitude, KolesaArray: []*model.Kolo{}}

	return &returnPostajalisce
}

// function for finding a Postajalisce in database by id
func (db *DB) FindPostajalisce(id string) *model.Postajalisce {
	ObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Default().Println("ERROR -> couldn't convert id to ObjectID")
		log.Fatal(err)
	}

	postajalisceColl := db.client.Database(DATABASE).Collection(POSTAJALISCA)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := postajalisceColl.FindOne(ctx, bson.M{"_id": ObjectID})

	postajalisce := model.Postajalisce{}

	err = result.Decode(&postajalisce)

	if err != nil {
		log.Default().Println("ERROR -> couldn't decode result into Postajalisce")
	}

	return &postajalisce
}

// function for finding all Postajalisce in database

func (db *DB) FindAllPostajalisce() []*model.Postajalisce {
	postajalisceColl := db.client.Database(DATABASE).Collection(POSTAJALISCA)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := postajalisceColl.Find(ctx, bson.D{})
	if err != nil {
		log.Default().Println("ERROR -> couldn't create cursor for Postajalisce")
		log.Fatal(err)
	}

	var postajalisca []*model.Postajalisce
	for cursor.Next(ctx) {
		var postajalisce *model.Postajalisce
		err := cursor.Decode(&postajalisce)

		if err != nil {
			log.Default().Println("ERROR -> couldn't decode result into Postajalisce")
			log.Fatal(err)
		}

		postajalisca = append(postajalisca, postajalisce)
	}

	return postajalisca
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
