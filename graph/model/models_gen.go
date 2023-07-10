// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Kolo struct {
	ID               string    `json:"_id" bson:"_id"`
	SerijskaStevilka string    `json:"serijska_stevilka" bson:"serijska_stevilka"`
	Mnenje           []*string `json:"mnenje" bson:"mnenje"`
	Rezervirano      bool      `json:"rezervirano" bson:"rezervirano"`
}

type NewKolo struct {
	SerijskaStevilka string `json:"serijska_stevilka" bson:"serijska_stevilka"`
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
	ID               string `json:"_id" bson:"_id"`
	SerijskaStevilka string `json:"serijska_stevilka" bson:"serijska_stevilka"`
}

type UpdatePostajalisce struct {
	ID        string   `json:"_id" bson:"_id"`
	Ime       *string  `json:"ime,omitempty" bson:"ime,omitempty"`
	Naslov    *string  `json:"naslov,omitempty" bson:"naslov,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty" bson:"longitude,omitempty"`
}
