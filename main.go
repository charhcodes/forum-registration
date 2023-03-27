// https://earthly.dev/blog/golang-sqlite/
// https://softchris.github.io/golang-book/05-misc/05-sqlite/

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

var err error
var tmpl *template.Template

type Users struct {
	mu sync.Mutex // sync access to database connection object so
	// only one database is accessed at a time
	db *sql.DB
}

const file string = "sqlite/users.sqlite"

// const create string = `
//   CREATE TABLE IF NOT EXISTS users (
// 	id INTEGER PRIMARY KEY AUTOINCREMENT,
//   	uname TEXT,
//   	email TEXT,
//   	pword TEXT
//   );`

func main() {
	// connect to database
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		fmt.Println("could not open database")
		panic(err.Error())
	}
	fmt.Println("--connection success--")
	defer db.Close()

	createData(nil, nil, db)
	readDb(db)

	pageHandlers()
}

func pageHandlers() {
	path := "templates"
	fs := http.FileServer(http.Dir(path))
	http.Handle("/templates/", http.StripPrefix("/templates/", fs))

	// datab, err := sql.Open("mysql", "root:password@tcp(localhost:8080)/testdb")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer datab.Close()

	http.HandleFunc("/", registerHandler)
	// http.HandleFunc("/registerauth", registerAuthHandler)
	fmt.Printf("Fetching server...")
	http.ListenAndServe("localhost:8080", nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****registerHandler running*****")
	tmpl.ExecuteTemplate(w, "register.html", nil)
}

// func registerAuthHandler(w http.ResponseWriter, r *http.Request) {
// 	/*
// 		1. check username criteria
// 		2. check password criteria
// 		3. check if username is already exists in database
// 		4. create bcrypt hash from password
// 		5. insert username and password hash in database
// 		(email validation will be in another video)
// 	*/

// 	fmt.Println("*****registerAuthHandler running*****")
// 	r.ParseForm()
// 	username := r.FormValue("username")

// 	// check username for only alphaNumeric characters
// 	var nameAlphaNumeric = true
// 	for _, char := range username {
// 		// func IsLetter(r rune) bool, func IsNumber(r rune) bool
// 		// if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
// 		if unicode.IsLetter(char) == false && unicode.IsNumber(char) == false {
// 			nameAlphaNumeric = false
// 		}
// 	}

// 	// check username length
// 	var nameLength bool
// 	if 5 <= len(username) && len(username) <= 50 {
// 		nameLength = true
// 	}

// 	// check password criteria
// 	password := r.FormValue("password")
// 	fmt.Println("password:", password, "\npswdLength:", len(password))

// 	// variables that must pass for password creation criteria
// 	var pswdLowercase, pswdUppercase, pswdNumber, pswdSpecial, pswdLength, pswdNoSpaces bool
// 	pswdNoSpaces = true
// 	for _, char := range password {
// 		switch {
// 		// func IsLower(r rune) bool
// 		case unicode.IsLower(char):
// 			pswdLowercase = true
// 		// func IsUpper(r rune) bool
// 		case unicode.IsUpper(char):
// 			pswdUppercase = true
// 		// func IsNumber(r rune) bool
// 		case unicode.IsNumber(char):
// 			pswdNumber = true
// 		// func IsPunct(r rune) bool, func IsSymbol(r rune) bool
// 		case unicode.IsPunct(char) || unicode.IsSymbol(char):
// 			pswdSpecial = true
// 		// func IsSpace(r rune) bool, type rune = int32
// 		case unicode.IsSpace(int32(char)):
// 			pswdNoSpaces = false
// 		}
// 	}
// 	if 11 < len(password) && len(password) < 60 {
// 		pswdLength = true
// 	}
// 	fmt.Println("pswdLowercase:", pswdLowercase, "\npswdUppercase:", pswdUppercase, "\npswdNumber:", pswdNumber, "\npswdSpecial:", pswdSpecial, "\npswdLength:", pswdLength, "\npswdNoSpaces:", pswdNoSpaces, "\nnameAlphaNumeric:", nameAlphaNumeric, "\nnameLength:", nameLength)
// 	if !pswdLowercase || !pswdUppercase || !pswdNumber || !pswdSpecial || !pswdLength || !pswdNoSpaces || !nameAlphaNumeric || !nameLength {
// 		tmpl.ExecuteTemplate(w, "register.html", "please check username and password criteria")
// 		return
// 	}

// 	// check if username already exists for availability
// 	// stmt := "SELECT UserID FROM bcrypt WHERE username = ?"
// 	// row := datab.QueryRow(stmt, username)
// 	// var uID string
// 	// err := row.Scan(&uID)
// 	// if err != sql.ErrNoRows {
// 	// 	fmt.Println("username already exists, err:", err)
// 	// 	tmpl.ExecuteTemplate(w, "register.html", "username already taken")
// 	// 	return
// 	// }

// 	// create hash from password
// 	var hash []byte

// 	// func GenerateFromPassword(password []byte, cost int) ([]byte, error)
// 	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		fmt.Println("bcrypt err:", err)
// 		tmpl.ExecuteTemplate(w, "register.html", "there was a problem registering account")
// 		return
// 	}
// 	fmt.Println("hash:", hash)
// 	fmt.Println("string(hash):", string(hash))

// 	// func (db *DB) Prepare(query string) (*Stmt, error)
// 	// var insertStmt *sql.Stmt
// 	// insertStmt, err = datab.Prepare("INSERT INTO bcrypt (Username, Hash) VALUES (?, ?);")
// 	// if err != nil {
// 	// 	fmt.Println("error preparing statement:", err)
// 	// 	tmpl.ExecuteTemplate(w, "register.html", "there was a problem registering account")
// 	// 	return
// 	// }
// 	// defer insertStmt.Close()
// 	// var result sql.Result

// 	// result, err = insertStmt.Exec(username, hash)
// 	// rowsAff, _ := result.RowsAffected()
// 	// lastIns, _ := result.LastInsertId()
// 	// fmt.Println("rowsAff:", rowsAff)
// 	// fmt.Println("lastIns:", lastIns)
// 	// fmt.Println("err:", err)

// 	if err != nil {
// 		fmt.Println("error inserting new user")
// 		tmpl.ExecuteTemplate(w, "register.html", "there was a problem registering account")
// 		return
// 	}
// 	fmt.Fprint(w, "congrats, your account has been successfully created")
// }

func createData(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	data, err := db.Prepare("INSERT INTO users(id, uname, email, pword) values(?,?,?,?)")
	if err != nil {
		fmt.Println("could not open database")
		panic(err.Error())
	}

	r.ParseForm()
	username := r.FormValue("username")
	useremail := r.FormValue("email")
	userpw := r.FormValue("password")

	res, err := data.Exec(nil, username, useremail, userpw)
	if err != nil {
		fmt.Println("could not insert new data")
		panic(err.Error())
	}
	affected, _ := res.RowsAffected()
	log.Printf("Affected rows %d", affected)
}

func readDb(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Println("could not read database")
		panic(err.Error())
	}

	// iterate over database
	for rows.Next() {
		var id int
		var uname string
		var email string
		var pword string

		err = rows.Scan(&id, &uname, &email, &pword)
		if err != nil {
			fmt.Println("could not read database 2")
			panic(err.Error())
		}
		fmt.Println(id)
		fmt.Println(uname)
		fmt.Println(email)
		fmt.Println(pword)
	}
}
