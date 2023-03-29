package main

import (
	"fyne.io/fyne/v2"
)

type fillSpace struct{}

func (d *fillSpace) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(600, 400)
}

func (d *fillSpace) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, 0)
	// for _, o := range objects {
	newHeight := containerSize.Height - 10
	newSize := fyne.NewSize(containerSize.Width, newHeight)
	objects[0].Resize(newSize)
	objects[0].Move(pos)
	pos = pos.Add(fyne.NewPos(0, newHeight+10))
	// }
}
