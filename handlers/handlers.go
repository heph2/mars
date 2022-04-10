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
	StateFile *os.File
}

func (fh Filesystem) LockStateFile(res http.ResponseWriter) {
	if err := syscall.Flock(int(fh.StateFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		res.WriteHeader(423)
		log.Println("Get Exclusive Lock failed", err)
	}

	log.Println("Lock acquired on", fh.StateFile.Name())
}

func (fh Filesystem) UnlockStateFile(res http.ResponseWriter) {
	if err := syscall.Flock(int(fh.StateFile.Fd()), syscall.LOCK_UN); err != nil {
		log.Println("Error while unlocking file", err)
	}

	log.Println("Lock released on", fh.StateFile.Name())
}

func (fh Filesystem) GetStateFile(res http.ResponseWriter) {
	io.Copy(res, fh.StateFile)
}

func (fh Filesystem) UpdateStateFile(req *http.Request, res http.ResponseWriter) {
	if _, err := io.Copy(fh.StateFile, req.Body); err != nil {
		log.Println("Error streaming data", err)
	}

	log.Println(fh.StateFile.Name())
	res.WriteHeader(200)
}

func (fh Filesystem) DeleteStateFile(res http.ResponseWriter) {
	log.Println("Removing", fh.StateFile.Name())
	if os.RemoveAll(fh.StateFile.Name()) != nil {
		log.Println("Error while deleting statefile")
	}

	res.WriteHeader(200)
}
