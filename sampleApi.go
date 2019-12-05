// sampleApi project sampleApi.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Person struct {
	ID      uint   `json:"id"` // id is supported only when both the charectors are  capitalized
	Name    string `json:"name"`
	Address string `json:"address"`
}

var dba *gorm.DB

var err error

func init() {
	if dba, err = getDBConnection(); err != nil {
		panic(err)
	}

	if err != nil {
		fmt.Println("error")
	}
	// defer db.Close()
	dba.AutoMigrate(&Person{})
}
func getDBConnection() (db *gorm.DB, err error) {
	db, err = gorm.Open("sqlite3", "./gorm.db")
	return
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/personList", getPersons).Methods("GET")
	router.HandleFunc("/personList/{id}", getPerson).Methods("GET")
	router.HandleFunc("/personList", createPerson).Methods("POST")
	router.HandleFunc("/personList/{id}", updatePerson).Methods("PUT")
	router.HandleFunc("/personList/{id}", deletePerson).Methods("DELETE")
	http.ListenAndServe(":8082", router)
}

func getPersons(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	var users []Person
	dba.Find(&users)

	json.NewEncoder(w).Encode(users)

}
func getPerson(w http.ResponseWriter, r *http.Request) {
	personId := validateId(w, r)
	var person Person
	if err := dba.Where("id = ?", personId).First(&person).Error; err != nil {
		http.Error(w, "not found", 404)
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Context-tpe", "application/json")
	json.NewEncoder(w).Encode(person)
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	var person Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, "Invalid data provided", 401)
	}
	dba.Create(&person)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(person)
}

func updatePerson(w http.ResponseWriter, r *http.Request) {
	personId := validateId(w, r)
	var personToUpdate Person
	var person Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, "Invalid data provided", 401)
		return
	}

	if err := dba.Where("id = ?", personId).First(&personToUpdate).Error; err != nil {
		http.Error(w, "not found", 404)
		return
	}
	personToUpdate.Address = person.Address
	personToUpdate.Name = person.Name
	dba.Save(personToUpdate)
	json.NewEncoder(w).Encode(personToUpdate)
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	personId := validateId(w, r)
	var person Person
	dba.Where("id = ?", personId).Delete(&person)
	w.WriteHeader(http.StatusOK)
}

func validateId(w http.ResponseWriter, r *http.Request) (id string) {
	id = mux.Vars(r)["id"]

	if _, err := strconv.ParseUint(id, 8, 32); err != nil {
		http.Error(w, "Invalid Id Format detected ", 401)
		return
	}
	return
}
