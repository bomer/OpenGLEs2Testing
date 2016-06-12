// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin linux windows

// An app that draws a green triangle on a red background.
//
// Note: This demo is an early preview of Go 1.5. In order to build this
// program as an Android APK using the gomobile tool.
//
// See http://godoc.org/golang.org/x/mobile/cmd/gomobile to install gomobile.
//
// Get the basic example and use gomobile to build or install it on your device.
//
//   $ go get -d golang.org/x/mobile/example/basic
//   $ gomobile build golang.org/x/mobile/example/basic # will build an APK
//
//   # plug your Android device to your computer or start an Android emulator.
//   # if you have adb installed on your machine, use gomobile install to
//   # build and deploy the APK to an Android target.
//   $ gomobile install golang.org/x/mobile/example/basic
//
// Switch to your device or emulator to start the Basic application from
// the launcher.
// You can also run the application on your desktop by running the command
// below. (Note: It currently doesn't work on Windows.)
//   $ go install golang.org/x/mobile/example/basic && basic
package main

import (
	"encoding/binary"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
	"image"
	"image/png"
	"log"
	"os"

	"math/rand"
)

var (
	images   *glutil.Images
	fps      *debug.FPS
	program  gl.Program
	position gl.Attrib
	offset   gl.Uniform
	color    gl.Uniform
	buf      gl.Buffer

	green  float32
	touchX float32
	touchY float32
	img    glutil.Image
)

func main() {
	app.Main(func(a app.App) {
		var glctx gl.Context
		var sz size.Event
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx, _ = e.DrawContext.(gl.Context)
					onStart(glctx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					onStop(glctx)
					glctx = nil
				}
			case size.Event:
				sz = e
				touchX = float32(sz.WidthPx / 2)
				touchY = float32(sz.HeightPx / 2)
			case paint.Event:
				if glctx == nil || e.External {
					// As we are actively painting as fast as
					// we can (usually 60 FPS), skip any paint
					// events sent by the system.
					continue
				}

				onPaint(glctx, sz)
				a.Publish()
				// Drive the animation by preparing to paint the next frame
				// after this one is shown.
				a.Send(paint.Event{})
			case key.Event:
				if e.Code == key.CodeEscape {
					os.Exit(0)
					break
				}
			case touch.Event:
				touchX = e.X
				touchY = e.Y
			}
		}
	})
}

func onStart(glctx gl.Context) {
	var err error
	program, err = glutil.CreateProgram(glctx, vertexShader, fragmentShader)
	if err != nil {
		log.Printf("error creating GL program: %v", err)
		return
	}

	buf = glctx.CreateBuffer()
	glctx.BindBuffer(gl.ARRAY_BUFFER, buf)
	glctx.BufferData(gl.ARRAY_BUFFER, triangleData, gl.STATIC_DRAW)

	position = glctx.GetAttribLocation(program, "position")
	color = glctx.GetUniformLocation(program, "color")
	offset = glctx.GetUniformLocation(program, "offset")

	images = glutil.NewImages(glctx)
	fps = debug.NewFPS(images)

	// gl.TexImage2D(GL_TEXTURE_2D, 0, 3, SCREEN_WIDTH, SCREEN_HEIGHT, 0, GL_RGB, GL_UNSIGNED_BYTE, (GLvoid*)screenData);
	// test3 := []byte{2, 3, 5}
	// glctx.TexImage2D(gl.TEXTURE_2D, 0, 800, 800, gl.RGB, gl.UNSIGNED_BYTE, test3)
	// var img glutil.Image
	// rec := image.Rect(0, 0, 64, 32)
	img = *images.NewImage(64, 32)
	img.RGBA.Set(0, 0, image.Black)
	img.RGBA.Set(10, 10, image.Black)
	img.RGBA.Set(11, 10, image.Black)
	img.RGBA.Set(12, 10, image.Black)
	img.RGBA.Set(13, 10, image.Black)
	img.RGBA.Set(14, 10, image.Black)
	img.RGBA.Set(15, 10, image.Black)
	img.RGBA.Set(15, 11, image.Black)
	img.RGBA.Set(15, 12, image.Black)
	img.RGBA.Set(15, 13, image.Black)
	img.RGBA.Set(63, 31, image.Black)

	// images.NewImage(w, h)

	w, _ := os.Create("test.png")
	png.Encode(w, img.RGBA)

}

func onStop(glctx gl.Context) {
	glctx.DeleteProgram(program)
	glctx.DeleteBuffer(buf)
	fps.Release()
	images.Release()
}

func onPaint(glctx gl.Context, sz size.Event) {
	glctx.ClearColor(1, 1, 1, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)

	glctx.UseProgram(program)

	green += 0.01
	if green > 1 {
		green = 0
	}

	glctx.BindBuffer(gl.ARRAY_BUFFER, buf)
	// glctx.EnableVertexAttribArray(position)
	// glctx.VertexAttribPointer(position, coordsPerVertex, gl.FLOAT, false, 0, 0)
	// for i := 0; i < 64; i++ {
	// 	for j := 0; j < 32; j++ {
	// 		if i%2 == 0 {
	// 			glctx.Uniform4f(color, 1, 1, 1, 1)
	// 		} else {
	// 			glctx.Uniform4f(color, 0, 0, 0, 1)
	// 		}
	// 		// glctx.Uniform2f(offset, touchX/float32(sz.WidthPx), touchY/float32(sz.HeightPx))
	// 		glctx.Uniform2f(offset, float32(i)/float32(64), float32(j)/float32(32))
	// 		glctx.DrawArrays(gl.TRIANGLE_FAN, 0, vertexCount)

	// 	}
	// }
	// glctx.DisableVertexAttribArray(position)

	// //clear func
	// for i := 0; i < 64; i++ {
	// 	for j := 0; j < 32; j++ {
	// 		img.RGBA.Set(i, j, image.White)

	// 	}

	// }

	for j := 0; j < 5; j++ {

		img.RGBA.Set(rand.Intn(64), rand.Intn(32), image.Black)

	}

	tl := geom.Point{0, 0}
	tr := geom.Point{geom.Pt(sz.WidthPx / 4), 0}
	bl := geom.Point{0, geom.Pt(sz.HeightPx / 4)}
	// ptBottomRight := geom.Point{12 + 32, 16}
	img.Upload()

	// Set up the texture
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	// glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_NEAREST);
	// glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_NEAREST);
	// glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_S, GL_CLAMP);
	// glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_T, GL_CLAMP);

	img.Draw(sz, tl, tr, bl, img.RGBA.Bounds())
	// img.Draw(sz, , sz.WidthPx, 0, img.RGBA.Rect)
	fps.Draw(sz)
}

const squareoffset = 0.057

var triangleData = f32.Bytes(binary.LittleEndian,
	0.0, squareoffset, 0.0, // top left
	0.0, 0.0, 0.0, // bottom left
	squareoffset, 0.0, 0.0, // bottom right
	squareoffset, squareoffset, 0.0,
)

const (
	coordsPerVertex = 3
	vertexCount     = 4
)

const vertexShader = `#version 100
uniform vec2 offset;

attribute vec4 position;
void main() {
	// offset comes in with x/y values between 0 and 1.
	// position bounds are -1 to 1.
	vec4 offset4 = vec4(2.0*offset.x-1.0, 1.0-2.0*offset.y, 0, 0);
	gl_Position = position + offset4;
}`

const fragmentShader = `#version 100
precision mediump float;
uniform vec4 color;
void main() {
	gl_FragColor = color;
}`
