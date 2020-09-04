package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
)

const version = "1.1"

const sequenceResetTemplate = `
	begin
		perform (select setval('%s.%s'::regclass, coalesce((select max(%s) from %s.%s), 1)));
	exception
		when others then
			null;
	end;
`

const transactionTemplate = `
do
$$
begin
%s
end;
$$
`

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

func inputFilePaths() (tables string) {
	if len(os.Args) >= 2 {
		tables = os.Args[1]
	} else {
		tables = "tables.csv"
	}
	return
}

func main() {
	tablesFilepath := inputFilePaths()
	if tablesFilepath == "version" {
		fmt.Printf("rsq version: %s\n", version)
		return
	} else if tablesFilepath == "help" {
		fmt.Println("usage:")
		fmt.Println(">>> rsq [input filepath]")
		fmt.Println("default input filepath = tables.csv")
		fmt.Println("input file must be in csv format")
		fmt.Println("output filename = reset_sequences.sql")
		return
	}
	records, err := csvReadAll(tablesFilepath)
	if err != nil {
		fmt.Println(err)
		return
	}
	output := ""
	for _, record := range records[1:] {
		if len(record) >= 4 && record[2] != "" && record[3] != "" {
			output += fmt.Sprintf(sequenceResetTemplate,
				record[0],
				record[2],
				record[3],
				record[0],
				record[1],
			)
		}
	}
	output = fmt.Sprintf(transactionTemplate, output)
	err = ioutil.WriteFile("reset_sequences.sql", []byte(output), 0662)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
