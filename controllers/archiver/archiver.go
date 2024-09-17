package archiver

import (
	"contactapp/models/contact"
	"encoding/json"
)

const (
	Waiting = iota
	Running
	Complete
	Errored
)

type Archiver struct {
	Status   int
	Progress float64
	Archive  []byte
}

func (a *Archiver) GetStatus() int {
	return a.Status
}

func (a *Archiver) GetProgress() float64 {
	return a.Progress
}

func (a *Archiver) Run() {
	go a.runner()
}

func (a *Archiver) runner() {
	if a.Status == Waiting {
		a.Status = Running
		a.Progress = 0.33

		data, err := json.Marshal(contact.DB)
		if err != nil {
			a.Status = Errored
			a.Progress = 0.0
			return
		}

		a.Progress = 0.67
		a.Archive = data

		a.Progress = 1.0
		a.Status = Complete
	}
}

func (a *Archiver) Reset() {
	a.Status = Waiting
}

func (a *Archiver) ArchiveFile() []byte {
	return a.Archive
}

func New() Archiver {
	return Archiver{
		Status:   Waiting,
		Progress: 0.0,
		Archive:  []byte{},
	}
}
