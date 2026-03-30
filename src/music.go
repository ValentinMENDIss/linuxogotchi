package main

import (
	"log"
	"math/rand"

	"io/fs"
	"os"
	"path/filepath"

	"github.com/gopxl/beep"

	"fmt"
)

// TODO:
// - Add Debug Prints (instead of just regular Print, use Debug Prints, that will only run when debug flag is set to true)

const MUSIC_DIR string = "../data/music/"
const sr beep.SampleRate = 44100

func getMusicFiles() ([]string, error) {
	var files []string

	err := filepath.WalkDir(MUSIC_DIR, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			log.Print(err)
		}
		if !entry.IsDir() {
			switch filepath.Ext(entry.Name()) {
			case ".flac":
				files = append(files, path)
			case ".mp3":
				files = append(files, path)
			}
		}

		return nil
	})

	return files, err
}

func getMusicFile(files []string) *os.File {
	file_path := files[0]
	f, err := os.Open(file_path)
	if err != nil {
		log.Fatal(err)
	}

	return f
}

/* Fisher-Yates algorithm */
func shuffle(files *[]string) {
	last_index := len(*files) - 1
	for true {
		if last_index <= 0 {
			break
		}

		rand_index := rand.Intn(len(*files))

		buffer := (*files)[rand_index]
		(*files)[rand_index] = (*files)[last_index]
		(*files)[last_index] = buffer

		last_index -= 1
	}

	fmt.Println("Shuffled:", files)
}
