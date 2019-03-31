package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

const SKIP_LINES = 6
const SHEET_NAME = "Movimientos"
const DEFAULT_OUTPUT_FILENAME = "output.csv"
const ARCHIVE_PATH = "./archive"

// Golang Format layout
const DATEFORMAT = "01/02/2006"

type categoryMap struct {
	Categories map[string]string `json:"categories"`
	Names      map[string]string `json:"names"`
}

func main() {
	input, outputName := handleArgs()

	configDecoder := loadDecoder()
	records := append(make([][]string, 0), []string{"Date", "Category", "Note", "Amount"})

	fileInfo, err := ioutil.ReadDir(input)
	if err != nil {
		log.Fatalln("error opening file:", err)
	}

	for _, file := range fileInfo {
		filePath := fmt.Sprintf("%s/%s", input, file.Name())
		records = parseFile(filePath, records, configDecoder)

		archiveFile(filePath, file.Name())
	}

	fmt.Printf("Found %d operations, writing to %s \n", len(records), outputName)
	writeCSV(records, outputName)
}

func loadDecoder() categoryMap {
	jsonFile, err := os.Open("categories.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var configDecoder categoryMap
	json.Unmarshal([]byte(byteValue), &configDecoder)

	return configDecoder
}

func parseFile(inputName string, records [][]string, configDecoder categoryMap) [][]string {
	rows := readInputFile(inputName)

	for i := SKIP_LINES; i < len(rows); i++ {
		row := rows[i]
		records = append(records, []string{
			parseDate(row[0]),
			parseCategory(row[2], row[3], configDecoder),
			parseNote(row[3]),
			row[8],
		})
	}
	return records
}

func parseNote(note string) string {
	if strings.Contains(note, "Traspaso") || strings.Contains(note, "Transferencia") {
		return note
	}
	return strings.TrimSpace(strings.Split(note, "(")[0])
}

func parseDate(date string) string {
	// Excel Date starts in 1/1/1900 + 2 days
	start := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	days, _ := strconv.ParseInt(date, 0, 64)
	d := start.AddDate(0, 0, int(days)-2)

	return d.Format(DATEFORMAT)
}

func parseCategory(category, name string, configDecoder categoryMap) string {
	for k, v := range configDecoder.Names {
		if strings.Contains(name, k) {
			return v
		}
	}

	newCat, ok := configDecoder.Categories[category]
	if ok {
		return newCat
	}

	return category
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
		log.Fatalln("Usage: csvparser path/to/input [path/to/output]")
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

func archiveFile(filePath, filename string) {
	archiveFile := fmt.Sprintf("%s/%s", ARCHIVE_PATH, filename)
	err := os.Rename(filePath, archiveFile)
	if err != nil {
		fmt.Println(err)
	}
}
