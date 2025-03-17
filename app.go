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
	theme    *material.Theme
	ops      op.Ops
	settings Settings
	library  Library

	// Widgets
	listWidget widget.List
}

type AlbumIco struct {
	button widget.Clickable
	icon   *widget.Icon
	info   string
}

var listEntries []AlbumIco

func InitApp() (App, error) {
	a := App{
		theme: material.NewTheme(),
		ops:   op.Ops{},
	}

	// init settings
	a.settings = Settings{
		mediaDirectories: []string{"music"},
	}

	// Read media directories
	var err error
	a.library, err = GetLibraryFromMediaDirectories(a.settings.mediaDirectories)
	if err != nil {
		return a, err
	}

	listEntries = make([]AlbumIco, 0, 128)

	for albumName, albumData := range a.library.albums {
		ico, err := widget.NewIcon(albumData.coverArt)
		if err != nil {
			return App{}, err
		}

		listEntries = append(listEntries, AlbumIco{
			icon: ico,
			info: fmt.Sprintf("%s - %s - %d [%s]", albumName, albumData.artist, albumData.year, albumData.duration.String()),
		})
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
				return SplitWidget{}.Layout(gtx, 64,
					func(gtx C) D {
						ico := material.IconButton(a.theme, &listEntries[index].button, listEntries[index].icon, "Album Cover")
						return ico.Layout(gtx)
					},
					func(gtx C) D {
						lbl := material.Label(a.theme, unit.Sp(18), listEntries[index].info)
						return lbl.Layout(gtx)
					},
				)
			})
		}))
}
