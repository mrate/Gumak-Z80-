package host

import (
	"mutex/gumak"

	"github.com/veandco/go-sdl2/sdl"
)

var KeyMap = map[sdl.Keycode][]gumak.Key{
	sdl.K_LSHIFT: []gumak.Key{gumak.KeyShift},
	sdl.K_z:      []gumak.Key{gumak.KeyZ},
	sdl.K_x:      []gumak.Key{gumak.KeyX},
	sdl.K_c:      []gumak.Key{gumak.KeyC},
	sdl.K_v:      []gumak.Key{gumak.KeyV},

	sdl.K_a: []gumak.Key{gumak.KeyA},
	sdl.K_s: []gumak.Key{gumak.KeyS},
	sdl.K_d: []gumak.Key{gumak.KeyD},
	sdl.K_f: []gumak.Key{gumak.KeyF},
	sdl.K_g: []gumak.Key{gumak.KeyG},

	sdl.K_q: []gumak.Key{gumak.KeyQ},
	sdl.K_w: []gumak.Key{gumak.KeyW},
	sdl.K_e: []gumak.Key{gumak.KeyE},
	sdl.K_r: []gumak.Key{gumak.KeyR},
	sdl.K_t: []gumak.Key{gumak.KeyT},

	sdl.K_1: []gumak.Key{gumak.Key1},
	sdl.K_2: []gumak.Key{gumak.Key2},
	sdl.K_3: []gumak.Key{gumak.Key3},
	sdl.K_4: []gumak.Key{gumak.Key4},
	sdl.K_5: []gumak.Key{gumak.Key5},

	sdl.K_0: []gumak.Key{gumak.Key0},
	sdl.K_9: []gumak.Key{gumak.Key9},
	sdl.K_8: []gumak.Key{gumak.Key8},
	sdl.K_7: []gumak.Key{gumak.Key7},
	sdl.K_6: []gumak.Key{gumak.Key6},

	sdl.K_p: []gumak.Key{gumak.KeyP},
	sdl.K_o: []gumak.Key{gumak.KeyO},
	sdl.K_i: []gumak.Key{gumak.KeyI},
	sdl.K_u: []gumak.Key{gumak.KeyU},
	sdl.K_y: []gumak.Key{gumak.KeyY},

	sdl.K_RETURN: []gumak.Key{gumak.KeyEnter},
	sdl.K_l:      []gumak.Key{gumak.KeyL},
	sdl.K_k:      []gumak.Key{gumak.KeyK},
	sdl.K_j:      []gumak.Key{gumak.KeyJ},
	sdl.K_h:      []gumak.Key{gumak.KeyH},

	sdl.K_SPACE: []gumak.Key{gumak.KeySpace},
	sdl.K_LCTRL: []gumak.Key{gumak.KeySymbolShift},
	sdl.K_m:     []gumak.Key{gumak.KeyM},
	sdl.K_n:     []gumak.Key{gumak.KeyN},
	sdl.K_b:     []gumak.Key{gumak.KeyB},

	sdl.K_RSHIFT:       []gumak.Key{gumak.KeyShift},
	sdl.K_BACKSPACE:    []gumak.Key{gumak.KeyShift, gumak.Key0},
	sdl.K_LEFT:         []gumak.Key{gumak.KeyShift, gumak.Key5},
	sdl.K_RIGHT:        []gumak.Key{gumak.KeyShift, gumak.Key8},
	sdl.K_UP:           []gumak.Key{gumak.KeyShift, gumak.Key7},
	sdl.K_DOWN:         []gumak.Key{gumak.KeyShift, gumak.Key6},
	sdl.K_COMMA:        []gumak.Key{gumak.KeySymbolShift, gumak.KeyN},
	sdl.K_PERIOD:       []gumak.Key{gumak.KeySymbolShift, gumak.KeyM},
	sdl.K_SLASH:        []gumak.Key{gumak.KeySymbolShift, gumak.KeyV},
	sdl.K_EQUALS:       []gumak.Key{gumak.KeySymbolShift, gumak.KeyL},
	sdl.K_KP_PLUS:      []gumak.Key{gumak.KeySymbolShift, gumak.KeyK},
	sdl.K_KP_MINUS:     []gumak.Key{gumak.KeySymbolShift, gumak.KeyJ},
	sdl.K_KP_MULTIPLY:  []gumak.Key{gumak.KeySymbolShift, gumak.KeyB},
	sdl.K_KP_DIVIDE:    []gumak.Key{gumak.KeySymbolShift, gumak.KeyV},
	sdl.K_LEFTBRACKET:  []gumak.Key{gumak.KeySymbolShift, gumak.Key8},
	sdl.K_RIGHTBRACKET: []gumak.Key{gumak.KeySymbolShift, gumak.Key9},
	sdl.K_SEMICOLON:    []gumak.Key{gumak.KeySymbolShift, gumak.KeyO},
}
