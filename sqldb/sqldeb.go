package sqldb

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const dbname = "./test.db"
const mltable = "marketingLeads"

var p = fmt.Println

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

func CheckAdminCreds(adminName, pword string) bool {
	database, _ := sql.Open("sqlite3", dbname)
	rows, _ := database.Query("SELECT id, name, password FROM admin WHERE pword = ?", pword)
	var nr = rows.Next()
	rows.Close()
	database.Close()
	if nr == false {
		return false
	}

	var id int
	var name string
	var password string
	rows.Scan(&id, &name, &password)
	if adminName != name {
		return false
	}
	return true
}

func Login(name, pword string) {

}

func AllLeads() []Mlead {
	database, _ := sql.Open("sqlite3", dbname)
	mleads := extractLeads(database)
	database.Close()
	return mleads
}

func LeadById(id int) []Mlead {
	database, _ := sql.Open("sqlite3", dbname)
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

func AddLead(newLead Mlead) {
	database, _ := sql.Open("sqlite3", dbname)
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
	database, _ := sql.Open("sqlite3", dbname)
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
	p("updating ", lid)
	statement, _ := database.Prepare("UPDATE marketingLeads SET firstname = ?, lastname = ?, email = ?, company = ?, postcode = ?, acceptterms = ?, created = ? WHERE id = ?")
	statement.Exec(v.FirstName, v.LastName, v.Email, v.Company, v.Postcode, v.AcceptTerms, v.DateCreated, lid)
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

func createStubLeads() []Mlead {
	var tmp []Mlead
	tmp = stackLeads(tmp, "Baz", "Wong", "baz.wong@email.co", "Wongo Ltd", "pc1", true, "29-10-2019")
	tmp = stackLeads(tmp, "Larry", "Kong", "larry.kong@email.co", "Wongo Ltd", "pc1", true, "29-10-2019")
	tmp = stackLeads(tmp, "Nimrod", "Peabody", "nimibo@email.co", "Qongo Ltd", "pc1", true, "29-10-2019")
	return tmp
}

func stackLeads(stack []Mlead, first, last, email, company, postcode string, accept bool, date string) []Mlead {
	stack = append(stack, Mlead{0, first, last, email, company, postcode, accept, date})
	return stack
}

func printData(database *sql.DB) {
	mleads := extractLeads(database)
	for _, v := range mleads {
		fmt.Printf("%v\n", v)
	}
}

// Creates the leads db if it does not exist.
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

func Run() {
	database, _ := sql.Open("sqlite3", dbname)
	mleads := createStubLeads()
	prepDatabase(database)
	insertLeads(mleads, database)
	printData(database)
	database.Close()
}

func main() {
	Run()
}