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
	JSON(w, Speak(context.Get(r, "locale")).ErrorTwo)

	return
}

func JSON(w http.ResponseWriter, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Fprint(w, err)
	}
	fmt.Fprint(w, string(b))
}

var (
	EN = Global{
		ErrorOne:   ErrorType{Code: "01", Desc: "Eng"},
		ErrorTwo:   ErrorType{Code: "02", Desc: "Eng"},
		ErrorThree: ErrorType{Code: "03", Desc: "Eng"},
	}
	TH = Global{
		ErrorOne:   ErrorType{Code: "01", Desc: "ไทย"},
		ErrorTwo:   ErrorType{Code: "02", Desc: "ไทย"},
		ErrorThree: ErrorType{Code: "03", Desc: "ไทย"},
	}

	Language = map[interface{}]Global{
		"en": EN,
		"th": TH,
	}
)

func Speak(locale interface{}) Global {
	if locale == nil {
		return EN
	}

	if v, ok := Language[locale]; ok {
		return v
	}
	return EN
}

type Global struct {
	ErrorOne   ErrorType
	ErrorTwo   ErrorType
	ErrorThree ErrorType
}

type ErrorType struct {
	Code string `json:"code"`
	Desc string `json:"description"`
}
