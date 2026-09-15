package main

import (
	"hospital-manager/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

type LargeTextTheme struct {
	fyne.Theme
}

func (m *LargeTextTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return 16
	}

	if name == theme.SizeNameHeadingText {
		return 20
	}

	return theme.DefaultTheme().Size(name)
}

func main() {
	myApp := app.New()

	myApp.Settings().SetTheme(&LargeTextTheme{Theme: theme.DefaultTheme()})

	ui.ShowMainWindow(myApp)

	myApp.Run()
}
