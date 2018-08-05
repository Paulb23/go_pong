package main

import (
	"errors"
	"github.com/veandco/go-sdl2/sdl"
)

type Window struct {
	*sdl.Window
	*sdl.Renderer
}

func create_window(title string, width int32, height int32) (*Window, error) {

	window := Window{}

	sdl_window, err := sdl.CreateWindow(
		title,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		width,
		height,
		sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, errors.New("error creating window.")
	}
	window.Window = sdl_window

	sdl_renderer, err := sdl.CreateRenderer(
		window.Window,
		-1,
		sdl.RENDERER_ACCELERATED)
	if err != nil {
		window.Window.Destroy()
		return nil, errors.New("error creating renderer.")
	}
	window.Renderer = sdl_renderer

	window.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	window.SetLogicalSize(width, height)

	return &window, nil
}

func (w *Window) destroy() {
	w.Window.Destroy()
	w.Renderer.Destroy()
}
