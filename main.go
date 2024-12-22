package main

import (
	"fmt"
	"net/http"
	"regexp"
)

// Middleware для установки значения по умолчанию для имени
func SetDefaultName(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")

		// Если имя пустое, устанавливаем значение по умолчанию
		if name == "" {
			r.URL.Query().Set("name", "stranger")
		}

		// Обновляем параметры формы
		r.ParseForm()

		next(w, r)
	}
}

// Middleware для проверки корректности имени
func Sanitize(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")

		// Проверяем, содержит ли имя что-то кроме английских букв
		if isDirtyHacker(name) {
			fmt.Fprintf(w, "Hello, dirty hacker!")
			return
		}

		next(w, r)
	}
}

// Основной обработчик
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "stranger" {
		fmt.Fprintf(w, "Hello, stranger!")
	} else {
		fmt.Fprintf(w, "Hello, %s!", name)
	}
}

// Функция для проверки, содержит ли строка что-то кроме английских букв
func isDirtyHacker(name string) bool {
	// Регулярное выражение для проверки, состоит ли строка только из английских букв
	match, _ := regexp.MatchString("^[a-zA-Z]+$", name)
	return !match
}

func main() {
	// Устанавливаем обработчик для маршрута /hello с использованием middleware
	http.HandleFunc("/hello", SetDefaultName(Sanitize(HelloHandler)))

	// Запускаем сервер на порту 8080
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
