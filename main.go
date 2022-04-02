package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gofrs/flock"
)

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

	switch req.Method {
	case "GET":
		if err := getStateFile(res); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}
	case "POST":
		if err := updateStateFile(res); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}
	case "DELETE":
		if err := deleteStateFile(res); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}
	default:
		fmt.Printf("Method not supported\n")
	}
}

func getStateFile(res http.ResponseWriter) error {
	fileLock := flock.New("/tmp/terraform-state")
	locked, err := fileLock.TryLock()

	if err != nil {
		log.Println("Error while locking", err)
	}

	if locked {
		fh, err := os.Open(fileLock.Path())
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return errors.New("Cannot open statefile")
		}
		defer fh.Close()
		res.WriteHeader(200)

		io.Copy(res, fh)
		fileLock.Unlock()
	}
	return nil
}

func updateStateFile(res http.ResponseWriter) error {
	if err := os.Mkdir("/var/lib/debris", 0755); err != nil {
		log.Fatal(err)
	}
	return nil
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
