package counter

import (
	"encoding/json"
	"io"
	"os"
)

const COUNTFILE string = "count.json"

var Count int = 0

func Ptr() *int {
	return &Count
}

func Load() {
	dbfile, err := os.Open(COUNTFILE)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}
	defer dbfile.Close()

	file, err := io.ReadAll(dbfile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &Count)
	if err != nil {
		panic(err)
	}

	Count += Count
}

func Save() {
	file, err := os.Create(COUNTFILE)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := json.Marshal(Count)
	if err != nil {
		panic(err)
	}

	_, err = file.Write(data)
	if err != nil {
		panic(err)
	}
}
