package tamtam_sdl2_renderer

import (
	"fmt"
	"strconv"
	tt "tamtam/tamtam"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

const TEXTURE_SIZE = 1024
const TILE_SIZE = 64

var BACKGROUND_COLOR = [4]uint8{0.4 * 255, 0.4 * 255, 0.4 * 255}
var COLOR_WHEEL = [][4]uint8{{229, 198, 146, 255}, {20, 196, 52, 255}, {227, 121, 151, 255}}

// Absolute screen coordinates (as well as textureCoordinates) take the assumption that going NORTH is y + 1
// going EAST is x + 1. That does not match SDL internal convention. This gets corrected at render time.
type screenCoordinates tt.Vec2Di

// Relative coordinates inside a texture (from lower left corner)
type textureCoordinates tt.Vec2Di

// We cut the screen plane in square textures of size TEXTURE_SIZE*TEXTURE_SIZE
//
//          -----------------
//         |                 |
//         |(0,TEXTURE_SIZE) |
//          ----------------- -----------------
//         |                 |                 |
//         |(0,0)            |(TEXTURE_SIZE,0) |
//          ----------------- -----------------
//
//
// Tiles have an assembly position (tt.Vec2i) given by the model, for instance [-1, 2]
// which are mapped to screen positions (screenCoordinates) [-1*TILE_SIZE, 2*TILE_SIZE] and then rendered to
// the appropriate texture, here the texture with coordinates [0,0].
type SDL2AssemblyRenderer struct {
	sdlRenderer       *sdl.Renderer
	assembly          *tt.TileAssembly
	tilesTextureCache map[screenCoordinates]*sdl.Texture
	gridTextureCache  map[screenCoordinates]*sdl.Texture
}

func NewSDL2AssemblyRenderer(assembly *tt.TileAssembly, sdlRenderer *sdl.Renderer) (assemblyRenderer SDL2AssemblyRenderer) {

	assemblyRenderer.assembly = assembly
	assemblyRenderer.sdlRenderer = sdlRenderer
	assemblyRenderer.tilesTextureCache = make(map[screenCoordinates]*sdl.Texture)
	assemblyRenderer.gridTextureCache = make(map[screenCoordinates]*sdl.Texture)

	fmt.Println("Creating assembly renderer")
	assemblyRenderer.UpdateTextures()

	return assemblyRenderer
}

func assemblyPosToScreenCoordinates(tilePos tt.Vec2Di) screenCoordinates {
	return screenCoordinates{tilePos[0] * TILE_SIZE, tilePos[1] * TILE_SIZE}
}

// Returns the screen coordinates of the left corner of the texture on which the tile belongs
func getTileTextureLeftCornerCoord(tilePos tt.Vec2Di) screenCoordinates {
	modX := 0
	modY := 0

	if tilePos[0] < 0 {
		modX = -1
	}

	if tilePos[1] < 0 {
		modY = -1
	}

	screenCoord := assemblyPosToScreenCoordinates(tilePos)

	return screenCoordinates{(screenCoord[0]/TEXTURE_SIZE + modX) * TEXTURE_SIZE, (screenCoord[1]/TEXTURE_SIZE + modY) * TEXTURE_SIZE}
}

// Rendering local grid elements of a tile to the grid texture
func (assemblyRenderer *SDL2AssemblyRenderer) renderLocalGrid(texture *sdl.Texture, tilePos tt.Vec2Di) {

	screenCoord := assemblyPosToScreenCoordinates(tilePos)
	textureLeftCornerCoord := getTileTextureLeftCornerCoord(tilePos)

	coordInTexture := textureCoordinates{screenCoord[0] - textureLeftCornerCoord[0], screenCoord[1] - textureLeftCornerCoord[1]}

	successiveSquareVertices := [4][2]int32{{
		int32(coordInTexture[0]), int32(coordInTexture[1] + TILE_SIZE)}, {int32(coordInTexture[0] + TILE_SIZE), int32(coordInTexture[1] + TILE_SIZE)}, {int32(coordInTexture[0] + TILE_SIZE), int32(coordInTexture[1])}, {int32(coordInTexture[0]), int32(coordInTexture[1])}}

	assemblyRenderer.sdlRenderer.SetRenderTarget(texture)
	gfx.LineRGBA(assemblyRenderer.sdlRenderer, successiveSquareVertices[0][0]+1, successiveSquareVertices[0][1]-1, successiveSquareVertices[2][0]-1, successiveSquareVertices[2][1]+1, 0, 0, 0, 255)

	gfx.LineRGBA(assemblyRenderer.sdlRenderer, successiveSquareVertices[0][0]+1, successiveSquareVertices[0][1]-1, successiveSquareVertices[1][0]-1, successiveSquareVertices[1][1]-1, 0, 0, 0, 255)

	gfx.LineRGBA(assemblyRenderer.sdlRenderer, successiveSquareVertices[1][0]-1, successiveSquareVertices[1][1]-1, successiveSquareVertices[2][0]-1, successiveSquareVertices[2][1]+1, 0, 0, 0, 255)

	gfx.LineRGBA(assemblyRenderer.sdlRenderer, successiveSquareVertices[2][0]-1, successiveSquareVertices[2][1]+1, successiveSquareVertices[3][0]+1, successiveSquareVertices[3][1]+1, 0, 0, 0, 255)

	gfx.LineRGBA(assemblyRenderer.sdlRenderer, successiveSquareVertices[3][0]+1, successiveSquareVertices[3][1]+1, successiveSquareVertices[0][0]+1, successiveSquareVertices[0][1]-1, 0, 0, 0, 255)

	gfx.LineRGBA(assemblyRenderer.sdlRenderer, successiveSquareVertices[1][0]-1, successiveSquareVertices[1][1]-1, successiveSquareVertices[3][0]+1, successiveSquareVertices[3][1]+1, 0, 0, 0, 255)
	assemblyRenderer.sdlRenderer.SetRenderTarget(nil)
}

// Rendering the tile to the tile texture
func (assemblyRenderer *SDL2AssemblyRenderer) renderTile(texture *sdl.Texture, tile tt.SquareGlues, tilePos tt.Vec2Di) {
	assemblyRenderer.sdlRenderer.SetRenderTarget(texture)

	screenCoord := assemblyPosToScreenCoordinates(tilePos)
	textureLeftCornerCoord := getTileTextureLeftCornerCoord(tilePos)

	coordInTexture := textureCoordinates{screenCoord[0] - textureLeftCornerCoord[0], screenCoord[1] - textureLeftCornerCoord[1]}

	successiveSquareVertices := [4][2]int32{{
		int32(coordInTexture[0]), int32(coordInTexture[1] + TILE_SIZE)}, {int32(coordInTexture[0] + TILE_SIZE), int32(coordInTexture[1] + TILE_SIZE)}, {int32(coordInTexture[0] + TILE_SIZE), int32(coordInTexture[1])}, {int32(coordInTexture[0]), int32(coordInTexture[1])}}

	for i := 0; i < 4; i += 1 {

		glue := tile[i]
		glue_int, err := strconv.Atoi(glue)

		if glue != tt.NULL_GLUE {
			if err != nil {
				panic(err)
			}
		}

		color := BACKGROUND_COLOR
		if glue != tt.NULL_GLUE && glue_int < len(COLOR_WHEEL) {
			color = COLOR_WHEEL[glue_int]
		}

		assemblyRenderer.sdlRenderer.SetDrawColor(0, 255, 0, 255)

		if glue != tt.NULL_GLUE {
			gfx.FilledTrigonRGBA(assemblyRenderer.sdlRenderer, successiveSquareVertices[i][0], successiveSquareVertices[i][1], successiveSquareVertices[(i+1)%4][0], successiveSquareVertices[(i+1)%4][1], int32(coordInTexture[0]+TILE_SIZE/2), int32(coordInTexture[1]+TILE_SIZE/2), color[0], color[1], color[2], color[3])
		}

	}

	assemblyRenderer.sdlRenderer.SetRenderTarget(nil)
}

func (assemblyRenderer *SDL2AssemblyRenderer) UpdateTextures() {
	for _, tileAndPos := range assemblyRenderer.assembly.GetNewlyAddedTiles() {
		textureLeftCornerCoord := getTileTextureLeftCornerCoord(tileAndPos.Pos)
		// If the texture does not exists we create it
		if _, ok := assemblyRenderer.tilesTextureCache[textureLeftCornerCoord]; !ok {
			var err error
			assemblyRenderer.tilesTextureCache[textureLeftCornerCoord], err = assemblyRenderer.sdlRenderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, TEXTURE_SIZE, TEXTURE_SIZE)

			fmt.Println("Creating tiles texture with bottom left corner:", textureLeftCornerCoord)

			if err != nil {
				panic(err)
			}

			assemblyRenderer.gridTextureCache[textureLeftCornerCoord], err = assemblyRenderer.sdlRenderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, TEXTURE_SIZE, TEXTURE_SIZE)

			fmt.Println("Creating tiles texture with bottom left corner:", textureLeftCornerCoord)

			if err != nil {
				panic(err)
			}

			// Tiles texture
			assemblyRenderer.sdlRenderer.SetRenderTarget(assemblyRenderer.tilesTextureCache[textureLeftCornerCoord])
			assemblyRenderer.sdlRenderer.SetDrawColor(100, 0, 0, BACKGROUND_COLOR[3])
			assemblyRenderer.sdlRenderer.FillRect(&sdl.Rect{0, 0, TEXTURE_SIZE, TEXTURE_SIZE})

			assemblyRenderer.sdlRenderer.SetRenderTarget(nil)

			// Grid texture
			assemblyRenderer.gridTextureCache[textureLeftCornerCoord].SetBlendMode(sdl.BLENDMODE_BLEND)
			assemblyRenderer.sdlRenderer.SetRenderTarget(assemblyRenderer.gridTextureCache[textureLeftCornerCoord])
			assemblyRenderer.sdlRenderer.SetDrawColor(0, 0, 0, 0)
			assemblyRenderer.sdlRenderer.FillRect(&sdl.Rect{0, 0, TEXTURE_SIZE, TEXTURE_SIZE})
			assemblyRenderer.sdlRenderer.SetRenderTarget(nil)
		}

		assemblyRenderer.renderTile(assemblyRenderer.tilesTextureCache[textureLeftCornerCoord], tileAndPos.Tile, tileAndPos.Pos)

		// For debug
		assemblyRenderer.sdlRenderer.SetRenderTarget(assemblyRenderer.tilesTextureCache[textureLeftCornerCoord])
		assemblyRenderer.sdlRenderer.SetDrawColor(0, 100, 100, 255)
		assemblyRenderer.sdlRenderer.SetScale(2, 2)
		assemblyRenderer.sdlRenderer.DrawRect(&sdl.Rect{0, 0, TILE_SIZE / 2, TILE_SIZE / 2})
		assemblyRenderer.sdlRenderer.SetRenderTarget(nil)
		assemblyRenderer.sdlRenderer.SetScale(1, 1)

		assemblyRenderer.renderLocalGrid(assemblyRenderer.gridTextureCache[textureLeftCornerCoord], tileAndPos.Pos)
	}

	assemblyRenderer.assembly.FlushNewlyAddedTiles()
}

