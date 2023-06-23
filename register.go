package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Объект для подключения к базе данных MySQL
var db *sql.DB

// Структура данных для хранения информации о пользователе
type User struct {
	Username string
	Password string
	Email    string
}

// Обработчик GET-запроса на главную страницу
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("register.html"))
	tmpl.Execute(w, nil)
}

// Обработчик POST-запроса при отправке формы регистрации
func registerHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем данные из формы
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	// Создаем новый объект пользователя
	user := User{
		Username: username,
		Password: password,
		Email:    email,
	}

	// Сохраняем пользователя в базе данных
	err := saveUser(user)
	if err != nil {
		log.Fatal(err)
	}

	// Отправляем ответ об успешной регистрации
	fmt.Fprintf(w, "Регистрация успешно завершена!")
}

// Функция для сохранения пользователя в базе данных
func saveUser(user User) error {
	// Подключаемся к базе данных MySQL
	db, err := sql.Open("mysql", "root:***@tcp(localhost:3308)/golang")
	if err != nil {
		return err
	}
	defer db.Close()

	// Проверяем наличие пользователя с таким же email в базе данных
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", user.Email).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		// Пользователь с таким email уже существует, возвращаем ошибку
		return errors.New("Пользователь с таким email уже зарегистрирован")
	}

	// Выполняем INSERT-запрос для добавления пользователя
	_, err = db.Exec("INSERT INTO users (username, password, email) VALUES (?, ?, ?)", user.Username, user.Password, user.Email)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Устанавливаем обработчики запросов
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/register", registerHandler)

	// Запускаем веб-сервер на порту 8000
	log.Fatal(http.ListenAndServe(":8000", nil))
}
