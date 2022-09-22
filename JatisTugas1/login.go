package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"log"
	"net/http"
	"text/template"
)

type User struct {
	id       int
	username string
	password string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "login"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("templates/*"))

func main() {
	log.Println("Server started on: http://localhost:9090")
	http.HandleFunc("/", Index)
	http.HandleFunc("/login", Login)
	http.ListenAndServe(":9090", nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "Index", nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	p := r.Form.Get("password")

	password := GetMD5Hash(p)
	db := dbConn()
	sqlStmt := "SELECT * FROM User WHERE username=? AND password=?"

	var id int
	row := db.QueryRow(sqlStmt, username, password)

	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		tmpl.ExecuteTemplate(w, "Index", nil)
	default:
		tmpl.ExecuteTemplate(w, "Success", nil)
	}

	defer db.Close()
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}