package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

const SKIP_LINES = 6
const SHEET_NAME = "Movimientos"
const DEFAULT_OUTPUT_FILENAME = "output.csv"

func main() {
	inputName, outputName := handleArgs()

	rows := readInputFile(inputName)
	records := make([][]string, 0)
	records = append(records, []string{"Date", "Category", "Note", "Amount"})

	for i := SKIP_LINES; i < len(rows); i++ {
		row := rows[i]

		records = append(records, []string{
			parseDate(row[0]),
			row[2],
			parseNote(row[3]),
			row[6],
		})
	}

	writeCSV(records, outputName)
}

func parseNote(note string) string {
	return strings.TrimSpace(strings.Split(note, "(")[0])
}

func parseDate(date string) string {
	dateSplit := strings.Split(date, "/")

	return fmt.Sprintf("%s/%s/%s", dateSplit[1], dateSplit[0], dateSplit[2])
}

func writeCSV(records [][]string, filename string) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Fatalln("error opening file:", err)
		return
	}

	w := csv.NewWriter(file)

	err = w.WriteAll(records)
	if err != nil {
		log.Fatalln("error writing csv:", err)
	}
}

func handleArgs() (string, string) {
	args := os.Args
	if len(args) < 2 {
		log.Fatalln("You must specify the name of the file")
		return "", ""
	}

	filename := DEFAULT_OUTPUT_FILENAME
	if len(args) == 3 {
		filename = args[2]
	}

	return args[1], filename
}

func readInputFile(filename string) [][]string {
	xlsx, err := excelize.OpenFile(filename)
	if err != nil {
		log.Fatalln("error opening xlsx file", err)
		return nil
	}

	return xlsx.GetRows(SHEET_NAME)
}