// Rendering the scene and correcting here the difference in convention
// between our screen coordinates and SDL's.
func (assemblyRenderer *SDL2AssemblyRenderer) Render(uiParams UIParameters) {

	// Render tiles
	for textureLeftCornerCoord, texture := range assemblyRenderer.tilesTextureCache {

		if !uiParams.IsInCameraView(tt.Vec2Di(textureLeftCornerCoord)) {
			continue
		}

		assemblyRenderer.sdlRenderer.CopyExF(texture, nil, &sdl.FRect{float32(textureLeftCornerCoord[0]-uiParams.Camera.Translation[0]) * uiParams.Camera.ZoomFactor, float32(-1*textureLeftCornerCoord[1]+uiParams.Camera.Translation[1]) * uiParams.Camera.ZoomFactor, float32(TEXTURE_SIZE * uiParams.Camera.ZoomFactor), float32(TEXTURE_SIZE * uiParams.Camera.ZoomFactor)}, 0, nil, sdl.FLIP_VERTICAL)
	}

	if !uiParams.ShowGrid {
		return
	}

	// Render grid
	for textureLeftCornerCoord, texture := range assemblyRenderer.gridTextureCache {

		if !uiParams.IsInCameraView(tt.Vec2Di(textureLeftCornerCoord)) {
			continue
		}

		assemblyRenderer.sdlRenderer.CopyExF(texture, nil, &sdl.FRect{float32(textureLeftCornerCoord[0]-uiParams.Camera.Translation[0]) * uiParams.Camera.ZoomFactor, float32(-1*textureLeftCornerCoord[1]+uiParams.Camera.Translation[1]) * uiParams.Camera.ZoomFactor, float32(TEXTURE_SIZE * uiParams.Camera.ZoomFactor), float32(TEXTURE_SIZE * uiParams.Camera.ZoomFactor)}, 0, nil, sdl.FLIP_VERTICAL)
	}
}

func (assemblyRenderer SDL2AssemblyRenderer) CountTextures() int {
	return len(assemblyRenderer.tilesTextureCache)
}

func (assemblyRenderer *SDL2AssemblyRenderer) Destroy() {
	for _, texture := range assemblyRenderer.tilesTextureCache {
		texture.Destroy()
	}
	for _, texture := range assemblyRenderer.gridTextureCache {
		texture.Destroy()
	}
}
