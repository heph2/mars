package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"git.mrkeebs.eu/debris/handlers"
)

var filepath string

func serveHttp(address string) error {
	router := http.NewServeMux()
	router.HandleFunc("/", requestHandler)

	if err := http.ListenAndServe(address, router); err != nil {
		return err
	}
	return nil
}

func requestHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		http.NotFound(res, req)
		return
	}

	fs := handlers.Filesystem{
		Filepath: filepath + req.URL.Path,
	}

	switch req.Method {
	case "GET":
		fs.GetStateFile(res)
		return
	case "POST":
		fs.UpdateStateFile(req, res)
		return
	case "DELETE":
		fs.DeleteStateFile(res)
		return
	default:
		fmt.Printf("Method not supported\n")
	}
}

// This init function check if the directory that stores states
// exists. If not create it.
func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Println("Error while get home dir", err)
	}
	filepath = home + "/.debris"
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(filepath, 0755)
		if err != nil {
			log.Println("Error creating dir", err)
		}
	}
}

func main() {

	// how to handling errors
	// how to handling logs
	// how to handling init
	if err := serveHttp("0.0.0.0:8080"); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
