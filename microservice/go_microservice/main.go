package main

import (
	"bytes"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	cl "go_microservice/pkg/client"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/ledongthuc/pdf"
	"google.golang.org/grpc"
)

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
var values []interface{}
var emptyArray []int
var uploadedFilePath = "uploaded_file.txt"

func deletePreviousFile(filename string) error {
	err := os.Remove(filename)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func readPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	// Помните о закрытии файла
	defer f.Close()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}

func readTextFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func saveTextFile(filename, content string) error {
	err := ioutil.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}

func analyzeText(w http.ResponseWriter, text string) {
	conn, err := grpc.Dial("51.250.14.14:1111", grpc.WithInsecure())

	if err != nil {
		http.Error(w, "Failed to connect to analysis service", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	c := cl.NewTextAnalysServiceClient(conn)

	result, err := c.GetResult(context.Background(), &cl.SettingsTextPB{
		Text: text,
	})
	if err != nil {
		http.Error(w, "Failed to get analysis result", http.StatusInternalServerError)
		return
	}

	// values = append(values, result.GetHardReading())
	// values = append(values, result.GetWaterValue())
	// values = append(values, result.GetMood())
	// values[0] = result.HardReading
	// values[1] = result.WaterValue
	// values[2] = result.Mood

	values = append(values, result.GetHardReading())
	values = append(values, result.GetWaterValue())
	values = append(values, result.GetMood())
	if len(values) > 3 {
		values = values[len(values)-3:]
	}
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
		log.Println("Info", "DB OK")
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
	file, handler, err := r.FormFile("file")

	if err != nil {
		http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Получить расширение файла
	fileExtension := filepath.Ext(handler.Filename)

	// Удалить предыдущий файл, если существует
	deletePreviousFile(uploadedFilePath)

	// Сохранить загруженный файл на диск
	tempFile, err := ioutil.TempFile("", "uploaded_file_*"+fileExtension)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Обработать сохраненный файл
	var content string
	if fileExtension == ".pdf" {
		content, err = readPdf(tempFile.Name())
	} else {
		content, err = readTextFile(tempFile.Name())
	}
	if err != nil {
		http.Error(w, "Failed to process file", http.StatusInternalServerError)
		return
	}

	// Сохранить содержимое в текстовый файл
	err = saveTextFile(uploadedFilePath, content)
	if err != nil {
		http.Error(w, "Failed to save content", http.StatusInternalServerError)
		return
	}

	// Выполнить анализ текста
	analyzeText(w, content)

	// Connect to DB
	db, err := sql.Open("mysql", "EGOR:EGOR@tcp(127.0.0.1:3305)/calendar")

	if err != nil {
		// fmt.Println(err.Error())
		log.Println("Error", err.Error())
		panic(err)
	} else {
		log.Println("Info", "DB OK")
	}
	defer db.Close()

	//Внесение данных в DB
	insert, err := db.Query(fmt.Sprintf("INSERT INTO `states` (`title`, `reading`, `water`, `mood`) VALUES ('%s', '%d', '%d', '%s')", title, values[0], values[1], values[2]))

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	} else {
		log.Println("Info", "Внесли новую запись в БД")
	}
	defer insert.Close()

	var lastID int
	err = db.QueryRow("SELECT MAX(id) FROM `states`").Scan(&lastID)
	if err != nil {
		panic(err.Error())
	}

	log.Printf("Последний ID: %d", lastID)
	postID := lastID
	postIDString := strconv.Itoa(postID)
	redirectURL := "/post/" + postIDString

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)

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

	log.Printf("результат - %d, %d, %s", values[0], values[1], values[2])

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
	handleFunc()
}
