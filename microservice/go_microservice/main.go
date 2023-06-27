package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/ledongthuc/pdf"
	"google.golang.org/grpc"

	cl "go_microservice/pkg/client"
)

var uploadedFilePath = "uploaded_file.txt"

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", fileUploadHandler).Methods("GET")
	r.HandleFunc("/", fileUploadHandler).Methods("POST")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Отобразить форму для загрузки файла
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<body>
				<form method="post" enctype="multipart/form-data">
					<input type="file" name="file">
					<input type="submit" value="Upload">
				</form>
			</body>
			</html>
		`)
		return
	}

	if r.Method == "POST" {
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
	}
}

func deletePreviousFile(filename string) error {
	err := os.Remove(filename)
	if err != nil && !os.IsNotExist(err) {
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

	// Вывести результаты анализа на веб-страницу
	fmt.Fprintf(w, "Analysis Result:\n")
	fmt.Fprintf(w, "Timestamp: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "Hard Reading: %d\n", result.GetHardReading())
	fmt.Fprintf(w, "Water Value: %d\n", result.GetWaterValue())
	fmt.Fprintf(w, "Mood: %s\n", result.GetMood())

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
