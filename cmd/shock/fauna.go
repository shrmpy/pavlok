/*
 * reference https://
 */
package main

import (
	"log"

	f "github.com/fauna/faunadb-go/v4/faunadb"
)

// struct for key value pairs
type Pairs struct {
	sdk     *f.FaunaClient
	db      string
	table   string
	index   string
}

// record fields (table row)
type Broadcast struct {
	Universal string `fauna:"universal"`
	Name string `fauna:"name"`
}

var (
	dataField = f.ObjKey("data")
)

func NewPairs(secret string, db string) *Pairs {
	return &Pairs{
		sdk:   f.NewFaunaClient(secret),
		db:    db,
		table: "broadcasts",
		index: "token_by_uuid",
	}
}

// UUID is used to reference name/token
func (p *Pairs) Query(id string) string {
	log.Printf("Query by uuid - %s", id)

	// Search the token_by_uuid index for matching ID
	value, err := p.sdk.Query(
		f.Get(f.MatchTerm(f.Index(p.index), id)),
	)

	if err != nil {
		log.Printf("Query by failed - %s ", id)
		log.Print(err)
		return ""
	}

	log.Printf("Datarow matched - %s", id)
	log.Print(value)

	var sess Broadcast
	_ = value.At(dataField).Get(&sess)
	return sess.Name
}
