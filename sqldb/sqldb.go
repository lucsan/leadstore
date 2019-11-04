package sqldb

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const mltable = "marketingLeads"

var dbname string

func dbase() *sql.DB {
	database, _ := sql.Open("sqlite3", dbname)
	return database
	//
}

/*
	Authentication values.
*/
var publicKey string
var Token string

/*
	Customer data model record/lead structure
*/
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

/*
	Wrapper function for checkAdminCreds.
*/
func Login(name, pword string) bool {
	return checkAdminCreds(name, pword)
}

/*
	validate the admin credentials, name and password.
*/
func checkAdminCreds(adminName, pword string) bool {
	database := dbase()
	rows, _ := database.Query("SELECT id, name, password FROM admin WHERE name = ?", adminName)
	var nr = rows.Next()

	if nr == false {
		rows.Close()
		database.Close()
		return false
	}

	var id int
	var name string
	var password string
	rows.Scan(&id, &name, &password)
	if password != pword {
		rows.Close()
		database.Close()
		return false
	}

	rows.Close()
	database.Close()
	tokenEncryption(pword)
	return true
}

/*
	Placeholder function for authentication.
	Naturally full token encryption would be inserted here.
*/
func tokenEncryption(pword string) {
	// TODO: implement full token encryption.
	Token = pword + publicKey
}

func tokenChecker(token string) bool {
	if Token != token {
		return false
	}
	return true
}

func AllLeads() []Mlead {
	database := dbase()
	mleads := extractLeads(database)
	database.Close()
	fmt.Println(mleads)
	return mleads
}

/*
	Retrieves lead by id.
*/
func LeadById(id int) []Mlead {
	database := dbase()
	var mleads []Mlead
	var firstname string
	var lastname string
	var email string
	var company string
	var postcode string
	var acceptterms bool
	var created string

	rows, _ := database.Query("SELECT id, firstname, lastname, email, company, postcode, acceptterms, created FROM marketingLeads WHERE id = ?", id)
	for rows.Next() {
		rows.Scan(&id, &firstname, &lastname, &email, &company, &postcode, &acceptterms, &created)
		mleads = append(mleads, Mlead{id, firstname, lastname, email, company, postcode, acceptterms, created})
	}
	rows.Close()
	database.Close()
	return mleads
}

/*
	Retrieves all leads.
*/
func extractLeads(database *sql.DB) []Mlead {
	var mleads []Mlead
	var id int
	var firstname string
	var lastname string
	var email string
	var company string
	var postcode string
	var acceptterms bool
	var created string

	rows, _ := database.Query("SELECT id, firstname, lastname, email, company, postcode, acceptterms, created FROM marketingLeads")
	for rows.Next() {
		rows.Scan(&id, &firstname, &lastname, &email, &company, &postcode, &acceptterms, &created)
		mleads = append(mleads, Mlead{id, firstname, lastname, email, company, postcode, acceptterms, created})
	}
	rows.Close()
	return mleads
}

/*
	Inserts leads without id's, updates one's with ids.
*/
func AddLead(newLead Mlead) {
	database := dbase()
	if newLead.Id > -1 {
		updateLead(newLead, database)
	} else {
		var mleads []Mlead
		mleads = append(mleads, newLead)
		insertLeads(mleads, database)
	}
	printData(database)
	database.Close()
}

func DeleteLead(id int) {
	database := dbase()
	statement, _ := database.Prepare("DELETE FROM marketingLeads WHERE id = ?")
	statement.Exec(id)
	printData(database)
	database.Close()
}

func insertLeads(mleads []Mlead, database *sql.DB) {
	statement, _ := database.Prepare("INSERT INTO marketingLeads (firstname, lastname, email, company, postcode, acceptterms, created) VALUES (?, ?, ?, ?, ?, ?, ?)")
	dateCreated := time.Now().Format("2006-01-02 15:04:05")
	for _, v := range mleads {
		check := checkLeadNotExists(database, v.FirstName, v.LastName, v.Email)
		if check == true {
			statement.Exec(v.FirstName, v.LastName, v.Email, v.Company, v.Postcode, v.AcceptTerms, dateCreated)
		}
	}
}

func updateLead(v Mlead, database *sql.DB) {
	lid := v.Id
	statement, _ := database.Prepare("UPDATE marketingLeads SET firstname = ?, lastname = ?, email = ?, company = ?, postcode = ?, acceptterms = ? WHERE id = ?")
	statement.Exec(v.FirstName, v.LastName, v.Email, v.Company, v.Postcode, v.AcceptTerms, lid)
	printData(database)
}

func checkLeadNotExists(database *sql.DB, first, last, email string) bool {
	rows, _ := database.Query("SELECT id, firstname, lastname, email FROM marketingLeads WHERE firstname = ? AND lastname = ? AND email = ?", first, last, email)
	if rows.Next() == false {
		rows.Close()
		return true
	}
	rows.Close()
	return false
}

/*
	Stub customer/lead data to pre-poulate db with for test and dev.
*/
func CreateStubLeads() []Mlead {
	var tmp []Mlead
	tmp = stackLeads(tmp, "Baz", "Wong", "baz.wong@email.co", "Wongo Ltd", "PC1", true)
	tmp = stackLeads(tmp, "Ming", "Merciless", "ming@email.co", "Monogo Plc", "PC2", false)
	tmp = stackLeads(tmp, "Nimrod", "Peabody", "nimibo@email.co", "Qongo Ltd", "PC3", true)
	return tmp
}

func stackLeads(stack []Mlead, first, last, email, company, postcode string, accept bool) []Mlead {
	stack = append(stack, Mlead{0, first, last, email, company, postcode, accept, ""})
	return stack
}

/*
	Test/dev output function for visual confirmation of db activity.
*/
func printData(database *sql.DB) {
	mleads := extractLeads(database)
	for _, v := range mleads {
		fmt.Printf("%v\n", v)
	}
}

/*
	Creates the leads db if it does not exist.
*/
func prepDatabase(database *sql.DB) {
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS marketingLeads (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT, email TEXT, company TEXT, postcode TEXT, acceptterms INTEGER, created TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS admin (id INTEGER PRIMARY KEY, name TEXT, password TEXT)")
	statement.Exec()
	var adminName = "admin"
	var pword = "passy"
	rows, _ := database.Query("SELECT id, name, password FROM admin WHERE name = ?", adminName)

	if rows.Next() == false {
		statement, _ := database.Prepare("INSERT INTO admin (name, password) VALUES (? ,?)")
		statement.Exec(adminName, pword)
	}
	rows.Close()
}

/*
	Recieves the public key for use in authentication.
	insertLeads and printData are test/dev function to be removed prior to production.
*/
func Run(pk, dbn string) {
	publicKey = pk
	dbname = dbn
	database := dbase()
	mleads := CreateStubLeads()
	prepDatabase(database)
	insertLeads(mleads, database)
	printData(database)
	database.Close()
}

func RunTest(dbn string) {
	dbname = dbn
}
