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
	"github.com/TV4/graceful"
)

type Board struct {
	Lists []trello.List
}

type List struct {
	Name  string
	Cards []trello.Card
}

type Server struct{}

var (
	port int

	trelloAppKey  string
	trelloToken   string
	trelloBoardId string

	trelloClient *trello.Client
	trelloBoard  *trello.Board

	trelloStartColumn int
	trelloStopColumn  int
)

func init() {
	var err error

	// set log format
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	port, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Please provide a valid port number (e.g. 8080)")
	}

	trelloAppKey = os.Getenv("TRELLO_APP_KEY")
	trelloToken = os.Getenv("TRELLO_TOKEN")

	if trelloAppKey == "" || trelloToken == "" {
		log.Fatal("Please provide trello credentials")
	}
}

func main() {

	var err error

	// New Trello Client
	trelloClient, err = trello.NewAuthClient(trelloAppKey, &trelloToken)
	if err != nil {
		log.Errorf("Could not connect to Trello, err: %s", err)
		os.Exit(1)
	}

	log.Info("Up & running...")

	graceful.LogListenAndServe(&http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: &Server{},
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/favicon.ico" {
		return
	}

	var err error

	log.Info("Received request")

	trelloBoardId = r.URL.Query().Get("trelloBoardId")
	if trelloBoardId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Please provide a Trello Board ID"))
		log.Error("Trello Board ID missing")
		return

	} else {
		trelloBoard, err = trelloClient.Board(trelloBoardId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Could not get Trello board"))
			log.Errorf("Could not get Trello board %s, err: %s", trelloBoardId, err)
			return
		}
	}

	startColumn := r.URL.Query().Get("trelloStartColumn")
	if startColumn == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Please provide a trello start column (trelloStartColumn)"))
		log.Error("Trello start column is missing")
		return

	} else {
		trelloStartColumn, err = strconv.Atoi(startColumn)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Please provide a valid trello start column (e.g. 1)"))
			log.Error("No valid trello start column (e.g. 1) was set as get parameter")
		}
	}

	stopColumn := r.URL.Query().Get("trelloStopColumn")
	if stopColumn == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Please provide a trello stop column (trelloStopColumn)"))
		log.Error("Trello stop column is missing")
		return

	} else {
		trelloStopColumn, err = strconv.Atoi(stopColumn)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Please provide a valid trello stop column (e.g. 3)"))
			log.Error("No valid trello stop column (e.g. 3) was set as get parameter")
		}
	}

	board := Board{}

	// @trello Board Lists
	allLists, err := trelloBoard.Lists()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Failed getting Lists for Trello board"))
	}

	for _, list := range allLists[trelloStartColumn:trelloStopColumn] {
		board.Lists = append(board.Lists, list)
	}

	templatePath := filepath.Join("tmpl", "board.html")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Failed parsing template"))
	}
	err = tmpl.Execute(w, board)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Failed applying template"))
	}
}