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
const DEFAULT_OUTPUT_FILENAME = "./output.csv"
const ARCHIVE_PATH = "./archive"

// Golang Format layout
const DATEFORMAT = "01/02/2006"

type categoryMap struct {
	Categories map[string]string `json:"categories"`
	Names      map[string]string `json:"names"`
}

type rowConfig struct {
	Date     int `json:"date"`
	Category int `json:"category"`
	Note     int `json:"note"`
	Amount   int `json:"amount"`
}

func main() {
	input, bankName, outputName := handleArgs()
	configDecoder := loadDecoder()
	rowConfig := loadBankRowConfig(bankName)
	records := append(make([][]string, 0), []string{"Date", "Category", "Note", "Amount"})

	fileInfo, err := ioutil.ReadDir(input)
	if err != nil {
		log.Fatalln("error opening file:", err)
	}

	for _, file := range fileInfo {
		filePath := fmt.Sprintf("%s/%s", input, file.Name())
		records = parseFile(filePath, records, configDecoder, rowConfig)

		archiveFile(filePath, file.Name())
	}

	fmt.Printf("Found %d operations, writing to %s \n", len(records), outputName)
	writeCSV(records, outputName)
}

func loadJson(path string) *os.File {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}

	return jsonFile
}

func loadDecoder() categoryMap {
	jsonFile := loadJson("./config/categories.json")
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}
	var configDecoder categoryMap
	json.Unmarshal([]byte(byteValue), &configDecoder)

	return configDecoder
}

func loadBankRowConfig(bankName string) rowConfig {
	jsonFile := loadJson(fmt.Sprintf("./config/%sRowConfig.json", bankName))
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}
	var bankRowConfig rowConfig
	json.Unmarshal([]byte(byteValue), &bankRowConfig)

	return bankRowConfig
}

func parseRow(row []string, configDecoder categoryMap, rowConfig rowConfig) []string {
	return []string{
		parseDate(row[rowConfig.Date]),
		parseCategory(row[rowConfig.Category], row[rowConfig.Note], configDecoder),
		parseNote(row[rowConfig.Note]),
		row[rowConfig.Amount],
	}
}

func parseFile(inputName string, records [][]string, configDecoder categoryMap, rowConfig rowConfig) [][]string {
	rows := readInputFile(inputName)

	for i := SKIP_LINES; i < len(rows); i++ {
		row := rows[i]
		records = append(records, parseRow(row, configDecoder, rowConfig))
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

func handleArgs() (string, string, string) {
	args := os.Args
	filename := DEFAULT_OUTPUT_FILENAME
	bankName := ""

	if len(args) < 3 {
		log.Fatalln("Usage: csvparser path/to/input bankName [path/to/output]")
		return "", "", ""
	}
	if len(args) == 3 {
		bankName = args[2]
	} else {
		bankName = args[2]

		filename = args[3]
	}

	return args[1], bankName, filename
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
