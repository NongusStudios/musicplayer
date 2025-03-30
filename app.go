package main

import (
	"fmt"
	"image"
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

type Settings struct {
	mediaDirectories []string
}

type ListEntryType = int

const (
	ListEntryTypeAlbum ListEntryType = iota
	ListEntryTypeSong
)

type ListEntryData struct {
	typ         ListEntryType
	albumName   string
	trackNumber int

	// Potential Widgets
	clickable *widget.Clickable
}

type App struct {
	theme    *material.Theme
	ops      op.Ops
	settings Settings
	player   Player

	// Widgets
	listWidget      widget.List
	listEntriesData []ListEntryData
	listEntries     []layout.Widget

	playButton widget.Clickable
}

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
	library, err := GetLibraryFromMediaDirectories(a.settings.mediaDirectories)
	if err != nil {
		return a, err
	}

	// init player
	a.player, err = InitPlayer(library)
	if err != nil {
		return a, err
	}

	return a, nil
}

func (a *App) SetupWidgets() error {
	// init widgets
	a.listWidget.List = layout.List{
		Axis: layout.Vertical,
	}

	// Add list entries
	a.listEntries = make([]layout.Widget, 0, 128)
	a.listEntriesData = make([]ListEntryData, 0, 128)

	for albumName, albumData := range a.player.Lib.Albums {
		a.listEntriesData = append(a.listEntriesData, ListEntryData{
			typ: ListEntryTypeAlbum,
		})
		a.listEntries = append(a.listEntries, func(gtx C) D {
			lbl := material.Label(a.theme, unit.Sp(24),
				fmt.Sprintf("%s - %s - %d [%s]", albumName, albumData.artist, albumData.year, albumData.duration.String()))
			return lbl.Layout(gtx)
		})

		for i, song := range a.player.Lib.Songs[albumName] {
			a.listEntriesData = append(a.listEntriesData, ListEntryData{
				typ:         ListEntryTypeSong,
				albumName:   albumName,
				trackNumber: i,
				clickable:   new(widget.Clickable),
			})

			listEntryData := &a.listEntriesData[len(a.listEntriesData)-1]
			a.listEntries = append(a.listEntries, func(gtx C) D {
				return layout.Inset{Left: unit.Dp(25)}.Layout(gtx, func(gtx C) D {
					return SplitWidget{Ratios: []float32{0.1, 0.1, 0.8}}.Layout(gtx, 18,
						func(gtx C) D {
							btn := material.Button(a.theme, listEntryData.clickable, "Play")
							return btn.Layout(gtx)
						},
						layout.Spacer{}.Layout,
						func(gtx C) D {
							lbl := material.Label(a.theme, unit.Sp(16),
								fmt.Sprintf("%s - %s", song.Title, song.Duration.String()))
							return lbl.Layout(gtx)
						},
					)
				})
			})
		}
	}

	return nil
}

func (a *App) Update(gtx C) {
	if a.playButton.Clicked(gtx) {
		a.player.TogglePlayBack()
	}

	for _, listEntryData := range a.listEntriesData {
		if listEntryData.typ == ListEntryTypeSong && listEntryData.clickable.Clicked(gtx) {
			if listEntryData.albumName == a.player.CurrentAlbum && listEntryData.trackNumber == a.player.CurrentTrack {
				a.player.Pause()
				continue
			}
			a.player.PlaySong(listEntryData.albumName, listEntryData.trackNumber)
		}
	}
}

func (a *App) EachFrame() {
	a.player.Update()
}

func (a *App) Draw(gtx C) {
	layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceEnd,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16), Top: unit.Dp(16), Bottom: unit.Dp(16)}.Layout(gtx, func(gtx C) D {
				return SplitWidget{Ratios: []float32{0.1, 0.01, 0.9}}.Layout(gtx, 24,
					func(gtx C) D {
						txt := "Play"
						if a.player.IsPlaying() {
							txt = "Pause"
						}

						btn := material.Button(a.theme, &a.playButton, txt)
						return btn.Layout(gtx)
					},
					layout.Spacer{}.Layout,
					func(gtx C) D {
						slider := material.Slider(a.theme, a.player.CurrentSongProgress)
						return slider.Layout(gtx)
					})
			})
		}),
		layout.Rigid(func(gtx C) D {
			return ColorBox(gtx, image.Pt(gtx.Constraints.Max.X, 1), color.NRGBA{R: 0, G: 0, B: 0, A: 255})
		}),
		layout.Rigid(layout.Spacer{Height: 11}.Layout),
		layout.Rigid(func(gtx C) D {
			list := material.List(a.theme, &a.listWidget)

			return list.Layout(gtx, len(a.listEntries), func(gtx C, index int) D {
				return a.listEntries[index](gtx)
			})
		}),
	)
}
