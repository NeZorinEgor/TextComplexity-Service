package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/ledongthuc/pdf"
)

var uploadedFilePath = "uploaded_file.txt"

func main() {
	http.HandleFunc("/", fileUploadHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
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

		// Вывести содержимое на веб-страницу

		fmt.Fprintf(w, "%s>", content)
	}
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
