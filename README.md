# Gumak-Z80-
ZX Spectrum emulator written in GO. 

This is my take on ZX Spectrum emulator for the occasion of ZX Spectrum 40th aniversary. I always wanted to write an emulator and I also wanted to learn new programming language so here we are :-)

I also wrote an [essay](https://www.linkedin.com/pulse/40th-anniversary-zx-spectrum-tom%2525C3%2525A1%2525C5%2525A1-kot%2525C3%2525A1l%3FtrackingId=CohAdhCnSdGVl0V1uOYX%252Bg%253D%253D) with some details about ZX Spectrum internals and emulator development.

There are 3 parts of the project:
 - **gumak** - implementation of an emulator as a GO package
 - **gumak_sdl** - SDL front-end for emulator that can be run and you can start programming in Basic or play your favourite game
 - **gumak_cbind** - simple wrapper of gumak GO module to C so you can use it in a C/C++ project (or embed it in Unreal Engine if you wish)

[![ZX Spectrum in Unreal Engine](https://img.youtube.com/vi/RsxvStoXF08/0.jpg)](https://www.youtube.com/watch?v=RsxvStoXF08)

The emulator is capable of running 48K ROM of the original ZX Spectrum as well as the newer 128K ROM of the ZX Spectrum 128K+ version. There is also a simple implementation of the original AY-3-8912 sound chip found in the newer versions of ZX Spectrum.

[![AY-3-8912 sound demo](https://img.youtube.com/vi/flPLISOoE8s/0.jpg)](https://www.youtube.com/watch?v=flPLISOoE8s)
