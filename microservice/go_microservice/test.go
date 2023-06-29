package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestReadPdf(t *testing.T) {
	// Создание временного файла с содержимым PDF
	tempFile, err := ioutil.TempFile("", "test_*.pdf")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Запись содержимого PDF во временный файл
	content := "This is a sample PDF content."
	err = ioutil.WriteFile(tempFile.Name(), []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write PDF content to file: %v", err)
	}

	// Вызов функции для чтения PDF
	result, err := readPdf(tempFile.Name())

	// Проверка результатов
	if err != nil {
		t.Errorf("Error reading PDF: %v", err)
	}
	if result != content {
		t.Errorf("Expected PDF content: %s, but got: %s", content, result)
	}
}

func TestReadTextFile(t *testing.T) {
	// Создание временного файла с текстовым содержимым
	tempFile, err := ioutil.TempFile("", "test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Запись текстового содержимого во временный файл
	content := "This is a sample text content."
	err = ioutil.WriteFile(tempFile.Name(), []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write text content to file: %v", err)
	}

	// Вызов функции для чтения текстового файла
	result, err := readTextFile(tempFile.Name())

	// Проверка результатов
	if err != nil {
		t.Errorf("Error reading text file: %v", err)
	}
	if result != content {
		t.Errorf("Expected text content: %s, but got: %s", content, result)
	}
}

// func readTextFile(s string) {
// 	panic("unimplemented")
// }

func TestSaveTextFile(t *testing.T) {
	// Создание временного файла
	tempFile, err := ioutil.TempFile("", "test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Запись текстового содержимого во временный файл
	content := "This is a sample text content."
	err = saveTextFile(tempFile.Name(), content)
	if err != nil {
		t.Fatalf("Failed to save text content to file: %v", err)
	}

	// Чтение содержимого временного файла
	result, err := ioutil.ReadFile(tempFile.Name())

	// Проверка результатов
	if err != nil {
		t.Errorf("Error reading saved file: %v", err)
	}
	if string(result) != content {
		t.Errorf("Expected saved content: %s, but got: %s", content, string(result))
	}
}
