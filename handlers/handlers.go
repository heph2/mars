package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"syscall"
	"encoding/json"
)

type Handler interface {
	getStateFile(http.ResponseWriter)
	updateStateFile(http.Request, http.ResponseWriter)
	deleteStateFile(http.ResponseWriter)
	lockStateFile(http.Request, http.ResponseWriter)
	unlockStateFile(http.ResponseWriter)
}

type Filesystem struct {
	StateFile string
}

type State struct {
	Locked bool `json:"locked"`
	State interface{} `json:"state"`
}

var state State

func (fh Filesystem) LockStateFile(req *http.Request, res http.ResponseWriter) {
	b, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return
	}
	
	var data interface{}
	err = json.Unmarshal(b, &data)
	state.Locked = true
	log.Println("Current state", state.Locked)
	state.State = data
}

func (fh Filesystem) UnlockStateFile(res http.ResponseWriter) {
	state.Locked = false
	log.Println("Current state", state.Locked)
}

func (fh Filesystem) GetStateFile(res http.ResponseWriter) {
	f, err := os.OpenFile(fh.StateFile, os.O_RDWR, 0666)
	if err != nil {
		log.Println(err)
		return
	}

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_SH|syscall.LOCK_NB); err != nil {
		log.Println("Error while acquiring _non_ blocking lock")
		return
	}

	if _, err := io.Copy(res, f); err != nil {
		log.Println("Error while streming data", err)
		return
	}

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
		log.Println("Unlock share lock failed", err)
		return
	}
}

func (fh Filesystem) UpdateStateFile(req *http.Request, res http.ResponseWriter) {
	f, err := os.OpenFile(fh.StateFile, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Println(err)
		return
	}

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		log.Println("Error while acquiring _non_ blocking lock")
		return
	}

	// CHECK STRUCT BEFORE COPYING //
	if !state.Locked {
		if _, err := io.Copy(f, req.Body); err != nil {
			log.Println("Error streaming data", err)
			return
		}
	} else {
		log.Println("File is locked")
	}

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
		log.Println("Unlock share lock failed", err)
		return
	}

	res.WriteHeader(200)
}

func (fh Filesystem) DeleteStateFile(res http.ResponseWriter) {
	if os.RemoveAll(fh.StateFile) != nil {
		log.Println("Error while deleting statefile")
		return
	}

	res.WriteHeader(200)
}
