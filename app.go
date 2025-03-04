package main

import (
	"fmt"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type C = layout.Context
type D = layout.Dimensions

const NAV_HEIGHT int = 50

type Page int

const (
	PageLibrary Page = iota
	PageSettings
)

type Settings struct {
	mediaDirectories []string
}

type App struct {
	theme      *material.Theme
	ops        op.Ops
	settings   Settings
	mediaFiles []MediaFile

	// State
	currentPage Page

	// Widgets
	btnLibrary  widget.Clickable
	btnSettings widget.Clickable

	listLibrary widget.List
}

func initApp() (App, error) {
	a := App{
		theme: material.NewTheme(),
		ops:   op.Ops{},
		settings: Settings{
			mediaDirectories: []string{"music"},
		},
		currentPage: PageLibrary,
	}

	files, err := ReadFilesInDirectory(a.settings.mediaDirectories[0])
	if err != nil {
		return a, err
	}
	a.mediaFiles = files

	a.listLibrary.List = layout.List{
		Axis: layout.Vertical,
	}

	return a, nil
}

func (a *App) handleNavInput(gtx C) {
	if a.btnLibrary.Clicked(gtx) {
		a.currentPage = PageLibrary
	} else if a.btnSettings.Clicked(gtx) {
		a.currentPage = PageSettings
	}
}

func (a *App) Update(gtx C) {
	a.handleNavInput(gtx)
}

func (a *App) Draw(gtx C) {
	pages := make(map[Page]func(gtx C) D)

	// Library Page
	pages[PageLibrary] = func(gtx C) D {
		lst := material.List(a.theme, &a.listLibrary)
		return lst.Layout(gtx, len(a.mediaFiles), func(gtx layout.Context, index int) layout.Dimensions {
			lbl := material.Label(a.theme, unit.Sp(16), fmt.Sprintf("%d - %s", index, a.mediaFiles[index].name))
			return lbl.Layout(gtx)
		})
	}

	// Settings Page
	pages[PageSettings] = func(gtx C) D {
		return FillWithLabel(gtx, a.theme, "Settings", color.NRGBA{R: 0, G: 0, B: 255, A: 255})
	}

	layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceEnd,
	}.Layout(gtx,
		// Navbar
		layout.Rigid(func(gtx C) D {
			return SplitWidget{}.Layout(gtx, NAV_HEIGHT,
				func(gtx C) D {
					btn := material.Button(a.theme, &a.btnLibrary, "Library")
					return btn.Layout(gtx)
				},
				func(gtx C) D {
					btn := material.Button(a.theme, &a.btnSettings, "Settings")
					return btn.Layout(gtx)
				},
			)
		}),
		layout.Rigid(pages[a.currentPage]),
	)
}
