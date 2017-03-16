package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/VojtechVitek/go-trello"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Board struct {
	Lists []trello.List
}

type List struct {
	Name  string
	Cards []trello.Card
}

var (
	trelloAppKey  string
	trelloToken   string
	trelloBoardId string

	trelloClient *trello.Client
	trelloBoard  *trello.Board
)

func main() {

	// set log format
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Please provide a valid port number (e.g. 8080)")
	}

	trelloAppKey = os.Getenv("TRELLO_APP_KEY")
	trelloToken = os.Getenv("TRELLO_TOKEN")
	trelloBoardId = os.Getenv("TRELLO_BOARD_ID")

	// New Trello Client
	trelloClient, err = trello.NewAuthClient(trelloAppKey, &trelloToken)
	if err != nil {
		log.Fatalf("Could not connect to Trello, err: %s", err)
		os.Exit(1)
	}

	trelloBoard, err = trelloClient.Board(trelloBoardId)
	if err != nil {
		log.Fatalf("Could not get Trello board %s, err: %s", trelloBoardId, err)
		os.Exit(1)
	}

	log.Info("Up & running...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	})

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	log.Info("Received request")

	tariffBoard := Board{}

	// @trello Board Lists
	lists, err := trelloBoard.Lists()
	if err != nil {
		log.Errorf("Failed getting Lists of Trello board, err: %s", err)
		os.Exit(1)
	}

	// TODO: make slice dynamic
	// Get first 3 lists
	for _, list := range lists[1:4] {
		tariffBoard.Lists = append(tariffBoard.Lists, list)
	}

	templatePath := filepath.Join("tmpl", "board.html")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Errorf("Failed to parse template, err: %s", err)
		return 500, err
	}
	err = tmpl.Execute(w, tariffBoard)
	if err != nil {
		log.Errorf("Failed to apply template, err: %s", err)
		return 500, err
	}
}
