package apis

/*
	apis (http API service)
	This component handles reciving and sending http signals.
	Data is transmitted in JSON in both directions.
	The endpoints are restricted, which is to say, only the necessary endpoints are exposed.
	These are:-
	common root /api/v1/
	/login for authentication
	/leads/add
	/leads/leadID
	Updates are handled via the submitted JSON by calling /leads/add with a specified id in the JSON
	Delete is handled via /leads/leadID endpoint where the http verb is DELETE

*/

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

/*
	Requires name and password in JSON call.
*/
func login(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var loginData map[string]interface{}
	json.Unmarshal([]byte(body), &loginData)
	name := loginData["name"].(string)
	pword := loginData["password"].(string)
	// TODO: Implement X-Public header value.
	valid := sqldb.Login(name, pword)

	w.Header().Set("Content-Type", "application/json")
	if valid == false {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("fobidden")
		return
	} else {
		json.NewEncoder(w).Encode("{\"status\":\"validated\"}")
	}
}

/*
	Checks http headers for X-Token, resonds with fobidden if not found or incorrect.
	Token encryption is handled by the sqldb component.
*/
func authenticate(w http.ResponseWriter, r *http.Request) bool {
	// pl := r.Header.Get("X-public")
	tk := r.Header.Get("X-Token")
	if tk != sqldb.Token || tk == "" || sqldb.Token == "" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("fobidden")
		return false
	}
	return true
}

/*
	DELETE handler for incomming delete signal.
*/
func deleteLead(w http.ResponseWriter, r *http.Request) {
	if authenticate(w, r) == false {
		return
	}
	pathParams := mux.Vars(r)
	leadid := pathParams["leadID"]
	var id int
	id, _ = strconv.Atoi(leadid)
	sqldb.DeleteLead(id)
}

/*
	POST Handler for add and update signals.
*/
func addLead(w http.ResponseWriter, r *http.Request) {
	if authenticate(w, r) == false {
		return
	}
	body, _ := ioutil.ReadAll(r.Body)

	var leadData map[string]interface{}

	json.Unmarshal([]byte(body), &leadData)

	// Detect if an id is specified in the JSON package, if so, assume update.
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
	// TODO: return confrimation
	json.NewEncoder(w).Encode(fmt.Sprintf("%v", mlead))
}

/*
	GET handler for all records or records by specified leadID.
*/
// TODO: add more search/response options eg: by first name.
func leads(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	if authenticate(w, r) == false {
		return
	}

	leadid := pathParams["leadID"]
	if leadid == "all" {
		allLeads(w)
	}

	if id, err := strconv.Atoi(leadid); err == nil {
		leadById(id, w)
	}
}

func leadById(id int, w http.ResponseWriter) {
	mleads := sqldb.LeadById(id)
	responseMaker(mleads, w)
}

func allLeads(w http.ResponseWriter) {
	mleads := sqldb.AllLeads()
	responseMaker(mleads, w)
}

/*
	Blank endpoints for development testing of server.
*/
func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Gets Still alive!")
}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Posts Still alive!")
}

func responseMaker(list []sqldb.Mlead, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	var tmp []string
	for _, v := range list {
		mlead := sqldb.Mlead{v.Id, v.FirstName, v.LastName, v.Email, v.Company, v.Postcode, v.AcceptTerms, v.DateCreated}
		b, _ := json.Marshal(mlead)
		tmp = append(tmp, string(b))
	}
	json.NewEncoder(w).Encode(tmp)
}

/*
	Vestigal function for the component as stand alone app.
*/
func main() {
	Routerer()
}
