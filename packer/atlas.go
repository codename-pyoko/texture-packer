package packer

import "image"

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Dimension struct {
	Width  int `json:"w"`
	Height int `json:"h"`
}

type Frame struct {
	Position
	Dimension
}

type AtlasEntry struct {
	Frame            Frame     `json:"frame"`
	SpriteSourceSize Frame     `json:"spriteSourceSize"`
	Trimmed          bool      `json:"trimmed"`
	Rotated          bool      `json:"rotated"`
	SourceSize       Dimension `json:"sourceSize"`
}

type Atlas struct {
	Frames map[string]AtlasEntry `json:"frames"`
}

func NewAtlas() Atlas {
	return Atlas{
		Frames: make(map[string]AtlasEntry),
	}
}

func (a *Atlas) Add(x, y int, path string, img image.Image) {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	a.Frames[path] = AtlasEntry{
		Frame: Frame{
			Position:  Position{x, y},
			Dimension: Dimension{w, h},
		},
		SpriteSourceSize: Frame{
			Position:  Position{0, 0},
			Dimension: Dimension{w, h},
		},
		Rotated:    false,
		Trimmed:    false,
		SourceSize: Dimension{w, h},
	}
}
