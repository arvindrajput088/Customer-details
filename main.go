package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type customer struct {
	ID           string `json:"Id"`
	Name         string `json:"Name"`
	Phone_number int    `json:"Phone_nmuber"`
	Email        string `json:"Email"`
}

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "oorja173"
	DB_NAME     = "Customer"
)

// DB set up
func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	return db
}

// Function for handling errors
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	var newCustomer customer
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the customer name,phone number and email only in order to create")
	}
	checkErr(err)
	json.Unmarshal(reqBody, &newCustomer)
	//fmt.Println(reqBody)
	//fmt.Println(newCustomer.ID)
	db := setupDB()
	db.QueryRow(`INSERT INTO "Customer_details"("ID","Name","Phone_number","Email") VALUES($1, $2, $3, $4)`, newCustomer.ID, newCustomer.Name, newCustomer.Phone_number, newCustomer.Email)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newCustomer)
}
func getOneCustomer(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	customerID := mux.Vars(r)["id"]
	row := db.QueryRow(`SELECT * FROM "Customer_details" WHERE "ID" = $1`, customerID)

	var id string
	var Name string
	var Phone_number int
	var Email string

	err := row.Scan(&id, &Name, &Phone_number, &Email)

	var customers = customer{id, Name, Phone_number, Email}

	// check errors
	checkErr(err)

	json.NewEncoder(w).Encode(customers)
}

func getAllCustomers(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	rows, err := db.Query(`SELECT * FROM "Customer_details"`)
	checkErr(err)
	var customers []customer

	// Foreach movie
	for rows.Next() {
		var id string
		var Name string
		var Phone_number int
		var Email string
		err = rows.Scan(&id, &Name, &Phone_number, &Email)
		// check errors
		checkErr(err)
		customers = append(customers, customer{ID: id, Name: Name, Phone_number: Phone_number, Email: Email})
	}
	json.NewEncoder(w).Encode(customers)
}
func updateCustomer(w http.ResponseWriter, r *http.Request) {
	customerID := mux.Vars(r)["id"]
	db := setupDB()
	var updatedValues customer
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the customer name,phone number and email only in order to update")
	}
	json.Unmarshal(reqBody, &updatedValues)
	var id string
	var Name string
	var Phone_number int
	var Email string
	db.QueryRow(`UPDATE "Customer_details" SET "Name"=$1,"Phone_number"=$2,"Email"=$3 WHERE "ID"=$4 returning *;`, updatedValues.Name, updatedValues.Phone_number, updatedValues.Email, customerID).Scan(&id, &Name, &Phone_number, &Email)
	updatedCustomer := customer{id, Name, Phone_number, Email}
	json.NewEncoder(w).Encode(updatedCustomer)
}
func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	customerID := mux.Vars(r)["id"]
	db := setupDB()
	_, err := db.Exec(`DELETE FROM "Customer_details" where "ID" = $1`, customerID)

	// check errors
	checkErr(err)
}
func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/customer", createCustomer).Methods("POST")
	router.HandleFunc("/customers", getAllCustomers).Methods("GET")
	router.HandleFunc("/customer/{id}", getOneCustomer).Methods("GET")
	router.HandleFunc("/customers/{id}", updateCustomer).Methods("PATCH")
	router.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")
	fmt.Println("Server at 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
