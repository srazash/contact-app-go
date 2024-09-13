package archiver

type Archiver struct {
	Status   string
	Prograss float64
	Archive  string
}

func (a *Archiver) GetStatus() string {
	return a.Status
}

func (a *Archiver) GetProgress() float64 {
	return a.Prograss
}

func (a *Archiver) Run() {

}

func (a *Archiver) Reset() {

}

func (a *Archiver) ArchiveFile() string {
	return a.Archive
}

func (a *Archiver) Get() {

}
