package main

import (
	"fmt"
	"leadstore/sqldb"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const dbname = "test.db"

type Mlead struct {
	Id          int
	FirstName   string
	LastName    string
	Email       string
	Company     string
	Postcode    string
	AcceptTerms bool
	DateCreated string
}

func TestTest(t *testing.T) {
	// t.Errorf("Test %s tested", "test")
}

func TestRemoveDbFile(t *testing.T) {
	err := os.Remove(dbname)
	if err != nil {
		fmt.Println("delete error")
		fmt.Println(err)
		return
	}
}

func TestDatabase(t *testing.T) {
	sqldb.Run("123abc", dbname)
}

func TestStubLeads(t *testing.T) {
	mleads := sqldb.CreateStubLeads()

	if mleads[2].FirstName == "Ming" {
		t.Errorf("%s does not equal Ming", mleads[2].FirstName)
	}
}

func TestAllLeads(t *testing.T) {
	mleads := sqldb.AllLeads()
	if mleads[1].FirstName != "Ming" {
		t.Errorf("%s does not equal Ming", mleads[1].FirstName)
	}
}

func TestAdd(t *testing.T) {
	newLead := sqldb.Mlead{-1, "Sauroman", "White", "saucy@email.co", "Rivendale Co", "RI1", true, ""}
	sqldb.AddLead(newLead)
	mleads := sqldb.AllLeads()
	if mleads[3].LastName != "White" {
		t.Errorf("%s does not equal White", mleads[3].LastName)
	}
}

func TestUpdate(t *testing.T) {
	updatedLead := sqldb.Mlead{4, "Sauroman", "Black", "saucy@email.co", "Rivendale Co", "RI1", false, ""}
	sqldb.AddLead(updatedLead)
	mleads := sqldb.AllLeads()
	if mleads[3].LastName != "Black" {
		t.Errorf("%s does not equal Black", mleads[3].LastName)
	}
}

func TestDelete(t *testing.T) {
	sqldb.DeleteLead(4)
	mleads := sqldb.AllLeads()
	if len(mleads) > 3 {
		t.Errorf("%d is greater than 3", len(mleads))
	}
}
