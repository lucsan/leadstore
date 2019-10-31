package main

import (
	"fmt"
	"leadstore/apis"
	"leadstore/sqldb"
)

func AllUsers() {

}

func main() {
	fmt.Println("LeadStore Runing")
	sqldb.Run()
	apis.Routerer()

}
