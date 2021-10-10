/*
 * reference https://
 */
package main

import (
	"log"

	"github.com/fauna/faunadb-go/v4/faunadb"

)

// struct for key value pairs
type Pairs struct {
	sdk     *faunadb.FaunaClient
	//session *faunadb.FaunaClient
	db      string
	table   string
}

// record fields (table row)
type Broadcast struct {
	Universal string `fauna:"universal"`
	Name string `fauna:"name"`
}

func NewPairs(secret string, db string) *Pairs {
	return &Pairs{
		sdk:   faunadb.NewFaunaClient(secret),
		db:    db,
		table: "broadcasts",
	}
}

// create kv pair where ch id is used to reference name/token
func (p *Pairs) Insert(token string, channel string) {


	// Populate new record
	row := Broadcast{Name: token, Universal: channel}

	// Adding record as new table row
	_, err := p.sdk.Query(
		faunadb.Create(
			faunadb.Collection(p.table),
			faunadb.Obj{"data": row},
		))

	if err != nil {
		log.Printf("Insert %s table failed - %v", p.table, err)
	}

}
