package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/tealeg/xlsx"
)

func UploadData(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		fmt.Println("GET")
		t, _ := template.ParseFiles("./templates/index.html")
		t.Execute(w, nil)

	} else if req.Method == "POST" {
		fmt.Println("POST")
		file, handler, err := req.FormFile("uploadfile")
		defer file.Close()
		if err != nil {
			log.Printf("Error while Posting data")
			t, _ := template.ParseFiles("./templates/index.html")
			t.Execute(w, nil)
		} else {
			fmt.Println("error throws in else statement")
			fmt.Println("handler.Filename", handler.Filename)
			fmt.Printf("Type of handler.Filename:%T\n", handler.Filename)
			fmt.Println("Length:", len(handler.Filename))
			f, err := os.OpenFile("./data/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println("Error:", err)
				t, _ := template.ParseFiles("./templates/index.html")
				t.Execute(w, nil)
			}
			defer f.Close()
			io.Copy(f, file)
			filePath := "./data/" + handler.Filename
			var extension = filepath.Ext(filePath)
			parsedData := ExcelCsvParser(filePath, extension)
			fmt.Println(parsedData) // this is dummy to shut up the compiler
		}
	} else {
		log.Printf("Error while Posting data")
		t, _ := template.ParseFiles("./templates/index.html")
		t.Execute(w, nil)

	}
}

func ExcelCsvParser(filePath string, pathExtension string) (parsedData []map[string]interface{}) {

	if pathExtension == ".xlsx" {
		fmt.Println("----------------We are parsing an xlsx file.---------------")
		parsedData := ReadXlsxFile(filePath)
		return parsedData
	}
	return parsedData
}

func main() {
	fmt.Println("Web server running at localhost:8000")
	router := mux.NewRouter()
	router.HandleFunc("/", UploadData)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./templates/")))
	log.Fatal(http.ListenAndServe(":8000", router))
}

func ReadXlsxFile(filePath string) []map[string]interface{} {
	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		fmt.Println("Error reading the file")
	}

	parsedData := make([]map[string]interface{}, 0, 0)

	//sheet
	for _, sheet := range xlFile.Sheets {

		// rows
		for _, row := range sheet.Rows {

			cells := make([]string, 0)

			for _, cell := range row.Cells {
				cellText := cell.String()
				cellText = strings.Replace(cellText, ",", ".", -1)

				cells = append(cells, cellText)
			}
			rowText := strings.Join(cells, ",")
			fmt.Println(rowText) // TODO :(

		}
	}
	return parsedData
}
