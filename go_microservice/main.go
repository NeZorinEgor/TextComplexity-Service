package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	cl "go_microservice/pkg/client"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type MyObject struct {
	happiness float64
	water     float64
}

// Структура статьи
type State struct {
	Id        uint16
	Title     string
	Full_text string
	Happiness float64
	Water     float64
}

var posts = []State{}
var showPost = State{}

// Временная замена grpc
func GetObject() MyObject {
	obj := MyObject{
		happiness: 0.8,
		water:     0.5,
	}
	return obj
}

// Начальная страница
func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	// Connect to DB
	db, err := sql.Open("mysql", "EGOR:EGOR@tcp(127.0.0.1:3305)/calendar")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Выборка данных
	res, err := db.Query("Select * from `states`")
	if err != nil {
		panic(err)
	}

	//Создание списка статей
	posts = []State{}
	for res.Next() {
		var post State
		err = res.Scan(&post.Id, &post.Title, &post.Full_text, &post.Happiness, &post.Water)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}

	t.ExecuteTemplate(w, "index", posts)

}

// Обработка передачи статьи
func saveArticle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	full_text := r.FormValue("full_text")
	var add_info = GetObject()
	happiness := add_info.happiness
	water := add_info.water

	if title == "" || full_text == "" {
		fmt.Fprintf(w, "Не все данные заполнены")
	} else {
		// Connect to DB
		db, err := sql.Open("mysql", "EGOR:EGOR@tcp(127.0.0.1:3305)/calendar")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		//Внесение данных в DB
		insert, err := db.Query(fmt.Sprintf("INSERT INTO `states` (`title`, `full_text`, `happines`, `water`) VALUES ('%s', '%s', '%f', '%f')", title, full_text, happiness, water))

		if err != nil {
			panic(err)
		}
		defer insert.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
}

// Отображение уникального поста
func show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	// Connect to DB
	db, err := sql.Open("mysql", "EGOR:EGOR@tcp(127.0.0.1:3305)/calendar")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Выборка данных
	res, err := db.Query(fmt.Sprintf("Select * From `states` WHERE `id` = '%s'", vars["id"]))
	if err != nil {
		panic(err)
	}

	var showPost = State{}
	for res.Next() {
		var post State
		err = res.Scan(&post.Id, &post.Title, &post.Full_text, &post.Happiness, &post.Water)
		if err != nil {
			panic(err)
		}
		showPost = post
	}
	t.ExecuteTemplate(w, "show", showPost)

}

func handleFunc() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/save_article", saveArticle).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")

	http.Handle("/", rtr)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	http.ListenAndServe(":8080", nil)
}

func main() {
	conn, err := grpc.Dial("188.168.25.28:21112", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := cl.NewTextAnalysServiceClient(conn)
	fmt.Println(time.Now().String())

	result, err := c.GetResult(context.Background(), &cl.SettingsTextPB{
		Text: "I love you",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(time.Now().String())
	fmt.Println(result.GetHardReading())
	fmt.Println(result.GetWaterValue())
	fmt.Println(result.GetMood())
	// handleFunc()
}
