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

	"github.com/rivo/tview"

	_ "modernc.org/sqlite"

	"github.com/joho/godotenv"
)

const (
	env_database              = "DATABASE"
	env_wishlist_user         = "WISHLIST_USER"
	settings_location         = "settings.env"
	default_database_location = "./data.db"
)

func main() {

	argus := os.Args
	if len(argus) >= 2 {
		switch argus[1] {
		case "--init":
			init_programm(argus)
		case "--help":
			print_help()
		default:
		}
	}
	//Load ENV file----------------------------------------------
	if godotenv.Load(settings_location) != nil {
		log.Fatal("Error loading .env file")
	}

	//Check if user is registerd ----------------------------------------------
	if os.Getenv(env_wishlist_user) == "null" {
		if !register_user() {
			os.Exit(0)
		}
	}

	//Start main window ----------------------------------------------
	print_main_window()
}

func init_programm(argus []string) {
	var database_location string
	//setup database ----------------------------------------------
	if len(argus) > 2 {
		database_location = argus[2]
	} else {
		database_location = default_database_location
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

	//setup settings.env ----------------------------------------------

	if _, err := os.Stat(settings_location); errors.Is(err, os.ErrNotExist) {
		os.Create(settings_location)
		fmt.Println("\t + " + settings_location + " created")
	} else {
		fmt.Println("\t + " + settings_location + " is existing")
	}

	//write settings ----------------------------------------------
	data := []byte(env_database + "=" + database_location + "\n" +
		env_wishlist_user + "=" + "null\n")
	os.WriteFile(settings_location, data, 0644)
}

func print_help() {
	fmt.Print("Possible Options: \n " +
		"\t --help \t show this help site \n" +
		"\t --init <path> \t create database  at <path>. If path is empty the path is ./data.db \n")
}

func register_user() bool {
	return_value := true
	app := tview.NewApplication()
	form := tview.NewForm().
		AddInputField("Username", "", 20, nil, nil).
		AddButton("Save", func() {
			app.Stop()
		}).
		AddButton("Cancel", func() {
			return_value = false
		})
	form.SetBorder(false).SetTitle("Register your User").SetTitleAlign(tview.AlignLeft)
	if err := app.SetRoot(form, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}
	os.Setenv(env_wishlist_user, form.GetFormItem(0).(*tview.InputField).GetText())

	db, err := sql.Open("sqlite", os.Getenv(env_database))
	if err != nil {
		fmt.Print(err)
	}
	_, err = db.Exec("INSERT INTO Users VALUES ('" + os.Getenv(env_wishlist_user) + "');")
	if err != nil {
		fmt.Print(err)
	}
	write_settings()
	db.Close()
	return return_value
}

func write_settings() {
	data := []byte(env_database + "=" + os.Getenv(env_database) + "\n" +
		env_wishlist_user + "=" + os.Getenv(env_wishlist_user) + "\n")
	os.WriteFile(settings_location, data, 0644)
}

func print_main_window() {
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	menu := newPrimitive("Menu")
	main := newPrimitive("Main content")

	button := tview.NewButton("Hit Enter to close")
	button.SetBorder(true).SetRect(0, 0, 22, 3)

	grid := tview.NewGrid().
		SetRows(3, 0).
		SetColumns(30, 0).
		SetBorders(true).
		AddItem(tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText("Start Programm"), 0, 0, 1, 2, 0, 0, false).
		AddItem(menu, 1, 0, 1, 1, 0, 100, false).
		AddItem(main, 1, 1, 1, 1, 0, 100, false)

	if err := tview.NewApplication().SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
