package archiver

import (
	"contactapp/models/contact"
	"encoding/json"
	"time"
)

const (
	Waiting = iota
	Running
	Complete
	Errored
)

type Archiver struct {
	Status   int
	Progress int
	Archive  []byte
}

func (a *Archiver) GetStatus() int {
	return a.Status
}

func (a *Archiver) GetProgress() int {
	return a.Progress
}

func (a *Archiver) Run() {
	go a.runner()
}

func (a *Archiver) runner() {
	if a.Status == Waiting {
		a.Status = Running
		a.Progress = 33

		time.Sleep(1 * time.Second)

		data, err := json.Marshal(contact.DB)
		if err != nil {
			a.Status = Errored
			a.Progress = 0
			return
		}

		time.Sleep(1 * time.Second)

		a.Progress = 67
		a.Archive = data

		time.Sleep(1 * time.Second)

		a.Progress = 100
		a.Status = Complete
	}
}

func (a *Archiver) Reset() {
	a.Status = Waiting
}

func (a *Archiver) ArchiveFile() []byte {
	return a.Archive
}

func Get() Archiver {
	return Archiver{
		Status:   Waiting,
		Progress: 0.0,
		Archive:  []byte{},
	}
}
