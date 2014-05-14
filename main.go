package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	httpAddr := flag.String("http", ":8080", "http address to listen")
	postgresUrl := flag.String("postgres-url", "", "postgres url to listen")
	useImport := flag.Bool("import", false, "use data import instead of server")
	flag.Parse()

	if *postgresUrl == "" {
		log.Fatal(`"postgres-url" option is required`)
		os.Exit(1)
	}

	if _, err := url.Parse(*postgresUrl); err != nil {
		log.Fatal("Invalid postgres URL: ", err.Error())
		os.Exit(1)
	}

	var err error
	db, err = sql.Open("postgres", *postgresUrl)
	if err != nil {
		log.Fatal("Cannot open a connection with postgresql: ", err.Error())
		os.Exit(1)
	}

	// TODO max conn
	db.SetMaxOpenConns(4)

	if *useImport {
		for _, filepath := range flag.Args() {
			log.Printf("Importing %s\n", filepath)
			if err = ImportFixture(filepath); err != nil {
				log.Fatalf("Failed to import %s: %s", filepath, err.Error())
				os.Exit(1)
			} else {
				log.Println("Import succeed")
			}
		}
		os.Exit(0)
	}

	mux := App()
	http.Handle("/", mux)
	err = http.ListenAndServe(*httpAddr, nil)
	log.Fatal(err)
}
