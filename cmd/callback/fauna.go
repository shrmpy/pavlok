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
	keyField = f.ObjKey("ref")
)

func NewPairs(secret string, db string) *Pairs {
	return &Pairs{
		sdk:   f.NewFaunaClient(secret),
		db:    db,
		table: "broadcasts",
		index: "token_by_uuid",
	}
}

// ch id is used to reference name/token
func (p *Pairs) Query(id string) string {

	val, err := p.queryRecord(id)
	if err != nil {
		return ""
	}

	var br Broadcast
	_ = val.At(dataField).Get(&br)
	return br.Name
}

// create kv pair where ch id is used to reference name/token
func (p *Pairs) Update(token string, channel string) {
	// to prevent duplicates, we do a delete first
	// (approximate update behavior)
	p.deleteRecord(channel)

        // Populate new record
        row := Broadcast{Name: token, Universal: channel}

        // Adding record as new table row
        _, err := p.sdk.Query(
                f.Create(
                        f.Collection(p.table),
                        f.Obj{"data": row},
                ))

        if err != nil {
                log.Printf("Insert %s table failed - %v", p.table, err)
        }
}

func (p *Pairs) deleteRecord(id string)  {
	log.Printf("Delete by channel-id - %s", id)
	val, err := p.queryRecord(id)
	if err != nil {
		log.Print("Delete no-op, nothing to see here")
		return
	}

	var rowid f.RefV
	val.At(keyField).Get(&rowid)
	p.sdk.Query(f.Delete(rowid))
}

func (p *Pairs) queryRecord(id string) (val f.Value, err error) {
	log.Printf("Query by channel-id - %s", id)

	// Search the token_by_uuid index for matching ID
	val, err = p.sdk.Query(
		f.Get(f.MatchTerm(f.Index(p.index), id)),
	)

	if err != nil {
		log.Printf("Query by failed - %s ", id)
		log.Print(err)
		return 
	}

	log.Printf("Datarow matched - %s", id)
	log.Print(val)
	return 
}

