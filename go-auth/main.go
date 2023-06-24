package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// User представляет модель данных пользователя
type User struct {
	ID       int
	Username string
	Password []byte
}

// Создайте подключение к базе данных MySQL
func dbConn() (db *sql.DB) {

	db, err := sql.Open("mysql", "root:***@tcp(127.0.0.1:3308)/golang")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// Регистрация пользователя
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("register.html"))
		tmpl.Execute(w, nil)
	} else {
		db := dbConn()
		defer db.Close()

		username := r.FormValue("username")
		password := r.FormValue("password")

		// Хэширование пароля
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
			return
		}

		// Вставка пользователя в базу данных
		insertUserQuery := "INSERT INTO users(username, password) VALUES(?, ?)"
		_, err = db.Exec(insertUserQuery, username, hashedPassword)
		if err != nil {
			log.Fatal(err)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// Вход пользователя
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("login.html"))
		tmpl.Execute(w, nil)
	} else {
		db := dbConn()
		defer db.Close()

		username := r.FormValue("username")
		password := r.FormValue("password")

		// Получение хэшированного пароля из базы данных
		getUserQuery := "SELECT id, password FROM users WHERE username = ?"
		row := db.QueryRow(getUserQuery, username)

		var userID int
		var hashedPassword []byte
		err := row.Scan(&userID, &hashedPassword)
		if err != nil {
			log.Fatal(err)
			return
		}

		// Сравнение хэшированного пароля
		err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
		if err != nil {
			log.Fatal(err)
			return
		}

		// Вход выполнен успешно, перенаправление на index.html
		http.Redirect(w, r, "/index.html", http.StatusSeeOther)

	}
}

func main() {
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.Handle("/", http.FileServer(http.Dir("template")))
	fmt.Println("Server started on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
