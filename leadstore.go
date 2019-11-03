package main

import (
	"fmt"
	"leadstore/apis"
	"leadstore/sqldb"
)

/*
	LeadStore is prototype customer details http api and sqlite database.
*/

const databaseName = "leadstore.db"
const publicKey = "qwerty"

/*
	Initalise the two main components.
*/
func main() {
	fmt.Println("LeadStore Runing")
	sqldb.Run(publicKey, databaseName)
	apis.Routerer()

}
