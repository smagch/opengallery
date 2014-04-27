package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
)

func main() {
	httpAddr := flag.String("http", ":8080", "http address to listen")
	postgresUrl := flag.String("postgres-url", "", "postgres url to listen")
	flag.Parse()

	if *postgresUrl == "" {
		log.Fatal("an empty postgres-url given")
	}

	var err error
	db, err = sql.Open("postgres", *postgresUrl)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(4)

	mux := App()
	http.Handle("/", mux)
	err = http.ListenAndServe(*httpAddr, nil)
	log.Fatal(err)
}
