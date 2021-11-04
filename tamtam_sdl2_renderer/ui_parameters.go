package tamtam_sdl2_renderer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	tt "tamtam/tamtam"
)

const CAMERA_PARAMETERS_DOT_FILE = ".tamtam_camera"

type CameraParameters struct {
	Translation tt.Vec2Di `json:"translation"`
	ZoomFactor  float32   `json:"zoom_factor"`
}

type UIParameters struct {
	Camera        CameraParameters   `json:"camera"`
	GlueColors    map[string][]uint8 `json:"glue_colors"`
	ShowGrid      bool               `json:"show_grid"`
	ShowTilesText bool               `json:"show_tiles_text"`
}

func NewUIParameters() (toReturn UIParameters) {

	toReturn.Camera.ZoomFactor = 1
	if _, err := os.Stat(CAMERA_PARAMETERS_DOT_FILE); err == nil {
		fmt.Println("Loading ", CAMERA_PARAMETERS_DOT_FILE)
		file, err := ioutil.ReadFile(CAMERA_PARAMETERS_DOT_FILE)

		if err != nil {
			fmt.Println(err)
			return toReturn
		}

		var camera CameraParameters
		err = json.Unmarshal([]byte(file), &camera)

		if err != nil {
			fmt.Println(err)
			return toReturn
		}

		toReturn.Camera = camera
	}

	return toReturn
}

func (uiParams UIParameters) DumpCamera() {
	b, err := json.MarshalIndent(uiParams.Camera, "", "  ")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Writing camera parameters to", CAMERA_PARAMETERS_DOT_FILE)
	fmt.Println(string(b))

	err = ioutil.WriteFile(CAMERA_PARAMETERS_DOT_FILE, b, 0644)

	if err != nil {
		fmt.Println(err)
		return
	}
}

// Returns true if the screen coordinate is within the camera view
func (uiParams UIParameters) IsInCameraView(pos tt.Vec2Di) bool {

	// TODO
	return true

}
