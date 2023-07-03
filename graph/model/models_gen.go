// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Dog struct {
	ID        string `json:"_id" bson:"_id"`
	Name      string `json:"name"`
	IsGoodBoi bool   `json:"isGoodBoi"`
}

type Kolo struct {
	ID               string    `json:"_id" bson:"_id"`
	SerijskaStevilka string    `json:"serijska_stevilka" bson:"serijska_stevilka"`
	Mnenje           []*string `json:"mnenje" bson:"mnenje"`
}

type NewDog struct {
	Name      string `json:"name"`
	IsGoodBoi bool   `json:"isGoodBoi"`
}

type NewKolo struct {
	SerijskaStevilka string `json:"serijska_stevilka"`
}

type NewPostajalisce struct {
	Ime       string  `json:"ime"`
	Naslov    string  `json:"naslov"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Postajalisce struct {
	ID          string  `json:"_id"`
	Ime         string  `json:"ime"`
	Naslov      string  `json:"naslov"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	KolesaArray []*Kolo `json:"kolesaArray"`
}
