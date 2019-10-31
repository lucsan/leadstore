package apis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"leadstore/sqldb"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var p = fmt.Println

func Routerer() {
	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("", get).Methods(http.MethodGet)
	api.HandleFunc("", post).Methods(http.MethodPost)
	api.HandleFunc("/leads/{leadID}", leads).Methods(http.MethodGet)
	api.HandleFunc("/leads/add", addLead).Methods(http.MethodPost)
	api.HandleFunc("/leads/{leadID}", deleteLead).Methods(http.MethodDelete)
	api.HandleFunc("/login", login).Methods(http.MethodPost)

	fmt.Println("Running server!")
	log.Fatal(http.ListenAndServe(":3000", router))
}

type Test_struct struct {
	Name string
}

func login(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	a := r.Header.Get("X-public")
	p("p ", a)
	var loginData map[string]interface{}
	json.Unmarshal([]byte(body), &loginData)
	p("login data ", loginData)
	name := loginData["name"].(string)
	pword := loginData["password"].(string)

	p(name, pword)
	sqldb.Login(name, pword)

	//
}

func deleteLead(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	leadid := pathParams["leadID"]
	var id int
	id, _ = strconv.Atoi(leadid)
	sqldb.DeleteLead(id)
}

func addLead(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	var leadData map[string]interface{}

	json.Unmarshal([]byte(body), &leadData)

	var id int
	id = -1
	if leadData["id"] != nil {
		iid := leadData["id"].(string)
		id, _ = strconv.Atoi(iid)
	}

	first := leadData["first"].(string)
	last := leadData["last"].(string)
	email := leadData["email"].(string)
	company := leadData["company"].(string)
	postcode := leadData["postcode"].(string)
	var terms = false

	if leadData["terms"].(string) == "true" {
		terms = true
	}
	mlead := sqldb.Mlead{id, first, last, email, company, postcode, terms, ""}

	sqldb.AddLead(mlead)
	json.NewEncoder(w).Encode(fmt.Sprintf("recieved something %v", mlead))
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Gets Still alive!")
	json.NewEncoder(w).Encode(r.Method)
}

func post(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Posts Still alive!")
	json.NewEncoder(w).Encode(r.Method)
}

func leads(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	//json.NewEncoder(w).Encode("Posts Still alive!")
	//json.NewEncoder(w).Encode(r.Method)
	//json.NewEncoder(w).Encode(pathParams)
	leadid := pathParams["leadID"]
	if leadid == "all" {
		allLeads(w)
	}

	if id, err := strconv.Atoi(leadid); err == nil {
		leadById(id, w)
	}
}

func leadById(id int, w http.ResponseWriter) {
	fmt.Printf("user %d ", id)
	mleads := sqldb.LeadById(id)
	responseMaker(mleads, w)
}

func allLeads(w http.ResponseWriter) {
	//json.NewEncoder(w).Encode("getting all users")
	mleads := sqldb.AllLeads()
	responseMaker(mleads, w)
}

func responseMaker(list []sqldb.Mlead, w http.ResponseWriter) {
	var tmp []string
	for _, v := range list {
		mlead := sqldb.Mlead{v.Id, v.FirstName, v.LastName, v.Email, v.Company, v.Postcode, v.AcceptTerms, v.DateCreated}
		b, _ := json.Marshal(mlead)
		tmp = append(tmp, string(b))
	}
	json.NewEncoder(w).Encode(tmp)
}

func main() {
	Routerer()
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Still alive!")
}
