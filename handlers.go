package handlers

type handlers interface {
	getStateFile()
	updateStateFile()
	deleteStateFile()
}

type Filesystem struct {
	//
}

func (fh Filesystem) getStateFile() {
	//
}
