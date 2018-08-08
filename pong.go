package main

import (
	"errors"
	"github.com/veandco/go-sdl2/sdl"
	"image/color"
)

const (
	WINDOW_TITLE  = "Go Pong!"
	WINDOW_WIDTH  = 512
	WINDOW_HEIGHT = 256

	// state
	INIT    = 0
	SERVE   = 1
	PLAYING = 2
	RESET   = 3

	GOAL_GAP int32 = 21

	NET_WIDTH int32 = 4

	BAT_SPEED  int32 = 4
	BAT_HEIGHT int32 = 28
	BAT_WIDTH  int32 = 7

	BALL_SIZE  int32 = 10
	BALL_SPEED int32 = 5
)

var (
	DEFAULT_CLEAR_COLOR = color.RGBA{0, 0, 0, 0}

	NEXT_ROUND_DELAY float32 = 0
	RESET_DELAY      float32 = 2

	NET_X int32 = (WINDOW_WIDTH / 2) - (NET_WIDTH / 2)

	balls = []*Ball{}
	bats  = []*Bat{}

	game_state int32 = INIT
)

func tick_update(delta float32, uptime float32) {

	switch game_state {
	case INIT:

		balls = append(balls,
			create_ball(
				Vector4{
					(WINDOW_WIDTH / 2) - (BALL_SIZE / 2),
					(WINDOW_HEIGHT / 2) - (BALL_SIZE / 2),
					BALL_SIZE,
					BALL_SIZE},
				BALL_SPEED))

		bats = append(bats,
			create_bat(
				Vector4{
					GOAL_GAP,
					(WINDOW_HEIGHT / 2) - (BAT_HEIGHT / 2),
					BAT_WIDTH,
					BAT_HEIGHT,
				},
				BAT_SPEED,
				true))

		bats = append(bats,
			create_bat(
				Vector4{
					WINDOW_WIDTH - GOAL_GAP,
					(WINDOW_HEIGHT / 2) - (BAT_HEIGHT / 2),
					BAT_WIDTH,
					BAT_HEIGHT},
				BAT_SPEED,
				false))

		game_state = SERVE
		break
	case SERVE:
		for _, ball := range balls {
			ball.serve(delta)
		}
		game_state = PLAYING
		break
	case PLAYING:
		for _, bat := range bats {
			bat.tick_update(delta, uptime)
		}

		for _, ball := range balls {
			ball.tick_update(delta, uptime)
		}
	case RESET:
		for _, ball := range balls {
			ball.reset()
		}

		for _, bat := range bats {
			bat.reset()
		}

		if uptime > NEXT_ROUND_DELAY {
			game_state = SERVE
			break
		}
	}

}

func event_update(event sdl.Event, uptime float32) {

}

func render_update(window *Window) {
	window.SetDrawColor(255, 255, 255, 255)

	// net
	window.FillRect(&sdl.Rect{NET_X, 0, NET_WIDTH, WINDOW_HEIGHT})

	for _, bat := range bats {
		bat.draw(window)
	}

	for _, ball := range balls {
		ball.draw(window)
	}
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
		ns float32 = 1000.0 / 60.0
	)

	for running {
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
		r, g, b, a := DEFAULT_CLEAR_COLOR.RGBA()
		window.SetDrawColor(uint8(r), uint8(g), uint8(b), uint8(a))

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

	window, err := create_window(WINDOW_TITLE, WINDOW_WIDTH, WINDOW_HEIGHT)
	if err != nil {
		sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_ERROR, "Window Creation", "FATAL: "+err.Error()+"!", nil)
		return
	}

	game_loop(window)

	window.destroy()
	sdl.Quit()
}
