package main

import (
	"errors"
	"github.com/veandco/go-sdl2/sdl"
)

func tick_update(delta float32, uptime float32) {

}

func event_update(event sdl.Event, uptime float32) {

}

func render_update(window *Window) {
	window.SetDrawColor(255, 0, 0, 255)
	window.FillRect(&sdl.Rect{0, 0, 100, 100})
}

func game_loop(window *Window) {
	var (
		running           = true
		timer             = sdl.GetTicks()
		last_time         = sdl.GetTicks()
		delta     float32 = 0
		fps       float32 = 0
		tick      float32 = 0
		uptime    float32 = 0
	)

	const (
		max_fps float32 = 60
		ns      float32 = 1000.0 / 60.0
	)

	for running {
		if fps < max_fps {
			now := sdl.GetTicks()
			delta += float32(now-last_time) / ns
			last_time = now
			window.Present()
			window.Clear()

			for delta >= 1 {
				tick_update(delta, uptime)

				for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
					event_update(event, uptime)

					switch event.(type) {
					case *sdl.QuitEvent:
						running = false
						break
					}
				}
				tick++
				delta--
			}
			fps++

			render_update(window)
			window.SetDrawColor(0, 0, 0, 255)
		}

		if sdl.GetTicks()-timer > 1000 {
			timer += 1000
			uptime++
			fps = 0
			tick = 0
		}
	}
}

func init_sdl() error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_ERROR, "SDL 2", "FATAL: Could not start SDL 2!", nil)
		return errors.New("sdl could not be initilised.")
	}
	return nil
}

func main() {
	if err := init_sdl(); err != nil {
		return
	}

	window, err := create_window("Go Pong!", 800, 600)
	if err != nil {
		sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_ERROR, "Window Creation", "FATAL: "+err.Error()+"!", nil)
		return
	}

	game_loop(window)

	window.destroy()
	sdl.Quit()
	return
}
