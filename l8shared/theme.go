package l8shared

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/lucasb-eyer/go-colorful"
)

type MyTheme struct{}

func (m MyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// if name == theme.ColorNameBackground {
	// 	if variant == theme.VariantLight {
	// 		return color.White
	// 	}
	// 	return color.Black
	// }

	if name == theme.ColorNameButton {
		aColor, _ := colorful.Hex("#F3E7CF")
		return aColor
	}
	return theme.DefaultTheme().Color(name, theme.VariantLight)
}

func (m MyTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m MyTheme) Size(name fyne.ThemeSizeName) float32 {

	if name == theme.SizeNameInputBorder {
		return theme.DefaultTheme().Size(theme.SizeNameInputBorder) * 0.5
	}
	return theme.DefaultTheme().Size(name)
}

func (m MyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}
