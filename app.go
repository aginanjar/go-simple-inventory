package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"

	itemModel "github.com/aginanjar/go-simple-inventory/model/item"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (app *App) Initalize() {
	var err error

	app.DB, err = sql.Open("sqlite3", "inventory.db")
	if err != nil {
		fmt.Println("err 1")
		log.Fatal(err)
	}

	// test connection
	ping := app.DB.Ping()
	if ping == nil {
		fmt.Println("DB is connected")
	}

	app.Router = mux.NewRouter()
	app.initializeRoutes()
}

func (app *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, app.Router))

}

func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Home")
	})

	app.Router.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This only test")
	})

	app.Router.HandleFunc("/items", app.getItems).Methods("GET")
	app.Router.HandleFunc("/item/create", app.createItem).Methods("POST")
	app.Router.HandleFunc("/item/{id:[0-9]+}", app.getItem).Methods("GET")
	app.Router.HandleFunc("/item/{id:[0-9]+}", app.updateItem).Methods("PUT")
	app.Router.HandleFunc("/item/{id:[0-9]+}", app.deleteItem).Methods("DELETE")
}

// Item
func (app *App) getItems(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	items, err := itemModel.GetItems(app.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, items)
}

func (app *App) createItem(w http.ResponseWriter, r *http.Request) {
	var i itemModel.Item
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&i); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := i.CreateItem(app.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, i)
}

func (app *App) getItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Item ID")
		return
	}

	i := itemModel.Item{ID: id}
	if err := i.GetItem(app.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Item not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, i)
}

func (app *App) updateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	var i itemModel.Item
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&i); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	i.ID = id

	if err := i.UpdateItem(app.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, i)
}

func (app *App) deleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Item ID")
		return
	}

	i := itemModel.Item{ID: id}
	if err := i.DeleteItem(app.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
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
