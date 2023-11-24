package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Привет, мир!")
}

func main() {
	http.HandleFunc("/", handler)

	log.Println("Запуск сервера на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
