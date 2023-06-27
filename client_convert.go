package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/ledongthuc/pdf"
)

func main() {
	http.HandleFunc("/", uploadHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
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
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Сохранить загруженный файл на диск
		tempFile, err := ioutil.TempFile("", "uploaded_file.*.pdf")
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

		// Обработать сохраненный файл с помощью функции readPdf
		content, err := readPdf(tempFile.Name())
		if err != nil {
			http.Error(w, "Failed to process PDF", http.StatusInternalServerError)
			return
		}

		// Сохранить содержимое PDF в текстовый файл
		err = saveTextFile("uploaded_file.txt", content)
		if err != nil {
			http.Error(w, "Failed to save content", http.StatusInternalServerError)
			return
		}

		// Вывести содержимое PDF на веб-страницу
		fmt.Fprintf(w, "PDF content saved to uploaded_file.txt<br>")
		fmt.Fprintf(w, "PDF content:<br><pre>%s</pre>", content)
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

func saveTextFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
