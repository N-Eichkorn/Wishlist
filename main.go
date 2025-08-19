// +---------------------------------------------------+
// | Author: Niklas Eichkorn
// | Date: 07.08.25
// | Version: 1.0
// |---------------------------------------------------+
// | Notes: merke https://gist.github.com/jordansissel/1e08b1c65157bde0f30a87c4fb569237
// +---------------------------------------------------+

package main

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.
import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"

	"github.com/joho/godotenv"
)

func main() {

	argus := os.Args
	if len(argus) >= 2 {
		switch argus[1] {
		case "--init":
			init_programm(argus)
		case "--help":
			fmt.Print("Possible Options: \n " +
				"\t --help \t show this help site \n" +
				"\t --init <path> \t create database  at <path>. If path is empty the path is ./data.db \n")
		default:
		}
	}
	//Load ENV file-----------------------
	err := godotenv.Load("settings.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func init_programm(argus []string) {
	var database_location string
	//setup database ----------------------------------------------
	if len(argus) > 2 {
		database_location = argus[2]
	} else {
		database_location = "./data.db"
	}

	db, err := sql.Open("sqlite", database_location)
	if err != nil {
		fmt.Print(err)
	}
	sql_init, _ := os.ReadFile("sql_init.sql")
	_, err = db.Exec(string(sql_init))
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("\t + database initiated")
	db.Close()

	//setup settings.env ------------------------------------------
	settings_location := "settings.env"
	if _, err := os.Stat(settings_location); errors.Is(err, os.ErrNotExist) {
		os.Create(settings_location)
		fmt.Println("\t + " + settings_location + " created")
	} else {
		fmt.Println("\t + " + settings_location + " is existing")
	}

	//write settings ---------------------------------------------
	data := []byte("DATABASE=" + database_location)
	os.WriteFile(settings_location, data, 0644)
}
