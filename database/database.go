package database

import (
	"context"
	"go-graphql-mongodb-api/graph/model"
	"log"
	"math"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// type of database
type DB struct {
	client *mongo.Client
}

// type for Postajalisce and the distance to the location of the user
type Razdalja struct {
	postajalisce *model.Postajalisce
	distance     float64
}

// type of sorted array of Razdalja
type ByDistance []*Razdalja

// functions for sorting the array

func (a ByDistance) Len() int           { return len(a) }
func (a ByDistance) Less(i, j int) bool { return a[i].distance < a[j].distance }
func (a ByDistance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// database collection names constants
const (
	DATABASE     = "projekt"
	KOLESA       = "kolesa"
	POSTAJALISCA = "postajalisca"
	IZPOSOJE     = "izposoje"
)

func Connect(url string) *DB {
	// generates a new client to connect to the deployment
	client, err := mongo.NewClient(options.Client().ApplyURI(url))

	if err != nil {
		log.Default().Println("ERROR -> couldn't create a client")
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Default().Println("ERROR -> couldn't connect to database")
		log.Fatal(err)
	}

	// test if the connection to the mongo.Client was created successfully
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		log.Default().Println("ERROR -> couldn't ping to database")
		log.Fatal(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Default().Println("ERROR -> couldn't ping to database")
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
	inserg, err := koloColl.InsertOne(ctx, bson.D{{Key: "serijska_stevilka", Value: kolo.SerijskaStevilka}, {Key: "mnenje", Value: bson.A{}},
		{Key: "rezervirano", Value: false}, {Key: "jeIzposojen", Value: false}})

	if err != nil {
		log.Default().Println("ERROR -> couldn't insert Kolo into database")
		log.Fatal(err)
	}

	insertedID := inserg.InsertedID.(primitive.ObjectID).Hex()
	returnKolo := model.Kolo{ID: insertedID, SerijskaStevilka: kolo.SerijskaStevilka, Mnenje: make([]*int, 500)}

	return &returnKolo
}

// function for updating a Kolo in database
func (db *DB) UpdateKolo(kolo model.UpdateKolo) *model.Kolo {
	koloColl := db.client.Database(DATABASE).Collection(KOLESA)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ObjectID, err := primitive.ObjectIDFromHex(kolo.ID)
	if err != nil {
		log.Default().Println("ERROR -> couldn't convert id to ObjectID")
		log.Fatal()
	}

	filter := bson.M{"_id": ObjectID}

	updatedKolo := bson.M{}

	if kolo.SerijskaStevilka != nil {
		updatedKolo["serijska_stevilka"] = kolo.SerijskaStevilka
	}

	if kolo.Mnenje != nil {
		updatedKolo["mnenje"] = kolo.Mnenje
	}

	update := bson.M{"$set": updatedKolo}

	_, err = koloColl.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Default().Println("ERROR -> couldn't update Kolo in database")
		log.Fatal(err)
	}

	var novoKolo = db.FindKolo(kolo.ID)

	if novoKolo.JeIzposojen != true {
		var posodobljenoKolo = model.KoloInput{
			ID:               kolo.ID,
			SerijskaStevilka: novoKolo.SerijskaStevilka,
			Mnenje:           novoKolo.Mnenje,
			JeIzposojen:      novoKolo.JeIzposojen,
		}

		db.UpdateKoloInPostajalisce(kolo.ID, &posodobljenoKolo)
	}

	return db.FindKolo(kolo.ID)
}

// function for deleting a Kolo from database
func (db *DB) DeleteKolo(id string) string {
	koloColl := db.client.Database(DATABASE).Collection(KOLESA)

	ObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Default().Println("ERROR -> couldn't convert id to ObjectID")
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": ObjectID}

	_, err = koloColl.DeleteOne(ctx, filter)
	if err != nil {
		log.Default().Println("ERROR -> couldn't delete Kolo from database")
		log.Fatal(err)
	}

	return "OK -> Successfully deleted Kolo from database."
}

// function for finding a kolo in database
func (db *DB) FindKolo(id string) *model.Kolo {
	ObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Default().Println("ERROR -> couldn't convert id to ObjectID")
		log.Fatal(err)
	}

	koloColl := db.client.Database(DATABASE).Collection(KOLESA)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := koloColl.FindOne(ctx, bson.M{"_id": ObjectID})

	kolo := model.Kolo{}

	err = result.Decode(&kolo)

	if err != nil {
		log.Default().Println("ERROR -> couldn't decode object Kolo")
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
		log.Default().Println("ERROR -> couldn't create cursor on collection Kolesa")
		log.Fatal(err)
	}

	var kolesa []*model.Kolo
	for cursor.Next(ctx) {
		var kolo *model.Kolo
		err := cursor.Decode(&kolo)

		if err != nil {
			log.Default().Println("ERROR -> couldn't decode cursor")
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

// function for updating a Postajalisce in database
func (db *DB) UpdatePostajalisce(postajalisce model.UpdatePostajalisce) *model.Postajalisce {
	postajalisceColl := db.client.Database(DATABASE).Collection(POSTAJALISCA)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ObjectID, err := primitive.ObjectIDFromHex(postajalisce.ID)
	if err != nil {
		log.Default().Println("ERROR -> couldn't convert id to ObjectID")
		log.Fatal(err)
	}

	filter := bson.M{"_id": ObjectID}

	update := bson.M{}

	if postajalisce.Ime != nil {
		update["ime"] = postajalisce.Ime
	}

	if postajalisce.Naslov != nil {
		update["naslov"] = postajalisce.Naslov
	}

	if postajalisce.Latitude != nil {
		update["latitude"] = postajalisce.Latitude
	}

	if postajalisce.Longitude != nil {
		update["longitude"] = postajalisce.Longitude
	}

	if postajalisce.KolesaArray != nil {
		update["kolesaArray"] = postajalisce.KolesaArray
	}

	updatePost := bson.M{"$set": update}

	_, err = postajalisceColl.UpdateOne(ctx, filter, updatePost)
	if err != nil {
		log.Default().Println("ERROR -> couldn't update Postajalisce")
		log.Fatal(err)
	}

	return db.FindPostajalisce(postajalisce.ID)
}

// function for deleting a Postajalisce from database
func (db *DB) DeletePostajalisce(id string) string {
	postajalisceColl := db.client.Database(DATABASE).Collection(POSTAJALISCA)

	ObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Default().Println("ERROR -> couldn't convert id to ObjectID")
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": ObjectID}

	_, err = postajalisceColl.DeleteOne(ctx, filter)
	if err != nil {
		log.Default().Println("ERROR -> couldn't delete Postajalisce from database")
		log.Fatal(err)
	}

	return "OK -> Successfully deleted Postajalisce from database"
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
		log.Fatal(err)
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

// function for finding nearest Postajalisce to a given location
func (db *DB) FindNearestPostajalisce(latitude float64, longitude float64, stPostajalisc int) []*model.Postajalisce {
	postajalisceColl := db.client.Database(DATABASE).Collection(POSTAJALISCA)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := postajalisceColl.Find(ctx, bson.D{})
	if err != nil {
		log.Default().Println("ERROR -> couldn't create cursor for Postajalisce")
		log.Fatal(err)
	}

	var postajalisca []*Razdalja
	for cursor.Next(ctx) {
		var postajalisce *model.Postajalisce
		err := cursor.Decode(&postajalisce)

		if err != nil {
			log.Default().Println("ERROR -> couldn't decode result into Postajalisce")
			log.Fatal(err)
		}

		var razdalja = distance(latitude, longitude, postajalisce.Latitude, postajalisce.Longitude)

		postajalisca = append(postajalisca, &Razdalja{postajalisce, razdalja})

	}

	// Sort the array by distance
	sort.Sort(ByDistance(postajalisca))

	for _, distance := range postajalisca {
		log.Default().Println(distance.postajalisce.Ime, distance)
	}

	var urejenaPostajalisca []*model.Postajalisce

	for i, postajalisce := range postajalisca {
		if i < stPostajalisc {
			urejenaPostajalisca = append(urejenaPostajalisca, postajalisce.postajalisce)
		}
	}

	return urejenaPostajalisca
}

// function for getting Postajalisce, in which the Kolo is stored
func (db *DB) GetPostajalisceFromKolo(id string) *model.Postajalisce {
	postajalisceColl := db.client.Database(DATABASE).Collection(POSTAJALISCA)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := postajalisceColl.Find(ctx, bson.D{})
	if err != nil {
		log.Default().Println("ERROR -> couldn't create cursor for Postajalisce")
		log.Fatal(err)
	}

	for cursor.Next(ctx) {
		var postajalisce *model.Postajalisce
		err := cursor.Decode(&postajalisce)

		if err != nil {
			log.Default().Println("ERROR -> couldn't decode result into Postajalisce")
			log.Fatal(err)
		}

		for _, kolo := range postajalisce.KolesaArray {
			if kolo.ID == id {
				return db.FindPostajalisce(postajalisce.ID)
			}
		}
	}

	return nil
}

// function for borrowing a Kolo from Postajalisce
func (db *DB) BorrowKolo(input model.IzposojaKolesa) *model.Izposoja {
	izposojeColl := db.client.Database(DATABASE).Collection(IZPOSOJE)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	postajalisce := db.GetPostajalisceFromKolo(input.BikeID)

	updatedKolesaArray := make([]*model.KoloInput, 0)
	for _, kolo := range postajalisce.KolesaArray {
		if kolo.ID != input.BikeID {
			updatedKolo := model.KoloInput{ID: kolo.ID, SerijskaStevilka: kolo.SerijskaStevilka, Mnenje: kolo.Mnenje, JeIzposojen: kolo.JeIzposojen}
			updatedKolesaArray = append(updatedKolesaArray, &updatedKolo)
		}
	}

	updatedPostajalisce := model.UpdatePostajalisce{
		ID:          postajalisce.ID,
		Ime:         &postajalisce.Ime,
		Naslov:      &postajalisce.Naslov,
		Latitude:    &postajalisce.Latitude,
		Longitude:   &postajalisce.Longitude,
		KolesaArray: updatedKolesaArray,
	}

	db.UpdatePostajalisce(updatedPostajalisce)

	inserg, err := izposojeColl.InsertOne(ctx, bson.D{{Key: "start_date", Value: formattedTime}, {Key: "start_station_id", Value: postajalisce.ID},
		{Key: "bike_id", Value: input.BikeID}, {Key: "trenutna_zasedenost_start", Value: len(postajalisce.KolesaArray)}, {Key: "weather", Value: input.Weather}, {Key: "start_station", Value: postajalisce.Ime},
		{Key: "username", Value: input.Username}})

	if err != nil {
		log.Default().Println("ERROR -> couldn't insert Izposoja into database")
		log.Fatal(err)
	}

	insertedID := inserg.InsertedID.(primitive.ObjectID).Hex()
	returnIzposoja := model.Izposoja{ID: insertedID, StartDate: formattedTime, StartStationID: postajalisce.ID, BikeID: input.BikeID, Weather: input.Weather, Username: input.Username,
		TrenutnaZasedenostStart: len(postajalisce.KolesaArray), StartStation: postajalisce.Ime}

	return &returnIzposoja
}

// function for returning a Kolo to Postajalisce
func (db *DB) ReturnKolo(input model.VraciloKolesa) *model.Izposoja {
	izposojeColl := db.client.Database(DATABASE).Collection(IZPOSOJE)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ObjectID, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		log.Default().Println("ERROR -> couldn't convert id to ObjectID")
		log.Fatal(err)
	}

	trenutniCas := time.Now()
	formattedTime := trenutniCas.Format("2006-01-02 15:04:05")

	currentTimeUnix := time.Now().Unix()

	borrow := db.FindIzposoja(input.ID)

	parsedTime, err := time.Parse("2006-01-02 15:04:05", borrow.StartDate)
	if err != nil {
		log.Default().Println("ERROR -> couldn't parse formatted time")
		log.Fatal(err)
	}

	trajanje := parsedTime.Unix() - currentTimeUnix

	postajalisce := db.FindPostajalisce(input.EndStationID)

	kolo := db.FindKolo(input.BikeID)

	novoKolo := model.KoloInput{ID: kolo.ID, SerijskaStevilka: kolo.SerijskaStevilka, Mnenje: kolo.Mnenje, JeIzposojen: kolo.JeIzposojen}

	updatedKolesaArray := make([]*model.KoloInput, 0)
	for _, kolo := range postajalisce.KolesaArray {
		updatedKolo := model.KoloInput{ID: kolo.ID, SerijskaStevilka: kolo.SerijskaStevilka, Mnenje: kolo.Mnenje, JeIzposojen: kolo.JeIzposojen}
		updatedKolesaArray = append(updatedKolesaArray, &updatedKolo)
	}

	updatedKolesaArray = append(updatedKolesaArray, &novoKolo)

	updatedPostajalisce := model.UpdatePostajalisce{
		ID:          postajalisce.ID,
		Ime:         &postajalisce.Ime,
		Naslov:      &postajalisce.Naslov,
		Latitude:    &postajalisce.Latitude,
		Longitude:   &postajalisce.Longitude,
		KolesaArray: updatedKolesaArray,
	}

	db.UpdatePostajalisce(updatedPostajalisce)

	filter := bson.M{"_id": ObjectID}

	update := bson.M{"$set": bson.M{"end_date": formattedTime, "end_station_id": input.EndStationID, "trenutna_zasedenost_end": len(postajalisce.KolesaArray),
		"end_station": input.EndStation, "duration": trajanje}}

	_, err = izposojeColl.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Default().Println("ERROR -> couldn't return Kolo into Postajalisce")
		log.Fatal(err)
	}

	return db.FindIzposoja(input.ID)
}

// function for deleting an Izposoja from database
func (db *DB) DeleteIzposoja(id string) string {
	izposojeColl := db.client.Database(DATABASE).Collection(IZPOSOJE)

	ObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Default().Println("ERROR -> couldn't convert id to ObjectID")
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": ObjectID}

	_, err = izposojeColl.DeleteOne(ctx, filter)
	if err != nil {
		log.Default().Println("ERROR -> couldn't delete Izposoja from database")
		log.Fatal(err)
	}

	return "OK -> successfully deleted Izposoja from database"
}

// function for finding an Izposoja in database
func (db *DB) FindIzposoja(id string) *model.Izposoja {
	ObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Default().Println("ERROR -> couldn't convert id to ObjectID")
		log.Fatal(err)
	}

	izposojaColl := db.client.Database(DATABASE).Collection(IZPOSOJE)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := izposojaColl.FindOne(ctx, bson.M{"_id": ObjectID})

	izposoja := model.Izposoja{}

	err = result.Decode(&izposoja)

	if err != nil {
		log.Default().Println("ERROR -> couldn't decode result into Izposoja")
		log.Fatal(err)
	}

	return &izposoja
}

// function for finding all Izposoja in database
func (db *DB) FindAllIzposoja() []*model.Izposoja {
	izposojaColl := db.client.Database(DATABASE).Collection(IZPOSOJE)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := izposojaColl.Find(ctx, bson.D{})
	if err != nil {
		log.Default().Println("ERROR -> couldn't create cursor for Izposoja")
		log.Fatal(err)
	}

	var izposoje []*model.Izposoja
	for cursor.Next(ctx) {
		var izposoja *model.Izposoja
		err := cursor.Decode(&izposoja)

		if err != nil {
			log.Default().Println("ERROR -> couldn't decode result into Izposoja")
			log.Fatal(err)
		}

		izposoje = append(izposoje, izposoja)
	}

	return izposoje
}

// function for getting all Izposoja for a specific user
func (db *DB) FindAllIzposojaByUser(username string) []*model.Izposoja {
	izposojaColl := db.client.Database(DATABASE).Collection(IZPOSOJE)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := izposojaColl.Find(ctx, bson.D{})
	if err != nil {
		log.Default().Println("ERROR -> couldn't create cursor for Izposoja")
		log.Fatal(err)
	}

	var izposoje []*model.Izposoja
	for cursor.Next(ctx) {
		var izposoja *model.Izposoja
		err := cursor.Decode(&izposoja)

		if err != nil {
			log.Default().Println("ERROR -> couldn't decode result into Izposoja")
			log.Fatal(err)
		}

		if izposoja.Username == username {
			izposoje = append(izposoje, izposoja)
		}
	}

	return izposoje
}

// function for inserting Mnenje into Kolo
func (db *DB) InsertMnenje(_id string, mnenje int) string {
	var kolo = db.FindKolo(_id)

	kolo.Mnenje = append(kolo.Mnenje, &mnenje)

	var novoKolo = model.UpdateKolo{ID: kolo.ID, Mnenje: kolo.Mnenje}

	db.UpdateKolo(novoKolo)

	return "OK -> successfully inserted Mnenje into Kolo"
}

// function for updating the necessary Kolo in Postajalisce
func (db *DB) UpdateKoloInPostajalisce(id string, picikl *model.KoloInput) {
	postajalisceColl := db.client.Database(DATABASE).Collection(POSTAJALISCA)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := postajalisceColl.Find(ctx, bson.D{})
	if err != nil {
		log.Default().Println("ERROR -> couldn't create cursor for Postajalisce")
		log.Fatal(err)
	}

	var kolesa []*model.Kolo
	var postajalisceID = ""
	for cursor.Next(ctx) {
		var postajalisce *model.Postajalisce
		err := cursor.Decode(&postajalisce)

		if err != nil {
			log.Default().Println("ERROR -> couldn't decode result into Postajalisce")
			log.Fatal(err)
		}

		for _, kolo := range postajalisce.KolesaArray {
			if kolo.ID == id {
				postajalisceID = postajalisce.ID
				kolesa = postajalisce.KolesaArray
			}
		}
	}

	if postajalisceID != "" {
		var updatedKolesa []*model.KoloInput

		for _, kolo := range kolesa {
			if kolo.ID != id {
				var picikl = model.KoloInput{
					ID:               kolo.ID,
					SerijskaStevilka: kolo.SerijskaStevilka,
					Mnenje:           kolo.Mnenje,
					JeIzposojen:      kolo.JeIzposojen,
				}

				updatedKolesa = append(updatedKolesa, &picikl)
			}
		}

		updatedKolesa = append(updatedKolesa, picikl)

		updatedPostajalisce := model.UpdatePostajalisce{
			ID:          postajalisceID,
			KolesaArray: updatedKolesa,
		}

		db.UpdatePostajalisce(updatedPostajalisce)
	}
}

// function for updating postajalisce from admin console
func (db *DB) UpdatePostajalisceAdmin(postajalisce model.ChangePostajalisce) string {
	postaja := db.FindPostajalisce(postajalisce.ID)

	postajajalisce := model.UpdatePostajalisce{
		ID:        postaja.ID,
		Ime:       postajalisce.Ime,
		Naslov:    postajalisce.Naslov,
		Latitude:  postajalisce.Latitude,
		Longitude: postajalisce.Longitude,
	}

	db.UpdatePostajalisce(postajajalisce)

	return "OK"
}

// function for updating kolo from admin console
func (db *DB) UpdateKoloAdmin(kolo model.ChangeKolo) string {

	picikl := db.FindKolo(kolo.ID)

	koloUpdate := model.UpdateKolo{
		ID:               picikl.ID,
		SerijskaStevilka: kolo.SerijskaStevilka,
	}

	db.UpdateKolo(koloUpdate)

	return "OK"
}

// function for calculating distance from location to Postajalisce
func distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * lat1 / 180)
	radlat2 := float64(PI * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	return dist
}
