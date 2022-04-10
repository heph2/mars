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
	useSyslog   = flag.Bool("log", true, "Whether enable syslog")
	addressPtr  = flag.String("addr", "0.0.0.0", "Address where mars should listen")
	portPtr     = flag.String("port", "8080", "Port where mars should listen")
	certFilePtr = flag.String("cert", "", "TLS Cert file path")
	keyFilePtr  = flag.String("key", "", "TLS Key file path")
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

	stateFile, err := os.OpenFile(*stateDirPtr+req.URL.Path,
		os.O_RDWR, 0666)
	if err != nil {
		log.Println(err)
	}

	fs := handlers.Filesystem{
		StateFile: stateFile,
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
	case "LOCK":
		fs.LockStateFile(res)
		return
	case "UNLOCK":
		fs.UnlockStateFile(res)
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

	log.Println(*certFilePtr, *keyFilePtr)

	if *certFilePtr != "" && *keyFilePtr != "" {
		if err := http.ListenAndServeTLS(*addressPtr+":"+*portPtr,
			*certFilePtr, *keyFilePtr, logRequest(router)); err != nil {
			log.Println(err)
		}
	} else {
		if err := http.ListenAndServe(*addressPtr+":"+*portPtr,
			logRequest(router)); err != nil {
			log.Println(err)
		}
	}
}
