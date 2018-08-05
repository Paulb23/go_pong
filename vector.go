package main

type Vector2 struct {
	x, y int32
}

type Vector4 struct {
	x, y, w, h int32
}

func (v Vector4) collides(other Vector4) bool {
	return (v.x < other.x+other.w &&
		v.x+v.w > other.x &&
		v.y < other.y+other.h &&
		v.y+v.h > other.y)
}

func (v Vector4) get_mid_point() Vector2 {
	return Vector2{(v.x / 2) - (v.w / 2), (v.y / 2) - (v.h / 2)}
}
