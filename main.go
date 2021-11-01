package main

import (
	"fmt"
	tt "tamtam/tamtam"
	ttr "tamtam/tamtam_sdl2_renderer"

	"github.com/veandco/go-sdl2/sdl"
)

func modifyTexture(renderer *sdl.Renderer, texture *sdl.Texture) {
	renderer.SetRenderTarget(texture)

	renderer.SetDrawColor(255, 255, 255, 255)
	for i := 0; i < 256; i += 1 {
		for j := 0; j < 256; j += 1 {
			renderer.SetDrawColor(uint8(i/4), uint8(i/2), uint8(i/3), 255)
			renderer.DrawPoint(int32(i), int32(j))
		}
	}

	renderer.SetRenderTarget(nil)
}

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

func main() {

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
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
					case sdl.K_z:
						uiParameters.Zoom_factor *= 1.5
						break
					case sdl.K_a:
						uiParameters.Zoom_factor /= 1.5
						break
					case sdl.K_p:
						fmt.Println(frameTime, totalFrameTicks, totalFrames, framePerf)
						fmt.Println("Current FPS: ", 1/frameTime)
						fmt.Println("Average FPS: ", 1000/(totalFrameTicks/totalFrames))
						fmt.Println("Current Perf: ", framePerf)
					case sdl.K_LEFT:
						uiParameters.Translation[0] -= ttr.TILE_SIZE
						break
					case sdl.K_RIGHT:
						uiParameters.Translation[0] += ttr.TILE_SIZE
						break
					case sdl.K_UP:
						uiParameters.Translation[1] += ttr.TILE_SIZE
						break
					case sdl.K_DOWN:
						uiParameters.Translation[1] -= ttr.TILE_SIZE
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
