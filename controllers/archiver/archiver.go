package archiver

const (
	Waiting = iota
	Running
	Complete
)

type Archiver struct {
	Status   int
	Progress float64
	Archive  chan string
}

func (a *Archiver) GetStatus() int {
	return a.Status
}

func (a *Archiver) GetProgress() float64 {
	return a.Progress
}

func (a *Archiver) Run() {

}

func (a *Archiver) Reset() {
	a.Status = Waiting
}

func (a *Archiver) ArchiveFile() string {
	s := <-a.Archive
	return s
}

func Get() Archiver {
	return Archiver{
		Status:   Waiting,
		Progress: 0.0,
		Archive:  make(chan string, 1),
	}
}
