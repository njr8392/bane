package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

//add some flags in there but i need to parse them to ints... other option to pass col names?
//var colnum = flag.String("n", "", "Get column by names or numbers")

//abstract away the reader....
var (
	titles = flag.String("n", "", "List column headers")
	columns = flag.String("c", "", "List what columns you want from the sheet")
	//add remove flag as well
)
type Frame struct {
	clength int // column length of the frame
	rlength int // row length of the frame 
	Data [][]string
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
func(f Frame)GetColumns(cols ...int) Frame{
	newSet := f
	//have to change data to accomdate where slices won't work
	//Caller can pull columns that are not adjacent
	newSet.Data = make([][]string, f.clength)
	for t := range newSet.Data{
		newSet.Data[t] = make([]string, len(cols))	
	}
	for i, nums := range cols{
		for j :=0; j < f.clength; j++{
		newSet.Data[j][i] = f.Data[j][nums]
		}
	}
	return newSet
}

//returns rows of a certain various indexs
func (f Frame) GetRows(rows ...int)Frame{
	newSet := f
	newSet.Data = make([][]string, len(rows))
	for t := range newSet.Data{
		newSet.Data[t] = make([]string, f.rlength)	
	}
	for i, nums := range rows{
		for j :=0; j < f.rlength; j++{
		newSet.Data[i][j] = f.Data[nums][j]
		}
	}
	return newSet
}
func (f Frame) PrintHeaders(){
	var heads string
	for i, head := range f.Data[0]{
	heads += fmt.Sprintf("%d: %s\n", i, head)
	}
	fmt.Println(heads)
}

//Writes to Stderr for piping
func (f Frame) Write(){
	w := csv.NewWriter(os.Stdout)

	for _,row := range f.Data{
		if err := w.Write(row); err != nil {
			log.Fatal("Error writing to csv:", err)
		}
	}
		
	w.Flush()
	if err := w. Error(); err != nil{
		log.Fatal(err)
	}
}

//constructor to take an multidemensional array and turn it into a Frame
func NewFrame(set [][]string)Frame{
	f := Frame{}
	f.rlength, f.clength, f.Data = len(set[0][:]), len(set), set
	return f
}

func main() {
	flag.Parse()
	test := os.Args[len(os.Args)-1]
	file, err := os.Open(string(test))
	if err != nil {
		fmt.Println(err)
	}
	r := csv.NewReader(file)
	records, err := r.ReadAll()
	f := NewFrame(records)
	if err != nil {
		fmt.Println(err)
	}
	if *titles != "" {
		f.PrintHeaders()
	}
	if *columns != "" {
		inp := os.Args[1 : len(os.Args)-1]
		holder := []int{}
		for _, arg := range inp {
			num, _ := strconv.Atoi(arg)
			holder = append(holder, num)
		}
		sheet := f.GetRows(holder...)
		sheet.Write()
		fmt.Printf("%s\n", sheet)
	}
}
