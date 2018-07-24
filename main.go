package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	context "mod/github.com/gorilla/context@v1.1.1"
	mux "mod/github.com/gorilla/mux@v1.6.2"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/products/{locale}", ProductHandler)

	r.Use(Middleware)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("middleware", r.URL)
		context.Set(r, "locale", mux.Vars(r)["locale"])

		h.ServeHTTP(w, r)
	})
}

func ProductHandler(w http.ResponseWriter, r *http.Request) {
	JSON(w, Global.ErrOne[Locale(context.Get(r, "locale"))])

	return
}

func JSON(w http.ResponseWriter, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Fprint(w, err)
	}
	fmt.Fprint(w, string(b))
}

type Language string

func Locale(v interface{}) Language {
	if v == nil {
		return EN
	}
	return Language(v.(string))
}

const (
	TH Language = "th"
	EN Language = "en"
)

var Global = Multiple{
	ErrOne: map[Language]ErrorType{
		TH: ErrorType{Code: "01", Desc: "ไทย"},
		EN: ErrorType{Code: "01", Desc: "Eng"},
	},
}

type Multiple struct {
	ErrOne map[Language]ErrorType
}

type ErrorType struct {
	Code string `json:"code"`
	Desc string `json:"description"`
}
