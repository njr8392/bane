package main

import (
	"flag"
	"fmt"
	"github.com/jackc/pgx"
	"log"
	"strconv"
	"strings"
)

var (
	titles   = flag.Bool("h", false, "List column headers")
	columns  = flag.String("c", "", "List what columns you want from the sheet")
	database = flag.Bool("db", false, "active the database feature")
	table    = flag.String("t", "", "name of the table you are creating")
)

func Run(f Frame) {

	if *titles {
		f.PrintHeaders()
	}
	//./bane -c 1,4 *.csv will write new csv of those columns to stdout. Can pipe into new file
	//should add functionality to split by column name as well
	if *columns != "" {
		holder := []int{}
		for _, arg := range strings.Split(*columns, ",") {
			num, err := strconv.Atoi(string(arg))
			if err != nil {
				log.Fatalf("Columns need to be defined as an index: %s", err)
			}
			holder = append(holder, num)
		}
		sheet := f.GetColumns(holder...)
		sheet.Write()
	}
}

//./bane -db=true -t=person *csv.  will make a tempory sql table that you can query
func Database(f Frame, db *pgx.Conn) {

	if *database {
		colnames := f.SqlTable()
		fmt.Printf("Columns create: %s\n\n", colnames)

		err := CreateTable(*table, colnames)
		fmt.Printf("Table %s created\n", *table)

		if err != nil {
			fmt.Println(err)
		}

		err = InsertData(f)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("All data has been correctly inserted")
		}
		//fixed so users can enter custom queries. Enter infinite loop and scan stdin?
		Select()
		err = DeleteTable()

		if err != nil {
			fmt.Println(err)
		}

	}
}
