package main

import (
	"fmt"
	"leadstore/apis"
	"leadstore/sqldb"
)

const publicKey = "qwerty"

func AllUsers() {

}

func main() {
	fmt.Println("LeadStore Runing")
	sqldb.Run(publicKey)
	apis.Routerer()

}
