package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
	"math/rand"
)

type Ball struct {
	starting_transform Vector4
	transform          Vector4
	speed              int32
	speed_x            int32
	speed_y            int32
}

func create_ball(transform Vector4, speed int32) *Ball {
	ball := Ball{
		transform,
		transform,
		speed,
		0,
		0,
	}
	return &ball
}

func (b *Ball) reset() {
	b.transform = b.starting_transform
	b.speed_x = 0
	b.speed_y = 0
}

func (b *Ball) serve(delta float32) {
	b.speed_x = -b.speed
	if rand.Intn(1) == 1 {
		b.speed_x = b.speed
	}
}

func (b Ball) draw(window *Window) {
	window.FillRect(&sdl.Rect{b.transform.x, b.transform.y, b.transform.w, b.transform.h})
}

func (b *Ball) tick_update(delta float32, uptime float32) {
	b.transform.x += b.speed_x
	b.transform.y += b.speed_y

	for _, bat := range bats {
		if b.transform.collides(bat.get_transform()) {

			var ball_y = b.transform.get_mid_point().y
			var bat_y = bat.get_transform().get_mid_point().y

			if math.Abs(float64(bat_y-ball_y)) > 10 {
				b.speed_y = b.speed
			} else if bat_y <= ball_y {
				b.speed_y = -b.speed
			}

			b.speed_x = -b.speed_x
		}
	}

	if b.transform.y < 0 {
		b.speed_y = -b.speed_y
	}

	if b.transform.y+b.transform.h > WINDOW_HEIGHT {
		b.speed_y = -b.speed_y
	}

	if b.transform.x < 0 || b.transform.x >= WINDOW_WIDTH {
		NEXT_ROUND_DELAY = uptime + RESET_DELAY
		game_state = RESET
	}
}

func (b Ball) get_transform() Vector4 {
	return b.transform
}
