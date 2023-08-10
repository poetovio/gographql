// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type ChangeKolo struct {
	ID               string  `json:"_id" bson:"_id"`
	SerijskaStevilka *string `json:"serijska_stevilka,omitempty" bson:"serijska_stevilka,omitempty"`
	Mnenje           []*int  `json:"mnenje,omitempty" bson:"mnenje,omitempty"`
}

type ChangePostajalisce struct {
	ID        string   `json:"_id" bson:"_id"`
	Ime       *string  `json:"ime,omitempty" bson:"ime,omitempty"`
	Naslov    *string  `json:"naslov,omitempty" bson:"naslov,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty" bson:"longitude,omitempty"`
}

type Izposoja struct {
	ID                      string  `json:"_id" bson:"_id"`
	StartDate               string  `json:"start_date" bson:"start_date"`
	EndDate                 *string `json:"end_date,omitempty" bson:"end_date,omitempty"`
	StartStationID          string  `json:"start_station_id" bson:"start_station_id"`
	EndStationID            *string `json:"end_station_id,omitempty" bson:"end_station_id,omitempty"`
	BikeID                  string  `json:"bike_id" bson:"bike_id"`
	Duration                *int    `json:"duration,omitempty" bson:"duration,omitempty"`
	TrenutnaZasedenostStart int     `json:"trenutna_zasedenost_start" bson:"trenutna_zasedenost_start"`
	TrenutnaZasedenostEnd   *int    `json:"trenutna_zasedenost_end,omitempty" bson:"trenutna_zasedenost_end,omitempty"`
	Weather                 string  `json:"weather" bson:"weather"`
	StartStation            string  `json:"start_station" bson:"start_station"`
	EndStation              string  `json:"end_station" bson:"end_station"`
	Username                string  `json:"username" bson:"username"`
}

type IzposojaKolesa struct {
	BikeID   string `json:"bike_id" bson:"bike_id"`
	Weather  string `json:"weather" bson:"weather"`
	Username string `json:"username" bson:"username"`
}

type Kolo struct {
	ID               string    `json:"_id" bson:"_id"`
	SerijskaStevilka string    `json:"serijska_stevilka" bson:"serijska_stevilka"`
	Mnenje           []*int `json:"mnenje" bson:"mnenje"`
	Rezervirano      bool      `json:"rezervirano" bson:"rezervirano"`
	JeIzposojen      bool      `json:"jeIzposojen" bson:"jeIzposojen"`
}

type KoloInput struct {
	ID               string `json:"_id" bson:"_id"`
	SerijskaStevilka string `json:"serijska_stevilka" bson:"serijska_stevilka"`
	Mnenje           []*int `json:"mnenje" bson:"mnenje"`
	JeIzposojen      bool   `json:"jeIzposojen" bson:"jeIzposojen"`
}

type NewKolo struct {
	SerijskaStevilka string `json:"serijska_stevilka"`
}

type NewPostajalisce struct {
	Ime       string  `json:"ime" bson:"ime"`
	Naslov    string  `json:"naslov" bson:"naslov"`
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

type Postajalisce struct {
	ID          string  `json:"_id" bson:"_id"`
	Ime         string  `json:"ime" bson:"ime"`
	Naslov      string  `json:"naslov" bson:"naslov"`
	Latitude    float64 `json:"latitude" bson:"latitude"`
	Longitude   float64 `json:"longitude" bson:"longitude"`
	KolesaArray []*Kolo `json:"kolesaArray" bson:"kolesaArray"`
}

type UpdateKolo struct {
	ID               string  `json:"_id" bson:"_id"`
	SerijskaStevilka *string `json:"serijska_stevilka,omitempty" bson:"serijska_stevilka,omitempty"`
	Mnenje           []*int  `json:"mnenje,omitempty" bson:"mnenje,omitempty"`
}

type UpdatePostajalisce struct {
	ID          string       `json:"_id" bson:"_id"`
	Ime         *string      `json:"ime,omitempty" bson:"ime,omitempty"`
	Naslov      *string      `json:"naslov,omitempty" bson:"naslov,omitempty"`
	Latitude    *float64     `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude   *float64     `json:"longitude,omitempty" bson:"longitude,omitempty"`
	KolesaArray []*KoloInput `json:"kolesaArray,omitempty" bson:"kolesaArray,omitempty"`
}

type VraciloKolesa struct {
	ID           string `json:"_id" bson:"_id"`
	BikeID       string `json:"bike_id" bson:"bike_id"`
	EndStationID string `json:"end_station_id" bson:"end_station_id"`
	EndStation   string `json:"end_station" bson:"end_station"`
}
