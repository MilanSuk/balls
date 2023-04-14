/*
Copyright 2023 Milan Suk

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"log"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

func InitSDLGlobal() error {
	return sdl.Init(sdl.INIT_EVERYTHING)
}
func DestroySDLGlobal() {
	sdl.Quit()
}

type Ui struct {
	window *sdl.Window
	render *sdl.Renderer

	world *World
}

func NewUi() (*Ui, error) {
	var ui Ui
	var err error

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "2")

	ui.window, err = sdl.CreateWindow("Balls", 50, 50, 1280, 720, sdl.WINDOW_RESIZABLE)
	if err != nil {
		return nil, fmt.Errorf("CreateWindow() failed: %w", err)
	}
	ui.render, err = sdl.CreateRenderer(ui.window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return nil, fmt.Errorf("CreateRenderer() failed: %w", err)
	}
	sdl.EventState(sdl.DROPFILE, sdl.ENABLE)
	sdl.StartTextInput()

	ui.world = WorldTest()

	return &ui, nil
}

func (ui *Ui) Destroy() error {
	var err error

	err = ui.render.Destroy()
	if err != nil {
		return fmt.Errorf("Render.Destroy() failed: %w", err)
	}
	err = ui.window.Destroy()
	if err != nil {
		return fmt.Errorf("Window.Destroy() failed: %w", err)
	}

	return nil
}

func (ui *Ui) GetMousePosition() (int32, int32) {

	gx, gy, _ := sdl.GetGlobalMouseState()
	wx, wy := ui.window.GetPosition()

	return (gx - wx), (gy - wy)
}

func (ui *Ui) SaveScreenshot() error {

	w, h, err := ui.render.GetOutputSize()
	if err != nil {
		return fmt.Errorf("GetOutputSize() failed: %w", err)
	}

	surface, err := sdl.CreateRGBSurface(0, w, h, 32, 0, 0, 0, 0)
	if err != nil {
		return fmt.Errorf("CreateRGBSurface() failed: %w", err)
	}

	err = ui.render.ReadPixels(nil, surface.Format.Format, unsafe.Pointer(&surface.Pixels()[0]), int(surface.Pitch))
	if err != nil {
		return fmt.Errorf("ReadPixels() failed: %w", err)
	}

	err = img.SavePNG(surface, "screenshot_"+time.Now().Format("2006-1-2_15-4-5")+".png")
	if err != nil {
		return fmt.Errorf("SavePNG() failed: %w", err)
	}

	surface.Free()
	return nil
}

func (ui *Ui) Event() (bool, error) {

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() { // some cases have RETURN(don't miss it in tick), some (can be missed in tick)!

		switch val := event.(type) {
		case sdl.QuitEvent:
			fmt.Println("Exiting ..")
			return false, nil

		case sdl.MouseButtonEvent:
			if val.Type == sdl.MOUSEBUTTONDOWN {
				sdl.CaptureMouse(true) // keep getting info even mouse is outside window

			} else if val.Type == sdl.MOUSEBUTTONUP {
				sdl.CaptureMouse(false)
			}

			return true, nil

		case sdl.MouseWheelEvent:
			return true, nil
		}

		//ui.SaveScreenshot() ...
	}

	return true, nil
}

func (ui *Ui) Draw() error {

	const SUB_STEPS = 10
	for i := 0; i < SUB_STEPS; i++ {
		ui.world.Solve(1.0 / 60 / SUB_STEPS)
	}

	for _, obj := range ui.world.objs {

		p := obj.pos.Mult(100)

		gfx.AAEllipseRGBA(ui.render, int32(p.x), int32(p.y), 3, 3, 10, 10, 10, 255)
	}

	//gfx.AALineRGBA(render, start.x, start.y, end.x, end.y, cd.r, cd.g, cd.b, cd.a)
	//gfx.FilledEllipseRGBA(self.render, p.x, p.y, coord.size.x/2, coord.size.y/2, cd.r, cd.g, cd.b, cd.a)
	//gfx.AAEllipseRGBA(self.render, p.x, p.y, coord.size.x/2, coord.size.y/2, cd.r, cd.g, cd.b, cd.a)

	return nil
}

func (ui *Ui) Tick() bool {

	running := true
	for running {

		// clears frame
		err := ui.render.SetDrawColor(220, 220, 220, 255)
		if err != nil {
			log.Printf("SetDrawColor() failed: %v\n", err)
			return false
		}
		err = ui.render.Clear()
		if err != nil {
			log.Printf("RenderClear() failed: %v\n", err)
			return false
		}

		//draws frame
		err = ui.Draw()
		if err != nil {
			log.Printf("Draw() failed: %v\n", err)
			return false
		}

		// finishes frame
		ui.render.Present()

		running, _ = ui.Event()
	}

	return running
}
