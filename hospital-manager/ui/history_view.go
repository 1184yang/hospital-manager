package ui

import (
	"fmt"
	"hospital-manager/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func OpenHistoryWindow(myApp fyne.App) {
	historyWindow := myApp.NewWindow("📜 歷史看診紀錄封存專區")
	historyWindow.Resize(fyne.NewSize(650, 450))

	historyData := data.LoadHistoryData()

	historyList := widget.NewList(
		func() int { return len(historyData) },
		func() fyne.CanvasObject {
			textObj := canvas.NewText("", imageColorGray)
			textObj.TextSize = 16
			return container.NewHBox(textObj)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			box := o.(*fyne.Container)
			textObj := box.Objects[0].(*canvas.Text)
			item := historyData[i]

			hWeekdayStr := getWeekdayStr(item.Date)
			textObj.Text = fmt.Sprintf("✔ [%s%s] %s - %s (%s)", item.Date, hWeekdayStr, item.Hospital, item.Clinic, item.Note)
			textObj.Refresh()
		},
	)

	historyWindow.SetContent(container.NewBorder(
		widget.NewLabelWithStyle("這裡是您過去已清除的歷史醫療足跡：", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}),
		nil, nil, nil, historyList),
	)
	historyWindow.Show()
}
