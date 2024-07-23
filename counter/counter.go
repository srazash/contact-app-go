package counter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const COUNTFILE string = "count.json"

var Count int = 0

func Increment() {
	Count++
	Save()
}

func PaddedCount() string {
	return fmt.Sprintf("%06d", Count)
}

func Load() {
	countfile, err := os.Open(COUNTFILE)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}
	defer countfile.Close()

	file, err := io.ReadAll(countfile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &Count)
	if err != nil {
		panic(err)
	}
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
