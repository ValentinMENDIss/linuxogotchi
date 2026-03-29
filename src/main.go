package main

import (
	"log"
	"math/rand"

	"io/fs"
	"os"
	"path/filepath"

	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/flac"
	"github.com/gopxl/beep/speaker"

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

func main() {
	files, err := getMusicFiles()

	fmt.Println(files)

	shuffle(&files)

	f := getMusicFile(files)

	streamer, format, err := flac.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	// setting sample rate and buffer size(buffer size = 1/10(in my case it is 1) of a second; smaller buffer = lower latency; bigger buffer = better stability)
	// Generally, you only want to call it once at the beginning of your program!
	speaker.Init(sr, sr.N(time.Second))
	// our declared sample rate(if music has higher or lower sample Rate, it will be resampled to play at the same speed)
	fmt.Println(sr)
	resampled := beep.Resample(4, format.SampleRate, sr, streamer)

	done := make(chan bool)
	speaker.Play(beep.Seq(resampled, beep.Callback(func() {
		done <- true
	})))

	<-done
}
