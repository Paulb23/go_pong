package main

type Entity interface {
	draw(window *Window)
	tick_update(delta float32, uptime float32)
	get_transform() Vector4
}
