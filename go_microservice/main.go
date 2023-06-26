package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type MyObject struct {
	reading uint16
	water   uint16
	mood    string
}

// Структура статьи
type State struct {
	Id      uint16
	Title   string
	Reading uint16
	Water   uint16
	Mood    string
}

var posts = []State{}
var showPost = State{}

// Временная замена grpc
func GetObject() MyObject {
	obj := MyObject{
		reading: 3,
		water:   5,
		mood:    "sad",
	}
	return obj
}

// Начальная страница
func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		// fmt.Fprintf(w, err.Error())
		// log.Println(w, err.Error())
		log.Println("Error", err.Error())

	}

	// Connect to DB
	db, err := sql.Open("mysql", "EGOR:EGOR@tcp(127.0.0.1:3305)/calendar")
	if err != nil {
		// fmt.Println(err.Error())
		log.Println("Error", err.Error())
		panic(err)
	} else {
		log.Println("Info", "ДБ ОК")
	}
	defer db.Close()

	// Выборка данных
	res, err := db.Query("Select * from `states`")
	if err != nil {
		// fmt.Println(err.Error())
		log.Println("Error", err.Error())
		panic(err)
	}

	//Создание списка статей
	posts = []State{}
	for res.Next() {
		var post State
		err = res.Scan(&post.Id, &post.Title, &post.Reading, &post.Water, &post.Mood)
		if err != nil {
			// fmt.Println(err.Error())
			log.Println("Error", err.Error())
			panic(err)
		}
		posts = append(posts, post)
	}

	t.ExecuteTemplate(w, "index", posts)

}

// Обработка передачи статьи
func saveArticle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	var add_info = GetObject()
	reading := add_info.reading
	water := add_info.water
	mood := add_info.mood

	if title == "" {
		fmt.Fprintf(w, "Не все данные заполнены")
	} else {
		// Connect to DB
		db, err := sql.Open("mysql", "EGOR:EGOR@tcp(127.0.0.1:3305)/calendar")
		if err != nil {
			// fmt.Println(err.Error())
			log.Println("Error", err.Error())
			panic(err)
		} else {
			log.Println("Info", "ДБ ОК")
		}
		defer db.Close()

		//Внесение данных в DB
		insert, err := db.Query(fmt.Sprintf("INSERT INTO `states` (`title`, `reading`, `water`, `mood`) VALUES ('%s', '%d', '%d', '%s')", title, reading, water, mood))

		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		} else {
			log.Println("Info", "Внесли новую запись в БД")
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
		// fmt.Fprintf(w, err.Error())
		log.Println("Error", err.Error())
	}

	// Connect to DB
	db, err := sql.Open("mysql", "EGOR:EGOR@tcp(127.0.0.1:3305)/calendar")
	if err != nil {
		// fmt.Println(err.Error())
		// fmt.Fprintf(w, err.Error())
		log.Println("Error", err.Error())
		panic(err)
	}
	defer db.Close()

	// Выборка данных
	res, err := db.Query(fmt.Sprintf("Select * From `states` WHERE `id` = '%s'", vars["id"]))
	if err != nil {
		// fmt.Println(err.Error())
		log.Println("Error", err.Error())
		panic(err)
	}

	var showPost = State{}
	for res.Next() {
		var post State
		err = res.Scan(&post.Id, &post.Title, &post.Reading, &post.Water, &post.Mood)
		if err != nil {
			// fmt.Println(err.Error())
			log.Println("Error", err.Error())
			panic(err)
		}
		showPost = post
	}

	t.ExecuteTemplate(w, "show", showPost)

}

func handleFunc() {
	addr := flag.String("addr", ":8080", "Сетевой адрес веб-сервера")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/save_article", saveArticle).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")

	http.Handle("/", rtr)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	http.ListenAndServe(*addr, nil)

	// Применяем созданные логгеры к нашему приложению.
	infoLog.Printf("Запуск сервера на %s", *addr)
	err := http.ListenAndServe(*addr, rtr)
	errorLog.Fatalf("Этот хост %s уже занят. Ошибка - %s ", *addr, err)
}

func main() {

	// fileContent, err := ioutil.ReadFile("./test_text/test5.txt")
	// if err != nil {
	// 	fmt.Println("Ошибка чтения файла:", err)
	// 	return
	// }

	// conn, err := grpc.Dial("10.128.0.13:1111", grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }
	// defer conn.Close()
	// c := cl.NewTextAnalysServiceClient(conn)
	// fmt.Println(time.Now().String())

	// result, err := c.GetResult(context.Background(), &cl.SettingsTextPB{
	// 	Text: string(fileContent),
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(time.Now().String())
	// fmt.Println(result.GetHardReading())
	// fmt.Println(result.GetWaterValue())
	// fmt.Println(result.GetMood())
	handleFunc()
}
