package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"syscall"
)

type Handler interface {
	getStateFile(http.ResponseWriter)
	updateStateFile(http.Request, http.ResponseWriter)
	deleteStateFile(http.ResponseWriter)
	lockStateFile(http.ResponseWriter)
	unlockStateFile(http.ResponseWriter)
}

type Filesystem struct {
	Filepath string
}

func (fh Filesystem) lockStateFile(res http.ResponseWriter) {

}

func (fh Filesystem) unlockStateFile(res http.ResponseWriter) {

}

func (fh Filesystem) GetStateFile(res http.ResponseWriter) {
	file, err := os.Open(fh.Filepath)
	if err != nil {
		log.Println("Error while open file", err)
	}
	defer file.Close()
	res.WriteHeader(200)

	io.Copy(res, file)
}

func (fh Filesystem) UpdateStateFile(req *http.Request, res http.ResponseWriter) {
	f, err := os.Create(fh.Filepath)
	if err != nil {
		log.Println("create file failed", err)
	}
	defer f.Close()

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		log.Println("Add Exclusive Lock failed", err)
	}

	if _, err := io.Copy(f, req.Body); err != nil {
		log.Println("Error streaming data", err)
	}

	res.WriteHeader(200)
}

func (fh Filesystem) DeleteStateFile(res http.ResponseWriter) {
	if os.RemoveAll(fh.Filepath) != nil {
		log.Println("Error while deleting statefile")
	}

	res.WriteHeader(200)
}
