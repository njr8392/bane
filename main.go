package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/jackc/pgx"
	"log"
	"os"
)

func main() {
	flag.Parse()
	csvfile := os.Args[len(os.Args)-1]
	file, err := os.Open(string(csvfile))
	if err != nil {
		fmt.Println(err)
	}
	r := csv.NewReader(file)
	records, err := r.ReadAll()
	f := NewFrame(records)
	if err != nil {
		fmt.Println(err)
	}
	// handles headers and column flags
	Run(f)

	//./bane -db=true -t=person *csv.  will make a tempory sql table that you can query
	if *database {
		var host string = "localhost"
		var port int = 5432
		var user string = "user"
		var password string = "password"
		var dbname string = "db"

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s", host, port, user, password, dbname)
		Dbase, err = pgx.Connect(context.Background(), psqlInfo)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("db connected")
		Database(f, Dbase)
		Dbase.Close(context.Background())
	}

}
