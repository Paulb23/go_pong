package main

import (
	"errors"
	"github.com/veandco/go-sdl2/sdl"
	"image/color"
	"math"
	"math/rand"
)

const (
	WINDOW_TITLE  = "Go Pong!"
	WINDOW_WIDTH  = 512
	WINDOW_HEIGHT = 256

	// state
	SERVE   = 0
	PLAYING = 1
	RESET   = 2
)

var (
	DEFAULT_CLEAR_COLOR = color.RGBA{0, 0, 0, 0}

	NEXT_ROUND_DELAY float32 = 0
	RESET_DELAY      float32 = 2

	NET_WIDTH int32 = 4
	NET_X     int32 = (WINDOW_WIDTH / 2) - (NET_WIDTH / 2)

	BAT_SPEED  int32 = 4
	BAT_HEIGHT int32 = 28
	BAT_WIDTH  int32 = 7
	START_X    int32 = 21
	START_Y    int32 = (WINDOW_HEIGHT / 2) - (BAT_HEIGHT / 2)

	BALL_SIZE    int32 = 10
	BALL_SPEED   int32 = 5
	BALL_START_X int32 = (WINDOW_WIDTH / 2) - (BALL_SIZE / 2)
	BALL_START_Y int32 = (WINDOW_HEIGHT / 2) - (BALL_SIZE / 2)

	player_bat   = Bat{Vector4{START_X, START_Y, BAT_WIDTH, BAT_HEIGHT}, BAT_SPEED}
	opponent_bat = Bat{Vector4{WINDOW_WIDTH - START_X, START_Y, BAT_WIDTH, BAT_HEIGHT}, BAT_SPEED}

	ball = Ball{Vector4{BALL_START_X, BALL_START_Y, 10, 10}, BAT_SPEED, 0, 0}

	game_state int32 = SERVE
)

func tick_update(delta float32, uptime float32) {

	// handle player input here as its smoother
	sdl.PumpEvents()
	keyboard_state := sdl.GetKeyboardState()

	if keyboard_state[sdl.SCANCODE_DOWN] > 0 {
		player_bat.transform.y += int32(float32(player_bat.speed) * delta)
	}
	if keyboard_state[sdl.SCANCODE_UP] > 0 {
		player_bat.transform.y -= int32(float32(player_bat.speed) * delta)
	}

	if player_bat.transform.y < 0 {
		player_bat.transform.y = 0
	}

	if player_bat.transform.y+player_bat.transform.h > WINDOW_HEIGHT {
		player_bat.transform.y = WINDOW_HEIGHT - player_bat.transform.h
	}

	// ball
	switch game_state {
	case SERVE:
		ball.speed_x = -ball.speed
		if rand.Intn(1) == 1 {
			ball.speed_x = ball.speed
		}
		game_state = PLAYING
		break
	case PLAYING:
		ball.transform.x += ball.speed_x
		ball.transform.y += ball.speed_y

		if ball.transform.collides(player_bat.transform) {
			var ball_y = ball.transform.get_mid_point().y
			var player_y = player_bat.transform.get_mid_point().y

			if math.Abs(float64(player_y-ball_y)) > 10 {
				ball.speed_y = ball.speed
			} else if player_y <= ball_y {
				ball.speed_y = -ball.speed
			}
			ball.speed_x = -ball.speed_x
		}

		if ball.transform.collides(opponent_bat.transform) {
			var ball_y = ball.transform.get_mid_point().y
			var opponent_y = opponent_bat.transform.get_mid_point().y

			if math.Abs(float64(opponent_y-ball_y)) > 10 {
				ball.speed_y = ball.speed
			} else if opponent_y <= ball_y {
				ball.speed_y = -ball.speed
			}
			ball.speed_x = -ball.speed_x
		}

		if ball.transform.y < 0 {
			ball.speed_y = -ball.speed_y
		}

		if ball.transform.y+ball.transform.h > WINDOW_HEIGHT {
			ball.speed_y = -ball.speed_y
		}

		if ball.transform.x < 0 || ball.transform.x >= WINDOW_WIDTH {
			NEXT_ROUND_DELAY = uptime + RESET_DELAY
			game_state = RESET
			break
		}

		// "ai"
		if dist := math.Abs(float64(ball.transform.y - opponent_bat.transform.y)); dist > 5 && rand.Int31n(10) > 2 {
			if ball.transform.y > opponent_bat.transform.y {
				opponent_bat.transform.y += opponent_bat.speed
			} else if ball.transform.y < opponent_bat.transform.y {
				opponent_bat.transform.y -= opponent_bat.speed
			}
		}

	case RESET:
		ball.transform.x = BALL_START_X
		ball.transform.y = BALL_START_Y

		player_bat.transform.x = START_X
		player_bat.transform.y = START_Y

		opponent_bat.transform.x = WINDOW_WIDTH - START_X
		opponent_bat.transform.y = START_Y

		if uptime > NEXT_ROUND_DELAY {
			game_state = PLAYING
			break
		}
	}

}

func event_update(event sdl.Event, uptime float32) {

}

func render_update(window *Window) {
	window.SetDrawColor(255, 255, 255, 255)

	// bats
	window.FillRect(&sdl.Rect{player_bat.transform.x, player_bat.transform.y, player_bat.transform.w, player_bat.transform.h})
	window.FillRect(&sdl.Rect{opponent_bat.transform.x, opponent_bat.transform.y, opponent_bat.transform.w, opponent_bat.transform.h})

	// net
	window.FillRect(&sdl.Rect{NET_X, 0, NET_WIDTH, WINDOW_HEIGHT})

	// ball
	window.FillRect(&sdl.Rect{ball.transform.x, ball.transform.y, ball.transform.w, ball.transform.h})
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
	return
}
