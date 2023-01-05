package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Stock       int    `json:"stock"`
}

type Result struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_store")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/product", getAll).Methods("GET")
	router.HandleFunc("/product", create).Methods("POST")
	router.HandleFunc("/product/{id}", getDetail).Methods("GET")
	router.HandleFunc("/product/{id}", update).Methods("PUT")
	router.HandleFunc("/product/{id}", delete).Methods("DELETE")
	fmt.Println("server started at localhost:9000")
	http.ListenAndServe(":9000", router)
}

func getAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var posts []Product
	result, err := db.Query("SELECT * from tb_product")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var post Product
		err := result.Scan(&post.ID, &post.Name, &post.Description, &post.Stock)
		if err != nil {
			panic(err.Error())
		}
		posts = append(posts, post)
	}
	json.NewEncoder(w).Encode(posts)
}

func getDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	result, err := db.Query("SELECT * FROM tb_product WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var post Product
	for result.Next() {
		err := result.Scan(&post.ID, &post.Name, &post.Description, &post.Stock)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(post)
}

func create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stmt, err := db.Prepare("INSERT INTO tb_product(name, description, stock) VALUES(?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	var product Product
	input := make(map[string]interface{})
	json.Unmarshal(body, &product)
	name := input["name"]
	description := input["description"]
	stock := input["stock"]

	_, err = stmt.Exec(name, description, stock)
	if err != nil {
		panic(err.Error())
	}
	res := Result{Success: true, Message: "Success create product"}
	result, err := json.Marshal(res)
	w.Write(result)
}

func update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	params := vars["id"]
	stmt, err := db.Prepare("UPDATE tb_product SET name = ?, description = ?, stock = ? WHERE id = ?")
	if err != nil {
		fmt.Println("Error disini")
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	input := make(map[string]interface{})
	json.Unmarshal(body, &body)
	newName := input["name"]
	newDescription := input["description"]
	newSstock := input["stock"]

	_, err = stmt.Exec(newName, newDescription, newSstock, params)
	if err != nil {
		panic(err.Error())
	}

	res := Result{Success: true, Message: "Success update product"}
	result, err := json.Marshal(res)
	w.Write(result)
}

func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM tb_product WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	res := Result{Success: true, Message: "Success delete product"}
	result, err := json.Marshal(res)
	w.Write(result)
}
