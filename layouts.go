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
	// pos = pos.Add(fyne.NewPos(0, newHeight+10))
	// }
}

type longEntry struct{}

func (d *longEntry) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(100, 30)
}

func (d *longEntry) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, 0)
	// for _, o := range objects {
	newHeight := containerSize.Height
	newSize := fyne.NewSize(containerSize.Width, newHeight)
	objects[0].Resize(newSize)
	objects[0].Move(pos)
	// pos = pos.Add(fyne.NewPos(0, newHeight+10))
	// }
}

type halfes struct {
}

func (d *halfes) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for i, o := range objects {
		childSize := o.MinSize()

		w += childSize.Width
		if i == 0 {
			h += childSize.Height
		}
	}
	return fyne.NewSize(w, h)
}

func (d *halfes) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, containerSize.Height-d.MinSize(objects).Height)
	for _, o := range objects {
		size := o.MinSize()
		newWidth := (containerSize.Width / float32(len(objects)))
		newSize := fyne.NewSize(newWidth, size.Height)
		o.Resize(newSize)
		// o.Resize(size)
		o.Move(pos)

		pos = pos.Add(fyne.NewPos(newSize.Width, 0))
	}
}
