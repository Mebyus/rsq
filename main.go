package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
)

const version = "1.1"

type tableMetadata struct {
	Schema                string
	Name                  string
	SequencePkAttribute   string
	ColumnPkAttribute     string
	ConstraintPkAttribute string
}

func csvReadAll(path string) (records [][]string, err error) {
	f, err := os.Open(path)
	if err != nil {
		err = fmt.Errorf("Открытие файла \"%s\" для чтения: %v", path, err)
		return
	}
	defer f.Close()
	r := csv.NewReader(f)
	records, err = r.ReadAll()
	if err != nil {
		err = fmt.Errorf("Чтение \"%s\" csv файла: %v", path, err)
		return
	}
	return
}

//чтение шаблона
func templateRead(path string) (text []byte, err error) {

	text, err = ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Чтение \"%s\" csv файла: %v", path, err)
		return
	}
	return
}

func main() {
	var templateFilepath, tablesFilepath string
	flag.StringVar(&tablesFilepath, "table", "", "path to csv table")
	flag.StringVar(&templateFilepath, "template", "", "path to template")
	v := flag.Bool("v", false, "display version")
	help := flag.Bool("help", false, "display help")
	flag.Parse()

	//если фраг указан в строке вызова
	if *v {
		fmt.Println(version)
		return
	}
	if *help {
		fmt.Println(version)
		fmt.Println("usage:")
		fmt.Println(">>> rsq [input filepath]")
		fmt.Println("default input filepath = tables.csv")
		fmt.Println("input file must be in csv format")
		fmt.Println("output filename = reset_sequences.sql")
		return
	}
	records, err := csvReadAll(tablesFilepath)
	if err != nil {
		log.Println(err)
		return
	}

	//заполнение структуры dbMetadata
	dbMetadata := make([]tableMetadata, 0)
	var tableMetadata tableMetadata
	for _, record := range records {
		if len(record) >= 4 && record[2] != "" && record[3] != "" {
			tableMetadata.Schema = record[0]
			tableMetadata.Name = record[1]
			tableMetadata.SequencePkAttribute = record[2]
			tableMetadata.ColumnPkAttribute = record[3]
			dbMetadata = append(dbMetadata, tableMetadata)
		}
	}

	//читаем шаблон из файла
	templateText, err := templateRead(templateFilepath)
	if err != nil {
		log.Println(err)
		return
	}

	template := template.New("")
	template.Parse(string(templateText))

	file, err := os.Create("./reset_sequences.sql")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	//заполнение шаблона
	err = template.Execute(file, dbMetadata[1:])

	return
}
