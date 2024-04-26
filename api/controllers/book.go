package controllers

import (
	"database/sql"
	"encoding/json"

	"io"
	"strconv"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/lazarok09/treinandosql/database"
	"github.com/lazarok09/treinandosql/helpers"
)

type BookBody struct {
	Name string `json:"name"`
}
type BookResponse struct {
	Response interface{} `json:"response"`
}
type Book struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	connection, err := database.Connect()
	if err != nil {
		w.Write([]byte("An error ocurred when connecting the database"))

	}
	defer connection.Close()
	var book BookBody
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("An error ocurred when reading the body database"))

	}

	if bookMarshallError := json.Unmarshal(requestBody, &book); bookMarshallError != nil {
		w.Write([]byte("An error ocurred when parsing the message"))
	}

	statment, err := connection.Prepare("INSERT INTO Book (name) VALUES (?)")

	if err != nil {
		w.Write([]byte("An error ocurred when creating the SQL statment"))
	}
	defer statment.Close()

	resultOfInsertOperation, err := statment.Exec(book.Name)

	if err != nil {
		w.Write([]byte("An error ocurred when creating a book"))
	}

	bookIdResult, err := resultOfInsertOperation.LastInsertId()
	if err != nil {
		w.Write([]byte("An error ocurred when getting rows affected"))
	}

	response := BookResponse{Response: bookIdResult}

	finalResponse, err := json.Marshal(response)
	if err != nil {
		w.Write([]byte("An error ocurred when getting the final response"))

	}
	w.WriteHeader(http.StatusCreated)
	w.Write(finalResponse)
}
func GetBooks(w http.ResponseWriter, r *http.Request) {
	connection, err := database.Connect()
	if err != nil {
		w.Write([]byte("An error ocurred when connecting the database"))

	}
	defer connection.Close()
	rows, err := connection.Query("SELECT * FROM Book")

	if err != nil {
		w.Write([]byte("An error ocurred when scaning book"))

	}
	defer rows.Close()

	var books []Book

	for rows.Next() {
		var book Book
		err := rows.Scan(
			&book.ID,
			&book.Name,
		)

		if err != nil {
			w.Write([]byte("An error ocurred when scaning books"))

		}
		books = append(books, book)
	}

	response := BookResponse{Response: books}

	finalResponse, err := json.Marshal(response)
	if err != nil {
		w.Write([]byte("An error ocurred when getting the final response"))

	}

	w.Write(finalResponse)
}
func GetBook(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	paramName := "id"
	bookId, err := strconv.ParseUint(params[paramName], 10, 32)

	if err != nil {
		helpers.ThrowParamMissing(w, paramName)
		return
	}

	connection, err := database.Connect()
	if err != nil {
		helpers.ThrowDBConnectionError(w, err)
		return
	}
	defer connection.Close()

	var book Book
	queryRowError := connection.QueryRow("SELECT * FROM Book WHERE id = ?", bookId).Scan(&book.ID, &book.Name)
	if queryRowError != nil {
		if queryRowError == sql.ErrNoRows {
			helpers.ThrowEntityNotFounded(queryRowError.Error(), w, bookId)
			return
		}
		status := http.StatusInternalServerError
		w.WriteHeader(status)
		response := helpers.ResponseErrorShape{Message: "An error occurred when scanning book value to the struct", Error: queryRowError.Error(), Status: status}
		json.NewEncoder(w).Encode(response)
		return
	}

	json.NewEncoder(w).Encode(book)
}
func UpdateBook(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	paramName := "id"
	bookId, err := strconv.ParseUint(params[paramName], 10, 32)

	if err != nil {
		helpers.ThrowParamMissing(w, paramName)
		return
	}

	connection, err := database.Connect()
	if err != nil {
		helpers.ThrowDBConnectionError(w, err)
		return
	}
	defer connection.Close()

	var book Book
	queryRowError := connection.QueryRow("SELECT * FROM Book WHERE id = ?", bookId).Scan(&book.ID, &book.Name)
	if queryRowError != nil {
		if queryRowError == sql.ErrNoRows {
			helpers.ThrowEntityNotFounded(queryRowError.Error(), w, bookId)
			return
		}

		status := http.StatusInternalServerError
		w.WriteHeader(status)
		response := helpers.ResponseErrorShape{Message: "An error occurred when scanning book value to the struct", Error: queryRowError.Error(), Status: status}
		json.NewEncoder(w).Encode(response)
		return
	}

	json.NewEncoder(w).Encode(book)
}
