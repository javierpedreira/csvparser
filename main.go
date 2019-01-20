package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

const SKIP_LINES = 6
const SHEET_NAME = "Movimientos"
const DEFAULT_OUTPUT_FILENAME = "output.csv"

type categoryMap struct {
	Categories map[string]string `json:"categories"`
	Names      map[string]string `json:"names"`
}

func main() {
	inputName, outputName := handleArgs()

	rows := readInputFile(inputName)
	records := make([][]string, 0)
	records = append(records, []string{"Date", "Category", "Note", "Amount"})
	configDecoder := loadDecoder()

	for i := SKIP_LINES; i < len(rows); i++ {
		row := rows[i]

		records = append(records, []string{
			parseDate(row[0]),
			parseCategory(row[2], row[3], configDecoder),
			parseNote(row[3]),
			row[6],
		})
	}

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

func parseNote(note string) string {
	return strings.TrimSpace(strings.Split(note, "(")[0])
}

func parseDate(date string) string {
	dateSplit := strings.Split(date, "/")

	return fmt.Sprintf("%s/%s/%s", dateSplit[1], dateSplit[0], dateSplit[2])
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
