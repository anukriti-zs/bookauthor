package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

type Book struct {
	Book_Id     int    `json:"book_id"`
	Title       string `json:"title"`
	Author      Author `json:"author"`
	Publication string `json:"publication"`
	PublishDate string `json:"publishdate"`
}
type Author struct {
	Author_Id int    `json:"author_id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Dob       string `json:"dob"`
	PenName   string `json:"penname"`
}

var db *sql.DB

var err error

func main() {

	db, err = sql.Open("mysql", "root:Anukriti@1609@tcp(127.0.0.1:3306)/bookauthor")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBookById).Methods("GET")
	router.HandleFunc("/postsauthor", createPostAuthor).Methods("POST")
	//router.HandleFunc("/updatebook/{id}", updatePostBook).Methods("PUT")
	router.HandleFunc("/postsbook", createPostBook).Methods("POST")
	//router.HandleFunc("/updateauthor/{id}", updatePostAuthor).Methods("PUT")
	//router.HandleFunc("/deleteauthor/{id}", deleteAuthor).Methods("DELETE")
	//router.HandleFunc("/deletebook/{id}", deleteBook).Methods("DELETE")
	http.ListenAndServe(":8000", router)
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var books []Book
	//var authors []Author
	result, err := db.Query("SELECT b.book_id,b.title,b.publication,b.publishdate,a.author_id,a.fname,a.lname,a.dob,a.penname from book b INNER JOIN author a on b.author_id=a.author_id")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var book Book
		if err := result.Scan(&book.Book_Id, &book.Title, &book.Publication, &book.PublishDate, &book.Author.Author_Id, &book.Author.FirstName, &book.Author.LastName, &book.Author.Dob, &book.Author.PenName); err != nil {
			fmt.Println(err)
		}
		fmt.Println(book)
		books = append(books, book)
	}
	json.NewEncoder(w).Encode(books)

}
func getBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT b.book_id,b.title,b.publication,b.publishdate,a.author_id,a.fname,a.lname,a.dob,a.penname from book b INNER JOIN author a on b.author_id=a.author_id WHERE book_id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var book Book
	for result.Next() {
		err := result.Scan(&book.Book_Id, &book.Title, &book.Publication, &book.PublishDate, &book.Author.Author_Id, &book.Author.FirstName, &book.Author.LastName, &book.Author.Dob, &book.Author.PenName)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(book)
}
func createPostAuthor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stmt, err := db.Prepare("INSERT INTO author(fname,lname,dob,penname) VALUES(?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	fname := keyVal["fname"]
	lname := keyVal["lname"]
	dob := keyVal["dob"]
	penname := keyVal["penname"]
	_, err = stmt.Exec(fname, lname, dob, penname)
	if err != nil {
		panic(err.Error())
	}connect
	fmt.Fprintf(w, "New post was created")
}
func createPostBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stmt, err := db.Prepare("INSERT INTO book(book_id,title,author_id,publication,publishdate) VALUES(?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	book_id := keyVal["book_id"]
	title := keyVal["title"]
	author_id := keyVal["author_id"]
	publication := keyVal["publication"]
	publishdate := keyVal["publishdate"]

	_, err = stmt.Exec(book_id, title, author_id, publication, publishdate)
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "New post was created")
}

/*func updatePostAuthor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE author SET fname = ? WHERE author_id = ?")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newfname := keyVal["fname"]
	_, err = stmt.Exec(newfname, params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Post with ID = %s was updated", params["id"])
}
func updatePostBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE book SET title = ? WHERE book_id = ?")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newTitle := keyVal["title"]
	_, err = stmt.Exec(newTitle, params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Post with ID = %s was updated", params["id"])
}
func deleteAuthor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM author WHERE author_id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Post with ID = %s was deleted", params["id"])
}
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM book WHERE book_id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Post with ID = %s was deleted", params["id"])
}
*/
