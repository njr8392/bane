package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/jackc/pgx"
	"log"
	"os"
	"strconv"
)

var Dbase *pgx.Conn

//will take a create table statment,column names/types and make the table
//no hyphens in database column names.... need escape phrase??
//postgres automatically converts all column names to lowercase
// combined with use of 2 Frame methods it will return
//CREATE TABLE person (Username varchar(100),  Identifier int, Firstname varchar(100), Lastname varchar(100));
func CreateTable(name, cols string) error {
	t := fmt.Sprintf("CREATE TABLE %s (%s);", name, cols)
	_, err := Dbase.Exec(context.Background(), t)
	if err != nil {
		return err
	}
	return nil
}

// pass the frame and iterate over every row adding to to a new array of {}interface and execute the same for each row
// wil create an accepatble insert statement such as
// INSERT INTO person (username,  identifier, firstname, lastname) VALUES($1, $2, $3, $4)
func InsertData(frame Frame) error {
	colnames := frame.DbColNames()
	vals := frame.SqlVals()
	insert := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", *table, colnames, vals) //table is provided via a flag
	for _, data := range frame.Data {
		var row []interface{}
		for _, val := range data {

			num, err := strconv.Atoi(val)
			if err != nil {
				row = append(row, val)
				continue
			}
			row = append(row, num)
		}

		//remeber to use []T... when providing an argument to variadic function!
		//you will get now build error and pgx will not throw an error
		_, err := Dbase.Exec(context.Background(), insert, row...)
		if err != nil {
			fmt.Errorf("error inserting %s into the table: %s", row, err)
		}
	}
	return nil
}

//write function to query the table command needs to read from stdin and output new data until the user is satisfied
//command for to write that output to a file then close????
//yes, key press for the program to exit and delete the table
func Select() {
	//need to abide by query formula. query will have place holders and values will be enter seperatly
	var rows pgx.Rows
	var er error
	for {
		var vals []interface{}
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter a query or press q to quit: ")

		//query needs to be entered in the same format as you would writing code
		//select * from tmp where name = $1, pam
		//need to use a comma to seperate query and arguments
		input, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
		}
		if input[0] == 'q' && len(input) == 2 { //byte len is two because it contains the newline character
			return
		}
		query := bytes.Split(input, []byte(","))
		//if more than 2 values add to args

		if len(query) > 1 {
			for i := 1; i < len(query); i++ {
				vals = append(vals, string(query[i])) // ints?
			}
			rows, er = Dbase.Query(context.Background(), string(query[0]), vals...) // why isn't this working? a basic select statement works.

		} else {
			rows, er = Dbase.Query(context.Background(), string(query[0]))

		}
		if er != nil {
			log.Fatal(er)
		}

		for rows.Next() {
			row, err := rows.Values()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(row)

		}
		rows.Close()
		//should be able to convert the results of the query back to a csv:
	}
}

//make it so this will get call when program exits
//only want the table to tempory
func DeleteTable() error {
	del := fmt.Sprintf("DROP TABLE %s;", *table)
	_, err := Dbase.Exec(context.Background(), del)
	return err
}
