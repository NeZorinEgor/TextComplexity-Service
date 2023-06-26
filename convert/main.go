package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var uploadedFilePath = "uploaded_file.txt"

func main() {
	http.HandleFunc("/", fileUploadHandler)
	http.ListenAndServe(":8080", nil)
}

func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" && strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		content, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		text := strings.ReplaceAll(string(content), "\n", " ")

		f, err := os.Create(uploadedFilePath)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		writer := bufio.NewWriter(f)
		_, err = writer.WriteString(text)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = writer.Flush()
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("File uploaded and saved successfully!"))
	} else {
		http.ServeFile(w, r, "upload.html")
	}
}

func init() {
	// Удаляем предыдущий загруженный файл при запуске сервера
	err := os.Remove(uploadedFilePath)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error removing previous uploaded file:", err.Error())
	}
}
