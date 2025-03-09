package main

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type C = layout.Context
type D = layout.Dimensions

const NAV_HEIGHT int = 50

type Settings struct {
	mediaDirectories []string
}

type App struct {
	theme            *material.Theme
	ops              op.Ops
	settings         Settings
	mediaDirectories []MediaDirectory

	// Widgets
	listWidget widget.List
}

var listEntries []string

func InitApp() (App, error) {
	a := App{
		theme: material.NewTheme(),
		ops:   op.Ops{},
	}

	// init settings
	a.settings = Settings{
		mediaDirectories: []string{"music"},
	}

	for _, path := range a.settings.mediaDirectories {
		dir, err := ReadFilesInDirectory(path)
		if err != nil {
			return a, err
		}

		a.mediaDirectories = append(a.mediaDirectories, dir)
	}

	listEntries = make([]string, 0, 128)

	for _, directory := range a.mediaDirectories {
		for path, files := range directory {
			listEntries = append(listEntries, fmt.Sprintf("%s:", path))
			for i, file := range files {
				listEntries = append(listEntries, fmt.Sprintf("    %d - %s", i, file))
			}
		}
	}

	// init widgets
	a.listWidget.List = layout.List{
		Axis: layout.Vertical,
	}

	return a, nil
}

func (a *App) Update(gtx C) {

}

func (a *App) Draw(gtx C) {

	layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceEnd,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			list := material.List(a.theme, &a.listWidget)

			return list.Layout(gtx, len(listEntries), func(gtx layout.Context, index int) layout.Dimensions {
				lbl := material.Label(a.theme, unit.Sp(18), listEntries[index])
				return lbl.Layout(gtx)
			})
		}))
}
