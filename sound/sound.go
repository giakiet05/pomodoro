package sound

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var initialized = false

// InitSpeaker should be called once in main()
func InitSpeaker() {
	if initialized {
		return
	}
	// init với sample rate chuẩn, ví dụ 44100Hz
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10))
	initialized = true
}

// PlaySound plays an mp3 file (first 5s)
func PlaySound(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Println("Failed to open sound file:", err)
		return
	}
	defer f.Close()

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Println("Failed to decode mp3:", err)
		return
	}
	defer streamer.Close()

	// Lấy 5s đầu tiên
	samples := format.SampleRate.N(3 * time.Second)
	limited := beep.Take(samples, streamer)

	done := make(chan bool)
	speaker.Play(beep.Seq(limited, beep.Callback(func() {
		done <- true
	})))

	<-done
}
