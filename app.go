package main

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type C = layout.Context
type D = layout.Dimensions

type App struct {
	theme *material.Theme
	ops   op.Ops

	// Navbar Widgets
	navLibrary  widget.Clickable
	navSettings widget.Clickable
}

func initApp() App {
	a := App{
		theme: material.NewTheme(),
		ops:   op.Ops{},
	}

	return a
}

func (a *App) draw(e app.FrameEvent) {
	gtx := app.NewContext(&a.ops, e)

	layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceStart,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {

		}),
	)

	e.Frame(gtx.Ops)
}
