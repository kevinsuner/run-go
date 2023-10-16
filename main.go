package main

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

const SCREEN_WIDTH int32 = 640
const SCREEN_HEIGHT int32 = 400

func main() {
	var window *sdl.Window = nil
	var surface *sdl.Surface = nil
	var err error

	defer sdl.Quit()
	defer window.Destroy()

	if sdl.INIT_VIDEO < 0 {
		log.Fatalf("SDL could not initialize! SDL_Error: %v\n", sdl.GetError())
	} else {
		window, err = sdl.CreateWindow("SDL Tutorial", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, SCREEN_WIDTH, SCREEN_HEIGHT, sdl.WINDOW_SHOWN)
		if err != nil {
			log.Fatalf("Window could not be created! SDL_Error: %v\n", sdl.GetError())
		}

		surface, err = window.GetSurface()
		if err != nil {
			log.Fatalf("Surface could not be found! SDL_Error: %v\n", sdl.GetError())
		}

		surface.FillRect(nil, sdl.MapRGB(surface.Format, 0xFF, 0xFF, 0xFF))
		window.UpdateSurface()

		var quit bool = false
		for !quit {
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				if event.GetType() == sdl.QUIT {
					quit = true
				}
			}
		}
	}
}
