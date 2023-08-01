package host

// typedef unsigned char Uint8;
// void BeeperCB(void *userdata, Uint8 *stream, int len);
import "C"

import (
	"reflect"
	"sync"
	"unsafe"

	"mutex/gumak"
	"mutex/gumak/log"

	"github.com/veandco/go-sdl2/sdl"
)

type Sound struct {
	gumak      *gumak.Gumak
	isOn       bool
	soundMutex *sync.Mutex
	frameReady *chan bool
	vram       []byte
}

var sound *Sound

var soundTime float64

//export BeeperCB
func BeeperCB(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	n := int(length)
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.Uint8)(unsafe.Pointer(&hdr))

	sound.soundMutex.Lock()
	defer sound.soundMutex.Unlock()

	for i := 0; i < n; i++ {
		for !sound.gumak.AudioSampleReady() {
			if sound.gumak.Tick() {
				sound.gumak.CopyVRam(sound.vram)

				select {
				case *sound.frameReady <- true:
				default:
				}
			}
		}

		buf[i] = C.Uint8(sound.gumak.PopAudioSample())

		/*
				buf[i] = C.Uint8(32 * (0.5 + 0.5*math.Sin(soundTime*440*2*math.Pi)))

			soundTime += 1. / 44100.
		*/
	}
}

// Sound
func (s *Sound) Init(g *gumak.Gumak, freq int, samples int, mutex *sync.Mutex, frameReady *chan bool, vram []byte) {
	s.gumak = g
	s.isOn = true
	s.soundMutex = mutex
	s.frameReady = frameReady
	s.vram = vram

	sound = s

	spec := &sdl.AudioSpec{
		Freq:     int32(freq),
		Format:   sdl.AUDIO_U8,
		Channels: 1,
		Samples:  uint16(samples),
		Callback: sdl.AudioCallback(C.BeeperCB),
	}

	if err := sdl.OpenAudio(spec, nil); err != nil {
		log.Error("Error opening audio: %s", err)
		return
	}

	sdl.PauseAudio(s.isOn)
}

func (s *Sound) Destroy() {
	sdl.PauseAudio(true)
	sdl.CloseAudio()
}

func (s *Sound) TurnOnOff(on bool) {
	s.isOn = on
	sdl.PauseAudio(!on)
}

func (s *Sound) IsOn() bool {
	return s.isOn
}

//func (s *Sound) Test() {
// Sound check, one, two, one, two...
// Nota A4: 440Hz = 1/440 = 0.00227272727272 (/2 -> half up, half down)
//x := 0
//t := 0.0
//for {
//	e := event{value: uint8(x * 255), time: t}
//	x = 1 - x
//	t += 0.00227272727272 / 4
//	eventQueue <- e
//}
//}
