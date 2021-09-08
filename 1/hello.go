package main

import (
	"bufio"
	"encoding/csv"
	"strconv"
	"strings"

	//"encoding/csv"
	//"fmt"
	"fmt"
	"os"

	//"os"

	"github.com/olekukonko/tablewriter"
)

type step int

const (
	stepType step = iota
	createTable
	insertTable
	selectTable
	deleteTable
	//updateTable
)

type parser struct {
	i               int
	sql             string
	stepType		step
}

//func addColumns() (err error) {
//	var tableName string = "man"
//	var columns = []string{"id", "name", "last_name", "age"}
//
//	var resultTableName = fmt.Sprintf("%s.csv", tableName)
//
//	if err := addColumns(columns, resultTableName); err != nil {
//		panic(err)
//	} else {
//		fmt.Println("Создана таблица под названием", resultTableName)
//	}
//	csvOut, err := os.Create(resultTableName)
//
//	if err != nil {
//		return err
//	}
//
//
//	w := csv.NewWriter(csvOut)
//
//	var record []string
//	for _, value := range columns {
//		fmt.Printf(value, record)
//		//record = append(record[:0], value) // copy data set 1 to beginning of output record
//		record = append(record, value) // append data set 2 to output record
//
//	}
//	w.Write(record)
//	w.Flush()
//
//	return
//
//}


func insertData(args []string, input string) {
	if len(args) > 2 && args[0] == "insert" && args[1] == "into" {

		fmt.Println("INSERT DATA")
		var firstPart = input[:strings.Index(input, "values")]
		var table = strings.TrimPrefix(firstPart, "insert into")
		var values = strings.TrimPrefix(input, "insert into" + table + "values")

		table = strings.TrimSpace(table)

		values = strings.TrimSpace(values)
		values = strings.Trim(values,"()")
		var valuesArr = strings.Split(values, ", ")

		file, err := os.OpenFile(table + ".csv", os.O_RDWR, 0755)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		reader.Comma = ','
		writer := csv.NewWriter(file)

		var tableRows [][]string

		for {
			record, e := reader.Read()
			if e != nil {
				fmt.Println(e)
				break
			}

			tableRows = append(tableRows, record)
		}

		fmt.Println(tableRows[0], " ", len(tableRows[0]), " ", len(valuesArr))

		if len(tableRows[0]) - 1 == len(valuesArr)  {
			var len = len(tableRows)

			var record []string
			record = append(record, strconv.Itoa(len))
			for _, value := range valuesArr {
				record = append(record, value)
			}

			err := writer.Write(record)

			if err != nil {
				fmt.Println("Ошибка вставки ", err)
			}

			writer.Flush()

			for {
				record, e := reader.Read()
				if e != nil {
					fmt.Println(e)
					break
				}

				tableRows = append(tableRows, record)
			}

			table := tablewriter.NewWriter(os.Stdout)

			table.SetHeader(tableRows[0])

			for index, row := range tableRows {
				if index != 0 {
					table.Append(row)
				}
			}

			table.Render()

		} else {
			fmt.Println("∆ Ошибка! Ошибка вставки ∆")
		}

		fmt.Println(values)
	} else {
		fmt.Println("∆ Неправильно задана команда вставки строки в таблиц! ∆")
	}
}


func getDataForSelect(tableName string, fields []string) {
	file, err := os.Open(tableName + ".csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	var tableRows [][]string

	for {
		record, e := reader.Read()
		if e != nil {
			fmt.Println(e)
			break
		}

		tableRows = append(tableRows, record)
	}


	var rowsIndex []int

	//for index, item := range tableRows {
	//	if len(fields) > 1 && fields[0] == "*" {
	//		fmt.Println(item)
	//	} else if len(fields) > 1 && fields[0] != "*" {
	//		if index == 0 {
	//			for it, i := range item {
	//
	//			}
	//		}
	//	}
	//}

	table := tablewriter.NewWriter(os.Stdout)

	//table.SetHeader(fields)

	if len(fields) == 1 && fields[0] == "*" {
		for index, item := range tableRows {
			if index == 0 {
				table.SetHeader(item)
			} else {
				table.Append(item)
			}
		}
	} else if len(tableRows) > 1 {
		var columns = tableRows[0]
		for _, item := range fields {
			for index, row := range columns {
				if row == item {
					rowsIndex = append(rowsIndex, index)
				}
			}
		}

		table.SetHeader(fields)

		fmt.Println(fields)

		for index, row := range tableRows {

			if index != 0 {
				var result []string
				for _, i := range rowsIndex {
					result = append(result, row[i])
				}
				table.Append(result)
			}
		}
	}


	table.Render()

	//fmt.Println(record)
	//if len(fields) > 1 && fields[0] == "*" {
	//
	//}
}


func selectData(args []string, input string)  {
	//var dataArr = append(args[:0], args[1:]...)
	var parameters = input[:strings.Index(input, "from")]
	var fields = strings.TrimPrefix(parameters, "select")
	var table = strings.TrimPrefix(input, "select" + fields + "from")

	fields = strings.TrimSpace(fields)

	table = strings.TrimSpace(table)

	var fieldsArr = strings.Split(fields, ", ")

	getDataForSelect(table, fieldsArr)

	//for _, item := range dataArr {
	//
	//}

}

func execInput(input string) {
	// Remove the newline character.
	input = strings.TrimSuffix(input, "\n")

	// Split the input separate the command and the arguments.
	args := strings.Split(input, " ")

	for _, value := range args {
		switch value {
		case "select":
			selectData(args, input)
		case "insert":
			insertData(args, input)
		}
		
	}

}
 
func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		// Read the keyboad input.
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		execInput(input)
		//if err = execInput(input); err != nil {
		//	fmt.Fprintln(os.Stderr, err)
		//}
	}

}