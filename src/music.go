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

	"github.com/wtolson/go-taglib"
)

// TODO:
// - Add Debug Prints (instead of just regular Print, use Debug Prints, that will only run when debug flag is set to true)

const MUSIC_DIR string = "../data/music/"
const sr beep.SampleRate = 44100

type MusicFile struct {
	Path      string
	Extension string
}

func Init() {
	speaker.Init(sr, sr.N(time.Second))
}

func Play() {
	musicQueue, err := LoadMusicQueue()
	if err != nil {
		log.Fatal(err)
	}
	ShuffleQueue(&musicQueue)

	for true {
		queue_length := len(musicQueue)

		if queue_length == 0 {
			musicQueue, err = LoadMusicQueue()
			if err != nil {
				log.Print(err)
			}
			ShuffleQueue(&musicQueue)

			queue_length = len(musicQueue)
		}

		var streamer beep.StreamSeekCloser
		var format beep.Format
		currentFile, currentFileExtension := getFromQueue(musicQueue)

		fm := getMetadata(musicQueue)
		fmt.Println(fm.Title())
		fmt.Println(fm.Album())
		fmt.Println(fm.Artist())
		fmt.Println(fm.Genre())
		fmt.Println(fm.Year())

		switch currentFileExtension {
		case ".flac":
			streamer, format = decodeFLAC(currentFile)
		case ".mp3":
			streamer, format = decodeMP3(currentFile)
		}
		defer streamer.Close()

		resampled := beep.Resample(4, format.SampleRate, sr, streamer)

		done := make(chan bool)
		speaker.Play(beep.Seq(resampled, beep.Callback(func() {
			done <- true
		})))

		<-done

		musicQueue = musicQueue[1:queue_length]

	}
}

func getMetadata(q []MusicFile) (t *taglib.File) {
	f, err := taglib.Read(q[0].Path)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
	return f
}

func decodeFLAC(f *os.File) (streamer beep.StreamSeekCloser, format beep.Format) {
	streamer, format, err := flac.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	return streamer, format
}

func decodeMP3(f *os.File) (streamer beep.StreamSeekCloser, format beep.Format) {
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	return streamer, format
}

func LoadMusicQueue() ([]MusicFile, error) {
	var q []MusicFile

	err := filepath.WalkDir(MUSIC_DIR, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			log.Print(err)
		}
		if !entry.IsDir() {
			switch filepath.Ext(entry.Name()) {
			case ".flac":
				music_file := MusicFile{path, filepath.Ext(entry.Name())}
				q = append(q, music_file)
			case ".mp3":
				music_file := MusicFile{path, filepath.Ext(entry.Name())}
				q = append(q, music_file)
			}
		}

		return nil
	})

	return q, err
}

func getFromQueue(q []MusicFile) (*os.File, string) {
	file_path := q[0].Path
	f, err := os.Open(file_path)
	if err != nil {
		log.Fatal(err)
	}

	return f, q[0].Extension
}

/* Fisher-Yates algorithm */
func ShuffleQueue(q *[]MusicFile) {
	last_index := len(*q) - 1
	for true {
		if last_index <= 0 {
			break
		}

		rand_index := rand.Intn(len(*q))

		buffer := (*q)[rand_index]
		(*q)[rand_index] = (*q)[last_index]
		(*q)[last_index] = buffer

		last_index -= 1
	}

	fmt.Println("Shuffled Queue:", q)
}
