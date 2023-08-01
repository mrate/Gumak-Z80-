module mutex/gumak_cwrap

go 1.18

replace mutex/gumak => ../gumak

require mutex/gumak v0.0.0-00010101000000-000000000000

require (
	golang.org/x/sys v0.0.0-20220319134239-a9b59b0215f8
	mutex/gumak_sdl v0.0.0-00010101000000-000000000000
)

require (
	github.com/TheTitanrain/w32 v0.0.0-20180517000239-4f5cfb03fabf // indirect
	github.com/sqweek/dialog v0.0.0-20220227145630-7a1c9e333fcf // indirect
	github.com/veandco/go-sdl2 v0.4.18 // indirect
)

replace mutex/gumak_sdl => ../gumak_sdl
