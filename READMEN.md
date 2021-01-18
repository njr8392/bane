Bane
The Reckoning of CSV Files

I have 3 goals for this project
	1.  Create an interface that makes it easy for users to manipulate data expressed in a 2D matrix (csv and excel files)
		written in Go, giving them to tools to automate IO when it comes to csv files
	
	2.  A lot csv files that I have worked with have a lot of garbage that you can throw away. Sure, you can just open your
		spreadsheet editor of choice but why not do this from the command line?  Simply specify the columns you want and it 
		will out a new csv file of those contents.  I would like to add more functionality that could be chained together
		with pipes.
		
	3.	Sometimes it is not a matter of keeping/removing columns.  You need more power. You need to query these csv files, I
		chose SQL since it is what I am familar with,  Currently, it would be more be more robust to use a CLI provided by your
		database of choice, but writing create table statements for a large file is too much for me! You could use this to 
		generate a Create table statement automatically but it currently only supports strings and integers, however a lot of
		the time that is all you need.  As a work around you currently have to write queries to fit the SQL ie seperating
		the queries and values.  I would like to eventually move away from this as I don't care about risk of database 
		injections because this is intended to run locally.  I may have to write my own driver for this but I will do more 
		research.

Cut with ./bane -c 2,3,4 (columns of choice) *.csv

Query with ./bane -db -t (table name) *.csv.  NOTE you will have to change the database variables in main.go
Check it out, fork it, make changes!
