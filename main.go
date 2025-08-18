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

	_ "github.com/mattn/go-sqlite3"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

type model struct {
	choices  []string         // items on the to-do list
	cursor   int              // which to-do list item our cursor is pointing at
	selected map[int]struct{} // which to-do items are selected
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "What should we buy at the market?\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func initialModel() model {
	return model{
		// Our to-do list is a grocery list
		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func main() {

	var database_location string

	argus := os.Args
	if len(argus) >= 2 {
		switch argus[1] {
		case "--init":
			//setup database ----------------------------------------------
			if len(argus) > 2 {
				database_location = argus[2]
			} else {
				database_location = "./data.db"
			}

			os.Create(database_location)
			fmt.Println("\t + databasefile created")
			db, err := sql.Open("sqlite3", database_location)
			if err != nil {
				fmt.Print(err)
			}
			db.Exec("CREATE TABLE Test(id INTEGER PRIMARY KEY, t1 TEXT);")
			fmt.Println("\t + database initiated")
			db.Close()
			// STOP
			//setup settings.env ------------------------------------------
			settings_location := "settings.env"
			if _, err := os.Stat(settings_location); errors.Is(err, os.ErrNotExist) {
				os.Create(settings_location)
				fmt.Println("\t" + settings_location + " created")
			} else {
				fmt.Println("\t +" + settings_location + " is existing")
			}

			//write settings ---------------------------------------------
			data := []byte("DATABASE=" + database_location + "\n" +
				"SETTINGS=" + settings_location + "\n")
			os.WriteFile(settings_location, data, 0644)

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
