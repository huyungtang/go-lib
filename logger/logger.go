package logger

import (
	"encoding/json"
	"fmt"
	"log"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Fatal
// ****************************************************************************************************************************************
func Fatal(err error) {
	log.Fatal(err)
}

// Log
// ****************************************************************************************************************************************
func Log(args ...any) {
	log.Print(args...)
}

// Logf
// ****************************************************************************************************************************************
func Logf(msg string, args ...any) {
	log.Printf(msg, args...)
}

// Json
// ****************************************************************************************************************************************
func Json(dto any) {
	bs, err := json.Marshal(dto)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(bs))
}

// Print
// ****************************************************************************************************************************************
func Print(args ...any) {
	fmt.Print(args...)
}

// Printf
// ****************************************************************************************************************************************
func Printf(msg string, args ...any) {
	fmt.Printf(msg, args...)
}

// Println
// ****************************************************************************************************************************************
func Println(args ...any) {
	fmt.Println(args...)
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
