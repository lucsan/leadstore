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
	corsHandler(api)
	api.HandleFunc("", get).Methods(http.MethodGet)
	api.HandleFunc("", post).Methods(http.MethodPost)
	api.HandleFunc("/leads/{leadID}", leads).Methods(http.MethodGet)
	api.HandleFunc("/leads/add", addLead).Methods(http.MethodPost)
	api.HandleFunc("/leads/{leadID}", deleteLead).Methods(http.MethodDelete)
	api.HandleFunc("/login", login).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":3000", router))
	fmt.Println("Running server!")
}

func corsHandler(api *mux.Router) {
	api.HandleFunc("/login", options).Methods(http.MethodOptions)
	api.HandleFunc("/leads/{leadID}", options).Methods(http.MethodOptions)
	api.HandleFunc("/leads/add", options).Methods(http.MethodOptions)
}

/*
	Logs the admin user and provides a secure token.
*/
func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("login called")
	body, _ := ioutil.ReadAll(r.Body)
	var loginData map[string]interface{}
	json.Unmarshal([]byte(body), &loginData)
	name := loginData["name"].(string)
	pword := loginData["password"].(string)
	// TODO: Implement X-Public header value for multiple clients (ie: public key per licence or company).
	valid := sqldb.Login(name, pword)

	fmt.Println("token", sqldb.Token)
	w.Header().Set("Content-Type", "application/json")

	if valid == false {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("fobidden")
		return
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "X-Token")
		w.Header().Set("X-Token", sqldb.Token)
		w.WriteHeader(http.StatusOK)
		fmt.Println("hg", w.Header().Get("X-Token"))
		json.NewEncoder(w).Encode("{\"status\":\"validated\"}")
	}
}

/*
	Checks http headers for X-Token, resonds with forbidden if not found or incorrect.
	Token encryption is handled by the sqldb component.
*/
func authenticate(w http.ResponseWriter, r *http.Request) bool {
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
	specifyHeaders(w)
	if authenticate(w, r) == false {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	pathParams := mux.Vars(r)
	leadid := pathParams["leadID"]
	var id int
	id, _ = strconv.Atoi(leadid)
	sqldb.DeleteLead(id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("{\"status\": \"deleted\"}")
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
	// TODO: add error handling to prevent server bawking if it recieves an incorrect signal (eg: ID instead of id)
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
	specifyHeaders(w)
	json.NewEncoder(w).Encode(fmt.Sprintf("%v", mlead))
}

/*
	GET handler for all records or records by specified leadID.
*/
// TODO: add more search/response options eg: by first name.
func leads(w http.ResponseWriter, r *http.Request) {
	fmt.Println("leads called")
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

/*
	options responder for CORS pre-filght check response.
*/
func options(w http.ResponseWriter, r *http.Request) {
	specifyHeaders(w)
	return
}

func responseMaker(list []sqldb.Mlead, w http.ResponseWriter) {
	specifyHeaders(w)
	var tmp []string
	for _, v := range list {
		mlead := sqldb.Mlead{v.Id, v.FirstName, v.LastName, v.Email, v.Company, v.Postcode, v.AcceptTerms, v.DateCreated}
		b, _ := json.Marshal(mlead)
		tmp = append(tmp, string(b))
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tmp)
}

/*
	Common headers for CORS and authentication tokens.
*/
func specifyHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Token")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Expose-Headers", "X-Token")
}

/*
	Vestigal function for the component as stand alone app.
*/
func main() {
	Routerer()
}
