package main

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

type SplitWidget struct{}

func (n SplitWidget) Layout(gtx C, height int, widgets ...layout.Widget) D {
	size := gtx.Constraints.Min.X / len(widgets)

	for i, wd := range widgets {
		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(size, MinI(height, gtx.Constraints.Max.Y)))

		trans := op.Offset(image.Pt(size*i*ClampI(i, 0, 1), 0)).Push(gtx.Ops)
		wd(gtx)
		trans.Pop()
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
	imageOp.Filter = paint.FilterNearest
	imageOp.Add(ops)
	op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(4, 4))).Add(ops)
	paint.PaintOp{}.Add(ops)
}
