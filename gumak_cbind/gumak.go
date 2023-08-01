package main

/*
enum GumakKey {
	GumakKeyShift,
	GumakKeyZ,
	GumakKeyX,
	GumakKeyC,
	GumakKeyV,
	GumakKeyA,
	GumakKeyS,
	GumakKeyD,
	GumakKeyF,
	GumakKeyG,
	GumakKeyQ,
	GumakKeyW,
	GumakKeyE,
	GumakKeyR,
	GumakKeyT,
	GumakKey1,
	GumakKey2,
	GumakKey3,
	GumakKey4,
	GumakKey5,
	GumakKey0,
	GumakKey9,
	GumakKey8,
	GumakKey7,
	GumakKey6,
	GumakKeyP,
	GumakKeyO,
	GumakKeyI,
	GumakKeyU,
	GumakKeyY,
	GumakKeyEnter,
	GumakKeyL,
	GumakKeyK,
	GumakKeyJ,
	GumakKeyH,
	GumakKeySpace,
	GumakKeySymbolShift,
	GumakKeyM,
	GumakKeyN,
	GumakKeyB,
};
*/
import (
	"C"
)
import (
	"fmt"
	"mutex/gumak"
	"unsafe"

	"golang.org/x/sys/windows"
)

type Instance struct {
	gumak *gumak.Gumak
}

var instances = map[int]Instance{}
var instanceId int

func WCharPtrToString(p *C.wchar_t) string {
	return windows.UTF16PtrToString((*uint16)(unsafe.Pointer(p)))
}

//export GumakResolution
func GumakResolution() (iw, ih, ow, oh int) {
	iw, ih = gumak.InnerResolution()
	ow, oh = gumak.OuterResolution()
	return
}

//export GumakCreate
func GumakCreate(is48K bool, freq, samples int, romPath string) int {
	fmt.Printf("Loading roms from '%s'", romPath)

	var inst Instance
	var err error

	inst.gumak, err = gumak.CreateNew(is48K, freq, samples, romPath)
	if err != nil {
		fmt.Printf("Error '%s'", err)
		return 0
	}

	instanceId++
	instances[instanceId] = inst
	return instanceId
}

//export GumakCreateWChar
func GumakCreateWChar(is48K bool, freq, samples int, p *C.wchar_t) int {
	return GumakCreate(is48K, freq, samples, WCharPtrToString(p))
}

//export GumakDestroy
func GumakDestroy(id int) {
	delete(instances, id)
}

//export GumakUpdateFrame
func GumakUpdateFrame(id int) bool {
	return instances[id].gumak.Update()
}

//export GumakReset
func GumakReset(id int) {
	instances[id].gumak.Restart()
}

//export GumakDisplayDataRGB
func GumakDisplayDataRGB(id int, data []byte, pitch int) {
	instances[id].gumak.GetDisplayDataRGB(data, pitch)
}

//export GumakDisplayDataBGR
func GumakDisplayDataBGR(id int, data []byte, pitch int) {
	instances[id].gumak.GetDisplayDataBGR(data, pitch)
}

//export GumakBackgroundColor
func GumakBackgroundColor(id int) uint32 {
	r, g, b := instances[id].gumak.GetBackgroundColor()
	return (uint32(r) << 24) | (uint32(g) << 16) | (uint32(b) << 8) | 0xff
}

//export GumakHandleKey
func GumakHandleKey(id int, key int, down bool) {
	instances[id].gumak.HandleKey(gumak.Key(key), down)
}

//export GumakLoadSnapshot
func GumakLoadSnapshot(id int, file string) bool {
	err := instances[id].gumak.LoadSnapshot(file)
	if err != nil {
		fmt.Printf("Error loading snapshot: %s", err)
		return false
	}
	return true
}

//export  GumakLoadSnapshotWChar
func GumakLoadSnapshotWChar(id int, p *C.wchar_t) bool {
	return GumakLoadSnapshot(id, WCharPtrToString(p))
}

//export GumakPlayTape
func GumakPlayTape(id int, file string) bool {
	err := instances[id].gumak.PlayTape(file)
	if err != nil {
		fmt.Printf("Error loading tape: %s", err)
		return false
	}
	return true
}

//export GumakPlayTapeWChar
func GumakPlayTapeWChar(id int, p *C.wchar_t) bool {
	return GumakPlayTape(id, WCharPtrToString(p))
}

//export GumakFirstActiveInstance
func GumakFirstActiveInstance() int {
	for k := range instances {
		return k
	}

	return 0
}

//export GumakAudioBuffer
func GumakAudioBuffer(id int, bufferId int, buffer unsafe.Pointer, length int) {
	b := unsafe.Slice((*uint8)(buffer), length)

	instances[id].gumak.SetAudioBuffer(bufferId, b)
}

//export GumakAudioReady
func GumakAudioReady(id int) bool {
	return instances[id].gumak.AudioReady()
}

func main() {}
