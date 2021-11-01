package tamtam_sdl2_renderer

import (
	tt "tamtam/tamtam"
)

type UIParameters struct {
	Translation tt.Vec2Di
	Zoom_factor float32
	Glue_colors map[string][]uint8
}

func NewUIParameters() (toReturn UIParameters) {
	toReturn.Zoom_factor = 1
	return toReturn
}
