package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

type Bat struct {
	starting_transform Vector4
	transform          Vector4
	speed              int32
	player             bool
}

func create_bat(transform Vector4, speed int32, player bool) *Bat {
	bat := Bat{
		transform,
		transform,
		speed,
		player,
	}
	return &bat
}

func (b *Bat) reset() {
	b.transform = b.starting_transform
}

func (b Bat) draw(window *Window) {
	window.FillRect(&sdl.Rect{b.transform.x, b.transform.y, b.transform.w, b.transform.h})
}

func (b *Bat) tick_update(delta float32, uptime float32) {

	if b.player {

		// handle player input here as its smoother
		sdl.PumpEvents()
		keyboard_state := sdl.GetKeyboardState()

		if keyboard_state[sdl.SCANCODE_DOWN] > 0 {
			b.transform.y += int32(float32(b.speed) * delta)
		}
		if keyboard_state[sdl.SCANCODE_UP] > 0 {
			b.transform.y -= int32(float32(b.speed) * delta)
		}

		if b.transform.y < 0 {
			b.transform.y = 0
		}

		if b.transform.y+b.transform.h > WINDOW_HEIGHT {
			b.transform.y = WINDOW_HEIGHT - b.transform.h
		}

		return
	}

	// "ai"
	var closest_ball *Ball
	var closest_dist float64 = 999
	for _, ball := range balls {
		dist := math.Abs(float64(ball.transform.x - b.transform.x))
		if dist < closest_dist {
			closest_ball = ball
			closest_dist = dist
		}
	}

	if closest_ball == nil {
		return
	}

	if closest_dist > WINDOW_WIDTH/2 {
		if b.transform.y > b.starting_transform.y {
			b.transform.y -= b.speed
		} else if b.transform.y < b.starting_transform.y {
			b.transform.y += b.speed
		}
		return
	}

	if closest_ball.transform.y > b.transform.y {
		b.transform.y += b.speed
	} else if closest_ball.transform.y < b.transform.y {
		b.transform.y -= b.speed
	}
}

func (b Bat) get_transform() Vector4 {
	return b.transform
}
