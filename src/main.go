package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/flac"
	"github.com/gopxl/beep/speaker"
)

func main() {
	// setting sample rate and buffer size(buffer size = 1/10(in my case it is 1) of a second; smaller buffer = lower latency; bigger buffer = better stability)
	// Generally, you only want to call it once at the beginning of your program!
	speaker.Init(sr, sr.N(time.Second))

	files, err := getMusicFiles()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(files)
	shuffle(&files)

	for true {
		files_len := len(files)

		if files_len == 0 {
			files, err = getMusicFiles()
			fmt.Println(files)
			shuffle(&files)

			files_len = len(files)
		}
		f := getMusicFile(files)

		streamer, format, err := flac.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
		defer streamer.Close()

		// our declared sample rate(if music has higher or lower sample Rate, it will be resampled to play at the same speed)
		fmt.Println(sr)
		resampled := beep.Resample(4, format.SampleRate, sr, streamer)

		done := make(chan bool)
		speaker.Play(beep.Seq(resampled, beep.Callback(func() {
			done <- true
		})))

		<-done

		files = files[1:files_len]
	}
}
