package main

import (
	"errors"
	"flag"
	"log"
	"log/syslog"
	"net/http"
	"os"

	"git.mrkeebs.eu/debris/handlers"
)

var (
	stateDirPtr = flag.String("dir", ".", "Terraform state directory path")
	useSyslog   = flag.Bool("log", true, "Enable syslog")
	addressPtr  = flag.String("addr", "0.0.0.0", "Address where mars should listen")
	portPtr     = flag.String("port", "8080", "Port where mars should listen")
)

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s \n", r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
func requestHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		http.NotFound(res, req)
		log.Println("Terraform state not found")
		return
	}

	fs := handlers.Filesystem{
		Filepath: *stateDirPtr + req.URL.Path,
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
		log.Println("HTTP Method not supperted")
	}
}

func main() {
	flag.Parse()

	if *useSyslog {
		logwriter, err := syslog.New(syslog.LOG_NOTICE, "mars")
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(logwriter)
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	}

	err := os.Mkdir(*stateDirPtr, 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Println("Error creating dir", err)
	}

	router := http.NewServeMux()
	router.HandleFunc("/", requestHandler)
	router.HandleFunc("/lock", handler func(ResponseWriter, *Request))

	if err := http.ListenAndServe(*addressPtr+":"+*portPtr,
		logRequest(router)); err != nil {
		log.Println(err)
	}
}
