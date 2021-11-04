package main

import (
	"fmt"
	tt "tamtam/tamtam"
	ttr "tamtam/tamtam_sdl2_renderer"

	"github.com/veandco/go-sdl2/sdl"
)

func newAssembly() tt.TileAssembly {
	SIZE := 20
	tileSet, err := tt.NewCrtTileSet(2, 3)

	if err != nil {
		panic(err)
	}

	var initialAssembly = make(map[tt.Vec2Di]tt.SquareGlues)

	for i := 0; i < SIZE; i += 1 {
		initialAssembly[tt.Vec2Di{-1, i}] = tt.SquareGlues{tt.NULL_GLUE, "0", tt.NULL_GLUE, tt.NULL_GLUE}
	}

	initialAssembly[tt.Vec2Di{0, -1}] = tt.SquareGlues{"1", tt.NULL_GLUE, tt.NULL_GLUE, tt.NULL_GLUE}

	for i := 0; i < SIZE-1; i += 1 {
		initialAssembly[tt.Vec2Di{1 + i, -1}] = tt.SquareGlues{"0", tt.NULL_GLUE, tt.NULL_GLUE, tt.NULL_GLUE}
	}

	var assembly = tt.NewAssembly(tileSet, initialAssembly, 2)

	didGrow, err := assembly.GrowSync(true)

	for didGrow && err == nil {
		didGrow, err = assembly.GrowSync(true)
	}

	if err != nil {
		panic(err)
	}

	return assembly
}

func countNumberKeyPressed() (count int) {
	for _, value := range sdl.GetKeyboardState() {
		if value != 0 {
			count += 1
		}
	}
	return count
}

// Determines how much to move the camera when arrow keys are pressed
// in units of tiles
func translationSpeed(mod uint16) int {

	if mod&sdl.KMOD_LSHIFT != 0 {
		return 5 * (countNumberKeyPressed() - 1)
	}

	return 1
}

func main() {

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("tamtam - v0.0.1", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_TARGETTEXTURE)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	assembly := newAssembly()

	assemblyRender := ttr.NewSDL2AssemblyRenderer(&assembly, renderer)
	defer assemblyRender.Destroy()

	uiParameters := ttr.NewUIParameters()

	running := true

	totalFrameTicks := 0
	totalFrames := 0
	var framePerf uint64 = 0
	var frameTime float32 = 0

	for running {
		totalFrames += 1

		startTicks := sdl.GetTicks()
		startPerf := sdl.GetPerformanceCounter()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break

			case *sdl.KeyboardEvent:
				switch t.Type {
				case sdl.KEYDOWN:
					switch t.Keysym.Sym {

					case sdl.K_LEFT:
						uiParameters.Camera.Translation[0] -= ttr.TILE_SIZE * translationSpeed(t.Keysym.Mod)
						break
					case sdl.K_RIGHT:
						uiParameters.Camera.Translation[0] += ttr.TILE_SIZE * translationSpeed(t.Keysym.Mod)
						break
					case sdl.K_UP:
						uiParameters.Camera.Translation[1] += ttr.TILE_SIZE * translationSpeed(t.Keysym.Mod)
						break
					case sdl.K_DOWN:
						uiParameters.Camera.Translation[1] -= ttr.TILE_SIZE * translationSpeed(t.Keysym.Mod)
						break

					case sdl.K_z:
						uiParameters.Camera.ZoomFactor *= 1.5
						break
					case sdl.K_a:
						uiParameters.Camera.ZoomFactor /= 1.5
						break

					case sdl.K_n:
						assembly.GrowSync(true)
						assemblyRender.UpdateTextures()
						break

					case sdl.K_s:
						fmt.Println(" Summary\n", "========\n", "Number of tiles:", assembly.Size(), "\n", "Number of textures:", assemblyRender.CountTextures())
						break
					case sdl.K_f:
						fmt.Println("Current FPS: ", 1/frameTime)
						fmt.Println("Average FPS: ", 1000/(totalFrameTicks/totalFrames))
						fmt.Println("Current Perf: ", framePerf)

					case sdl.K_g:
						uiParameters.ShowGrid = !uiParameters.ShowGrid
						break

					// Dumping camera parameters
					case sdl.K_d:
						uiParameters.DumpCamera()
						break

					case sdl.K_ESCAPE:
						running = false
						break
					}

					break
				}

			}

		}

		renderer.SetDrawColor(ttr.BACKGROUND_COLOR[0], ttr.BACKGROUND_COLOR[1], ttr.BACKGROUND_COLOR[2], ttr.BACKGROUND_COLOR[3])

		renderer.Clear()
		assemblyRender.Render(uiParameters)

		renderer.Present()

		sdl.Delay(20)
		endTicks := sdl.GetTicks()
		endPerf := sdl.GetPerformanceCounter()
		framePerf = endPerf - startPerf
		frameTime = float32(endTicks-startTicks) / 1000
		totalFrameTicks += int(endTicks) - int(startTicks)

	}
}
