// app.go

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"strconv"
	"time"
	"web-service/models"
	_ "web-service/models"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (a *App) Initialize() {

	// Using cloudsql-proxy will help you avoid white listing IPs and handling SSL.
	// In this case, it will also run inside your Go program and will not require
	// an additional process or container.

	// Connection String details:
	// * user       - the user created inside the DB. You can see more details on how to create it without password here:
	//                https://cloud.google.com/sql/docs/sql-proxy#flags
	// * project-id - your project id
	// * zone       - your general zone (us-central1/us-west1/etc)
	// * db-name    - the name of the database instance as it appears in the console
	var dbConnectionString = os.Getenv("DBConnectionString")

	var err error
	if a.DB, err = gorm.Open("mysql", dbConnectionString); err != nil {
		panic(err)
	}
	a.DB.AutoMigrate(&models.Feedback{})
	a.DB.AutoMigrate(&models.Company{})

	a.Router = mux.NewRouter()

	a.initializeRoutes()
}

// Put new routes here
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/company-signup", a.createCompany).Methods("POST")
	a.Router.HandleFunc("/feedback/", a.createFeedback).Methods("GET").Queries("company", "{company:.*}", "level", "{level:.*}")
	a.Router.HandleFunc("/company/", a.getCompany).Methods("GET").Queries("company", "{company:.*}")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getCompany(w http.ResponseWriter, r *http.Request) {
	var c models.Company

	vars := mux.Vars(r)
	companyName, found := vars["company"]
	if found == false {
		respondWithError(w, http.StatusBadRequest, "Invalid company name")
		return
	}
	c.CompanyName = companyName

	if err := c.GetCompany(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}

func (a *App) createFeedback(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var f models.Feedback

	vars := mux.Vars(r)
	company, found := vars["company"]
	if found == false {
		respondWithError(w, http.StatusBadRequest, "Invalid company name")
		return
	}

	level, err := strconv.Atoi(vars["level"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid happiness level")
		return
	}
	f.Company = company
	f.Happiness = level
	f.Timestamp = time.Now()

	fmt.Println("Created Feedback")

	if err := f.CreateFeedback(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, f)
}

func (a *App) createCompany(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var c models.Company
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := c.CreateCompany(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, c)
}

// GetCompanies gets all companies so their manager's emails can be used for the feedback email
func (a *App) GetCompanies() ([]models.Company, error) {
	var companies []models.Company
	if err := a.DB.Model(models.Company{}).Find(&companies).Error; err != nil {
		fmt.Println(err.Error())
		return companies, err
	}
	return companies, nil
}
