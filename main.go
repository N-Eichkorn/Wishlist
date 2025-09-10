// +---------------------------------------------------+
// | Author: Niklas
// | Date: 01.09.2025
// | Version: 1.0
// |---------------------------------------------------+

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/rivo/tview"

	_ "modernc.org/sqlite"

	"github.com/joho/godotenv"
	"github.com/lxn/walk"
)

const (
	env_database              = "DATABASE"
	env_wishlist_user         = "WISHLIST_USER"
	env_wishlist_to           = "WISHLIST_TO"
	env_wishlist_wish         = "WISHLIST_WISH"
	env_broadcast             = "WISHLIST_BROADCAST"
	settings_location         = "settings.env"
	default_database_location = "./data.db"
	refresh_rate              = 30
)

type Wish struct {
	ID        int8
	FROM      string
	TO        string
	WISH      string
	TIMESTAMP string
	BROADCAST bool
}

func (w Wish) to_string() string {
	return string(w.FROM + " -> " + w.TO + "\n" + w.WISH)
}

var Wishes []Wish

var Users []string

func main() {

	argus := os.Args
	if len(argus) >= 2 {
		switch argus[1] {
		case "--init":
			init_programm(argus)
			os.Exit(0)
		case "--help":
			print_help()
			os.Exit(0)
		default:
		}
	}
	//Load ENV file----------------------------------------------
	load_env_file()

	//Check if user is registerd ----------------------------------------------
	if os.Getenv(env_wishlist_user) == "null" {
		if !register_user() {
			os.Exit(0)
		}
	}

	//Start Backgroundtask ----------------------------------------------
	go sceduler_get_wishes()

	//Start main window ----------------------------------------------
	print_main_window()
}

func load_env_file() {
	if godotenv.Load(settings_location) != nil {
		log.Fatal("Error loading .env file")
	}
}

func init_programm(argus []string) {
	var database_location string
	//setup database ----------------------------------------------
	if len(argus) > 2 {
		database_location = argus[2]
	} else {
		database_location = default_database_location
	}

	setup_database(database_location)

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

func setup_database(database_location string) {
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
}

func print_help() {
	fmt.Print("Possible Options: \n " +
		"\t --help \t show this help site \n" +
		"\t --init <path> \t create database  at <path>. Default is ./data.db \n")
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
	_, err = db.Exec("INSERT INTO Users VALUES (?);", os.Getenv(env_wishlist_user))
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
	get_wishes()
	app := tview.NewApplication()
	header := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("" +
		"__        ___     _     _ _     _    \n" +
		"\\ \\      / (_)___| |__ | (_)___| |_  \n" +
		" \\ \\ /\\ / /| / __| '_ \\| | / __| __| \n" +
		"  \\ V  V / | \\__ \\ | | | | \\__ \\ |_  \n" +
		"   \\_/\\_/  |_|___/_| |_|_|_|___/\\__| \n")

	alphabet := []int{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't'}
	w := tview.NewList()
	generate_wishlist := func() {
		for i := 0; i < len(Wishes); i++ {
			w.AddItem("Whish from "+Wishes[i].FROM+" to "+Wishes[i].TO, Wishes[i].TIMESTAMP+" "+Wishes[i].WISH, rune(alphabet[i]), nil)
		}
	}
	generate_wishlist()

	button_grid := tview.NewGrid().SetRows(3, 3, 3).
		AddItem(tview.NewButton("Close App").SetSelectedFunc(func() {
			app.Stop()
		}), 0, 0, 1, 1, 5, 5, true).
		AddItem(tview.NewButton("Refresh").SetSelectedFunc(func() {
			w.Clear()
			get_wishes()
			generate_wishlist()
		}), 1, 0, 1, 1, 5, 5, false).
		AddItem(tview.NewButton("Write Wish").SetSelectedFunc(func() {
			app.Stop()
			print_wish_form()
		}), 2, 0, 1, 1, 5, 5, false)

	grid := tview.NewGrid().
		SetRows(5, 0).
		SetColumns(20, 0).
		SetBorders(true).
		AddItem(header, 0, 0, 1, 2, 0, 0, false).
		AddItem(button_grid, 1, 0, 1, 1, 0, 100, false).
		AddItem(w, 1, 1, 1, 1, 0, 100, false)

	if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func get_wishes() {

	Wishes = nil
	db, _ := sql.Open("sqlite", os.Getenv(env_database))
	defer db.Close()
	rows, err := db.Query("Select * from Wishes ORDER by timestamp desc limit 20;")
	if err != nil {
		fmt.Println(err)
	}
	i := 0
	for rows.Next() {
		var wi Wish
		if err := rows.Scan(&wi.ID, &wi.FROM, &wi.TO, &wi.WISH, &wi.TIMESTAMP, &wi.BROADCAST); err != nil {
		}
		Wishes = append(Wishes, Wish{ID: wi.ID, FROM: wi.FROM, TO: wi.TO, WISH: wi.WISH, TIMESTAMP: wi.TIMESTAMP, BROADCAST: wi.BROADCAST})
		i++
	}

	Users = nil
	rows, err = db.Query("Select * from Users")
	if err != nil {
		fmt.Println(err)
	}
	i = 0
	for rows.Next() {
		var us string
		if err := rows.Scan(&us); err != nil {
		}
		Users = append(Users, us)
		i++
	}

}

func print_wish_form() {
	abort := false
	app := tview.NewApplication()
	form := tview.NewForm().
		AddTextView("Wish from: ", os.Getenv(env_wishlist_user), 0, 1, false, false).
		AddDropDown("Wish to: ", Users, 0, func(option string, optionIndex int) { os.Setenv(env_wishlist_to, option) }).
		AddCheckbox("Broadcast to all users", false, func(option bool) {
			if option {
				os.Setenv(env_broadcast, "0")
			} else {
				os.Setenv(env_broadcast, "1")
			}
		}).
		AddTextArea("Wish: ", "", 30, 5, 150, func(text string) { os.Setenv(env_wishlist_wish, text) }).
		AddButton("Save", func() {
			app.Stop()
		}).
		AddButton("Cancel", func() {
			app.Stop()
			abort = true
		})
	form.SetBorder(false).SetTitleAlign(tview.AlignLeft)
	if err := app.SetRoot(form, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}
	if !abort {
		db, err := sql.Open("sqlite", os.Getenv(env_database))
		if err != nil {
			fmt.Print(err)
		}
		_, err = db.Exec("INSERT INTO Wishes ('from', 'to', 'wish', 'broadcast') VALUES (?,?,?,?);", os.Getenv(env_wishlist_user), os.Getenv(env_wishlist_to), os.Getenv(env_wishlist_wish), os.Getenv(env_broadcast))
		if err != nil {
			fmt.Print(err)
		}
		db.Close()
	}
	print_main_window()
}

func sceduler_get_wishes() {

	for {
		if len(Wishes) > 0 {
			break
		}
		time.Sleep(refresh_rate * time.Second)
		get_wishes()
	}

	for {
		a := Wishes[0].TIMESTAMP
		time.Sleep(refresh_rate * time.Second)
		get_wishes()
		b := Wishes[0].TIMESTAMP
		if strings.Compare(a, b) != 0 {
			walk.MsgBox(nil, "Wunschliste", Wishes[0].to_string(), walk.MsgBoxIconInformation)
		}
	}
}
