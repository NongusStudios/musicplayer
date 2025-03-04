package main

import (
	"log"
	"os"

	"gioui.org/app"
)

func main() {
	go func() {
		win := new(app.Window)
		err := run(win)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(win *app.Window) error {
	a, err := initApp()
	if err != nil {
		return err
	}

	for {
		switch e := win.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&a.ops, e)
			a.Update(gtx)
			a.Draw(gtx)
			e.Frame(&a.ops)
		}
	}
}
