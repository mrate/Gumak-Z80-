package main

import (
	"flag"
	"mutex/gumak"
	"mutex/gumak/log"
	"mutex/gumak_sdl/host"
	"os"
	"runtime"
	"runtime/pprof"
)

func main() {
	// Flags.
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	var filtering = flag.Int("filtering", 0, "upscale filtering (0=nearest, 1=linear, 2=best, 3=hq2x, 4=hq3x)")
	var scale = flag.Float64("scale", 4, "screen scale multiplicator")
	var machine = flag.String("machine", "128", "machine (48=48K, 128=128K)")
	var sound = flag.Bool("sound", true, "turn on sound")
	var rom = flag.String("rom", "", "rom to load on startup")

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Error("Failed to create pprof")
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	iw, ih := gumak.InnerResolution()
	ow, oh := gumak.OuterResolution()

	freq, samples := 41800, 512

	gumak, err := gumak.CreateNew(*machine == "48", freq)
	if err != nil {
		panic(err)
	}

	if len(*rom) > 0 {
		err := gumak.LoadSnapshot(*rom, nil)
		if err != nil {
			panic(err)
		}
	}

	runtime.LockOSThread()

	// Host platform
	platform := new(host.Platform)
	platform.Init(ow, oh, *scale)
	defer platform.Destroy()

	gfx := platform.CreateGfx(host.Filtering(*filtering), iw, ih, ow, oh)
	snd := platform.CreateSound(gumak, freq, samples)
	snd.TurnOnOff(*sound)

	for {
		if platform.Closed() {
			break
		}

		platform.WaitFrame()
		platform.Update(gumak)
		gfx.Draw(gumak)
	}
}
