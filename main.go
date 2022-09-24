package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
)

var hellos []string = []string{"Welcome!", "hallo", "Përshëndetje", "ሰላም", "مرحبا", "Բարեւ", "Salam", "Kaixo", "добры, дзень", "হ্যালো", "zdravo", "Здравейте", "Hola", "Hello", "Moni", "您好", "您好", "Bonghjornu", "zdravo", "Ahoj", "Hej", "Hallo", "Hello", "Saluton", "Tere", "Kumusta", "Hei", "Bonjour", "Hello", "Ola", "გამარჯობა", "Hallo", "Γεια, σας", "હેલો", "Bonjou", "Sannu", "Alohaʻoe", "שלום", "नमस्ते", "Nyob, zoo", "Helló", "Halló", "Ndewo", "Halo", "Dia, duit", "Ciao", "こんにちは", "Hello", "ಹಲೋ", "Сәлем", "ជំរាបសួរ", "안녕하세요.", "Hello", "салам", "ສະບາຍດີ", "salve", "Labdien!", "Sveiki", "Moien", "Здраво", "Hello", "Hello", "ഹലോ", "Hello", "Hiha", "हॅलो", "Сайн, байна, уу", "မင်္ဂလာပါ", "नमस्ते", "Hallo", "سلام", "سلام", "Cześć", "Olá", "ਹੈਲੋ", "Alo", "привет", "Talofa", "Hello", "Здраво", "Hello", "Hello", "هيلو", "හෙලෝ", "ahoj", "Pozdravljeni", "Hello", "Hola", "halo", "Sawa", "Hallå", "Салом", "ஹலோ", "హలో", "สวัสดี", "Merhaba", "Здрастуйте", "ہیلو", "Salom", "Xin, chào", "Helo", "Sawubona", "העלא", "Kaabo", "Sawubona"}

func main() {
	port := os.Getenv("PORT")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"message\": \"" + hellos[rand.Intn(len(hellos))] + "\"}"))
	})

	fmt.Println("listening on", port)
	http.ListenAndServe(":"+port, nil)
}
