package ui

import (
	"fmt"
	"hospital-manager/data"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func OpenCalendarWindow(myApp fyne.App) {
	calendarWindow := myApp.NewWindow("🗓️ 看診日程月曆看板")
	calendarWindow.Resize(fyne.NewSize(550, 480)) // 調整成最舒適的黃金比例視窗

	currentDisplayTime := time.Now()

	mainContainer := container.NewVBox()

	var prevBtn *widget.Button
	var nextBtn *widget.Button
	var monthTitle *widget.Label

	prevBtn = widget.NewButton("⬅️ 上個月", func() {
		currentDisplayTime = currentDisplayTime.AddDate(0, -1, 0)
	})

	nextBtn = widget.NewButton("下個月 ➡️", func() {
		currentDisplayTime = currentDisplayTime.AddDate(0, 1, 0)
	})

	monthTitle = widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	headerControl := container.NewBorder(nil, nil, prevBtn, nextBtn, monthTitle)

	var drawCalendar func()
	drawCalendar = func() {
		mainContainer.Objects = nil

		appointments := data.LoadData()
		year, month, _ := currentDisplayTime.Date()

		monthTitle.SetText(fmt.Sprintf("%d 年 %d 月", year, month))

		mainContainer.Add(headerControl)
		mainContainer.Add(widget.NewSeparator())

		grid := container.NewGridWithColumns(7)

		weekTitles := []string{"日", "一", "二", "三", "四", "五", "六"}
		for _, title := range weekTitles {
			grid.Add(widget.NewLabelWithStyle(title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
		}

		firstDayOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, currentDisplayTime.Location())
		startWeekday := int(firstDayOfMonth.Weekday())
		nextMonth := firstDayOfMonth.AddDate(0, 1, 0)
		totalDays := int(nextMonth.Sub(firstDayOfMonth).Hours() / 24)

		hasEvent := make(map[int]bool)
		for _, appt := range appointments {
			t, err := time.Parse("2006-01-02", appt.Date)
			if err == nil && t.Year() == year && t.Month() == month {
				hasEvent[t.Day()] = true
			}
		}

		realToday := time.Now()

		for i := 0; i < startWeekday; i++ {
			grid.Add(widget.NewLabel(""))
		}

		for day := 1; day <= totalDays; day++ {
			dayStr := fmt.Sprintf("%d", day)
			textObj := canvas.NewText(dayStr, theme.Color(theme.ColorNameForeground))
			textObj.Alignment = fyne.TextAlignCenter
			textObj.TextSize = 22

			if hasEvent[day] {
				textObj.Color = imageColorBlue
				textObj.TextStyle = fyne.TextStyle{Bold: true}
				textObj.Text = dayStr + "·"
			}

			if year == realToday.Year() && month == realToday.Month() && day == realToday.Day() {
				bgRect := canvas.NewRectangle(&colorRGB{R: 185, G: 240, B: 200, A: 255})

				bgRect.SetMinSize(fyne.NewSize(30, 30))

				cellStack := container.NewStack(bgRect, container.NewCenter(textObj))
				grid.Add(cellStack)
			} else {
				grid.Add(container.NewCenter(textObj))
			}
		}

		mainContainer.Add(grid)
		mainContainer.Refresh()
	}

	prevBtn.OnTapped = func() {
		currentDisplayTime = currentDisplayTime.AddDate(0, -1, 0)
		drawCalendar()
	}
	nextBtn.OnTapped = func() {
		currentDisplayTime = currentDisplayTime.AddDate(0, 1, 0)
		drawCalendar()
	}

	drawCalendar()

	topTips := widget.NewLabelWithStyle("（ 💡 亮藍色粗體帶點的日期 代表當天有看診或領藥行程 ）", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})

	calendarWindow.SetContent(container.NewBorder(topTips, nil, nil, nil, mainContainer))
	calendarWindow.Show()
}
