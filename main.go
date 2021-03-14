package main

import (
	"database/sql"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Post struct {
	Id    int
	Title string
	Body  string
}

var db, err = sql.Open("mysql", "root:root@/golang?charset=utf8")

// CREATE TABLE POSTS (ID INT(6) AUTO_INCREMENT PRIMARY KEY, TITLE VARCHAR(50), BODY TEXT);

func main() {
	// stmt, err := db.Prepare("INSERT INTO POSTS(TITLE, BODY) VALUES (?,?)")
	// checkErr(err)

	// _, err = stmt.Exec("My First Post 2", "My first content 2")
	// checkErr(err)

	r := mux.NewRouter()
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static/"))))

	r.HandleFunc("/{id}/view", ViewHandler)
	r.HandleFunc("/", HomeHandler)

	http.ListenAndServe(":8080", r)
}

func ListPosts() []Post {
	rows, err := db.Query("SELECT * FROM POSTS")
	checkErr(err)

	items := []Post{}

	for rows.Next() {
		post := Post{}
		rows.Scan(&post.Id, &post.Title, &post.Body)
		items = append(items, post)
	}
	return items
}

func GetPostById(id string) Post {
	row := db.QueryRow("SELECT * FROM POSTS WHERE ID=?", id)
	post := Post{}
	row.Scan(&post.Id, &post.Title, &post.Body)
	return post
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/layout.html", "templates/list.html"))
	if err := t.ExecuteTemplate(w, "layout.html", ListPosts()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	t := template.Must(template.ParseFiles("templates/layout.html", "templates/view.html"))
	t.ExecuteTemplate(w, "layout.html", GetPostById(id))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
