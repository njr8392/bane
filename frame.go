package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// a Frame will create an easy way to manipulate data read from a csv
type Frame struct {
	clength int // column length of the frame
	rlength int // row length of the frame
	Data    [][]string
}

//add a row of digits across the top as well?
func (f Frame) String() string {
	var s string
	for i, row := range f.Data {
		s += fmt.Sprintf("%d: %s\n", i, row)
	}
	return s
}

//gets desired columns from the file
//problems .... go csv.Reader igonres whitespace. keep? would mean no blank data
//need to add column length to frame struct as well (tmp hard coded in main function)
func (f Frame) GetColumns(cols ...int) Frame {
	newSet := f
	//Caller can pull columns that are not adjacent
	newSet.Data = make([][]string, f.clength)
	for t := range newSet.Data {
		newSet.Data[t] = make([]string, len(cols))
	}
	for i, nums := range cols {
		for j := 0; j < f.clength; j++ {
			newSet.Data[j][i] = f.Data[j][nums]
		}
	}
	return newSet
}

//returns rows of a certain various indexs
func (f Frame) GetRows(rows ...int) Frame {
	newSet := f
	newSet.Data = make([][]string, len(rows))
	for t := range newSet.Data {
		newSet.Data[t] = make([]string, f.rlength)
	}
	for i, nums := range rows {
		for j := 0; j < f.rlength; j++ {
			newSet.Data[i][j] = f.Data[nums][j]
		}
	}
	return newSet
}

//Print headers of the frame
func (f Frame) PrintHeaders() {
	var heads string
	for i, head := range f.Data[0] {
		heads += fmt.Sprintf("%d: %s\n", i, head)
	}
	fmt.Println(heads)
}

//Writes to Stout for piping
func (f Frame) Write() {
	w := csv.NewWriter(os.Stdout)

	for _, row := range f.Data {
		if err := w.Write(row); err != nil {
			log.Fatal("Error writing to csv:", err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

func (f Frame) ConvertNums() Frame {
	newFrame := f
	for i, line := range f.Data {
		for j := 0; j < f.rlength; j++ {
			_, err := strconv.Atoi(line[j])
			if err != nil {
				continue
			}
			//need to fix! can't add an int to [][]string need to figure out which cols are ints
			//then need to make a new array of ints and add them to that but mix of ints and string?????
			//will probably have to chagne Data to [][]interface will need to add type checking
			newFrame.Data[i][j] = line[j]
		}
	}
	return newFrame
}

//makes the meat of the create table statment for a Postgres database
//can't have hyphens or spaces in columns names!!!!!--- need to fix
func (f Frame) SqlTable() string {
	var columns string

	for i, title := range f.Data[0] {
		columns += title + " "
		if _, err := strconv.Atoi(f.Data[1][i]); err != nil {
			columns += "varchar(100)"
		} else {
			columns += "int"
		}

		if i < f.rlength-1 {
			columns += ", "
		}
	}
	return columns
}

// postgres makes all column names lowercase, added it to hear to eliminate confusing to the user
//Have to adjust for inserting and querying
func (f Frame) DbColNames() string {
	var columns string

	for i, title := range f.Data[0] {
		t := strings.ToLower(title)
		columns += t
		//checking if can convert col name to int.. neeed to check on the first row of data
		if i < f.rlength-1 {
			columns += ", "
		}
	}
	return columns
}

//Returns a string of the for $1, $2, etc up to the length of the row
//necessary for the syntax of inserting values into the database
func (f Frame) SqlVals() string {
	var q string
	for i := 0; i < f.rlength; i++ {
		num := i + 1
		str := strconv.Itoa(num)
		q += "$" + str
		if i != f.rlength-1 {
			q += ", "
		}
	}
	return q
}

//constructor to take an multidemensional array and turn it into a Frame
func NewFrame(set [][]string) Frame {
	f := Frame{}
	f.rlength, f.clength, f.Data = len(set[0][:]), len(set), set
	return f
}
