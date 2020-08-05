package logger

import (
	"fmt"
	"log"
	"os"
)

// InfoLogger exported
var InfoLogger *log.Logger

// ErrorLogger exported
var ErrorLogger *log.Logger

func init() {
	//relPath, err := filepath.Rel(".")
	//if err != nil {
	//	fmt.Println("Error reading given path:", err)
	//}

	generalLog, err := os.OpenFile("rspw.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	InfoLogger = log.New(generalLog, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(generalLog, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
}
