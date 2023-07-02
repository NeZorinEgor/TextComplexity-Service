package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestIndexHandler(t *testing.T) {
	// Создаем новый HTTP-запрос с методом GET и пустым телом
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем ResponseRecorder (реализация http.ResponseWriter) для записи ответа
	rr := httptest.NewRecorder()

	// Создаем маршрутизатор и регистрируем обработчик для "/"
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)

	// Выполняем запрос, передавая RequestRecorder и Request
	r.ServeHTTP(rr, req)

	// Проверяем код состояния ответа
	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался статус код %d, получен %d", http.StatusOK, rr.Code)
	}

	// Проверяем ожидаемый контент ответа
	expectedResponse := "Hello, World!"
	if rr.Body.String() != expectedResponse {
		t.Errorf("Ожидался ответ %s, получен %s", expectedResponse, rr.Body.String())
	}
}

func TestSaveArticleHandler(t *testing.T) {
	// Создаем новый HTTP-запрос с методом POST и пустым телом
	req, err := http.NewRequest("POST", "/save_article", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем ResponseRecorder (реализация http.ResponseWriter) для записи ответа
	rr := httptest.NewRecorder()

	// Создаем маршрутизатор и регистрируем обработчик для "/save_article"
	r := mux.NewRouter()
	r.HandleFunc("/save_article", saveArticleHandler)

	// Выполняем запрос, передавая RequestRecorder и Request
	r.ServeHTTP(rr, req)

	// Проверяем код состояния ответа
	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался статус код %d, получен %d", http.StatusOK, rr.Code)
	}

	// Проверяем ожидаемый контент ответа
	expectedResponse := "Article saved successfully"
	if rr.Body.String() != expectedResponse {
		t.Errorf("Ожидался ответ %s, получен %s", expectedResponse, rr.Body.String())
	}
}

// Функция-обработчик для главной страницы
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Отправляем "Hello, World!" в качестве ответа
	fmt.Fprint(w, "Hello, World!")
}

// Функция-обработчик для сохранения статьи
func saveArticleHandler(w http.ResponseWriter, r *http.Request) {
	// Сохраняем статью и отправляем успешный ответ
	fmt.Fprint(w, "Article saved successfully")
}
