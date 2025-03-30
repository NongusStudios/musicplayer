package main

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

type SplitWidget struct {
	Ratios []float32 // Must add to 1.0 and be same length as widgets
}

func (n SplitWidget) Layout(gtx C, height int, widgets ...layout.Widget) D {
	offset := 0
	for i, wd := range widgets {
		size := int(float32(gtx.Constraints.Min.X) * n.Ratios[i])

		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(size, MinI(height, gtx.Constraints.Max.Y)))

		trans := op.Offset(image.Pt(offset, 0)).Push(gtx.Ops)
		wd(gtx)
		trans.Pop()

		offset += size
	}

	return layout.Dimensions{Size: image.Pt(gtx.Constraints.Min.X, MinI(height, gtx.Constraints.Max.Y))}
}

func MinI(x int, y int) int {
	if x <= y {
		return x
	}
	return y
}

func ClampI(n, min, max int) int {
	if n < min {
		return min
	} else if n > max {
		return max
	}
	return n
}

func LoadImage(pngFile string) (image.Image, error) {
	file, err := os.Open(pngFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func ColorBox(gtx C, size image.Point, color color.NRGBA) D {
	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}

func FillWithLabel(gtx C, th *material.Theme, text string, backgroundColor color.NRGBA) D {
	ColorBox(gtx, gtx.Constraints.Max, backgroundColor)
	lbl := material.H3(th, text)
	return layout.Center.Layout(gtx, lbl.Layout)
}

func DrawImage(ops *op.Ops, img image.Image) {
	imageOp := paint.NewImageOp(img)
	imageOp.Filter = paint.FilterLinear
	imageOp.Add(ops)
	op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(4, 4))).Add(ops)
	paint.PaintOp{}.Add(ops)
}
