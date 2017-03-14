package main

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	"github.com/VojtechVitek/go-trello"
	"html/template"
	"path/filepath"
	"os"
	"fmt"
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
	trelloAppKey string
	trelloToken string
	trelloBoardId string

	trelloClient *trello.Client
	trelloBoard *trello.Board
)

func main() {

	// set log format
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}

	trelloAppKey = os.Getenv("TRELLO_APP_KEY")
	trelloToken = os.Getenv("TRELLO_TOKEN")
	trelloBoardId = os.Getenv("TRELLO_BOARD_ID")

	// New Trello Client
	trelloClient, err = trello.NewAuthClient(trelloAppKey, &trelloToken)
	if err != nil {
		log.Fatal(err)
	}

	trelloBoard, err = trelloClient.Board(trelloBoardId)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	// TODO: make slice dynamic
	// Get first 3 lists
	for _, list := range lists[1:4] {
		tariffBoard.Lists = append(tariffBoard.Lists, list);
	}


	templatePath := filepath.Join("tmpl", "board.html")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, tariffBoard)
	if err != nil {
		panic(err)
	}
}