package main

import (
	"fmt"

	"log"
	"math/rand"
	"time"

	"io/fs"
	"os"
	"path/filepath"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/flac"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

// TODO:
// - Add Debug Prints (instead of just regular Print, use Debug Prints, that will only run when debug flag is set to true)

const MUSIC_DIR string = "../data/music/"
const sr beep.SampleRate = 44100

type MusicFile struct {
	FilePath      string
	FileExtension string
}

func Init() {
	speaker.Init(sr, sr.N(time.Second))
}

func Play() {
	queue, err := LoadMusicFiles()
	Shuffle(&queue)

	for true {
		queue_length := len(queue)

		if queue_length == 0 {
			queue, err = LoadMusicFiles()
			if err != nil {
				log.Print(err)
			}

			Shuffle(&queue)
			queue_length = len(queue)
		}

		var streamer beep.StreamSeekCloser
		var format beep.Format
		f, ext := getFromQueue(queue)
		switch ext {
		case ".flac":
			streamer, format, err = flac.Decode(f)
		case ".mp3":
			streamer, format, err = mp3.Decode(f)
		}
		defer streamer.Close()

		resampled := beep.Resample(4, format.SampleRate, sr, streamer)

		done := make(chan bool)
		speaker.Play(beep.Seq(resampled, beep.Callback(func() {
			done <- true
		})))

		<-done

		queue = queue[1:queue_length]

	}

}

func decodeFlac(f *os.File) (streamer beep.StreamSeekCloser, format beep.Format) {
	streamer, format, err := flac.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	return streamer, format
}

func LoadMusicFiles() ([]MusicFile, error) {
	var files []MusicFile

	err := filepath.WalkDir(MUSIC_DIR, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			log.Print(err)
		}
		if !entry.IsDir() {
			switch filepath.Ext(entry.Name()) {
			case ".flac":
				music_file := MusicFile{path, filepath.Ext(entry.Name())}
				files = append(files, music_file)
			case ".mp3":
				music_file := MusicFile{path, filepath.Ext(entry.Name())}
				files = append(files, music_file)
			}
		}

		return nil
	})

	return files, err
}

func getFromQueue(files []MusicFile) (*os.File, string) {
	file_path := files[0].FilePath
	f, err := os.Open(file_path)
	if err != nil {
		log.Fatal(err)
	}

	return f, files[0].FileExtension
}

/* Fisher-Yates algorithm */
func Shuffle(files *[]MusicFile) {
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
