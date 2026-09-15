package ui

import (
	"fmt"
	"hospital-manager/data"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func sortAppointments(appts []data.Appointment) {
	sort.Slice(appts, func(i, j int) bool {
		return appts[i].Date < appts[j].Date
	})
}

func ShowMainWindow(myApp fyne.App) {
	myWindow := myApp.NewWindow("🩺 看診與處方箋管理系統")
	myWindow.Resize(fyne.NewSize(1050, 550))

	appointments := data.LoadData()
	sortAppointments(appointments)

	hospitalOptions := data.LoadHospitalOptions()

	var hospitalSelect *widget.Select

	var list *widget.List
	list = widget.NewList(
		func() int { return len(appointments) },
		func() fyne.CanvasObject {
			textColor := theme.Color(theme.ColorNameBackground)
			textObj := canvas.NewText("", textColor)
			btn := widget.NewButton("❌ 清除", nil)
			textContainer := container.NewHBox(textObj)
			return container.NewBorder(nil, nil, nil, btn, textContainer)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			box := o.(*fyne.Container)

			textContainer := box.Objects[0].(*fyne.Container)
			textObj := textContainer.Objects[0].(*canvas.Text)
			btn := box.Objects[1].(*widget.Button)

			item := appointments[i]

			weekdayStr := getWeekdayStr(item.Date)

			displayText := fmt.Sprintf("[%s%s] %s - %s (%s)", item.Date, weekdayStr, item.Hospital, item.Clinic, item.Note)
			textObj.Text = displayText

			today := time.Now().Local()
			todayZero := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
			apptTime, err := time.ParseInLocation("2006-01-02", item.Date, today.Location())

			if err != nil {
				textObj.Color = theme.Color(theme.ColorNameForeground)
			} else {
				if apptTime.Before(todayZero) {
					textObj.Color = imageColorGray
				} else if apptTime.Sub(todayZero) <= 24*time.Hour*2 {
					textObj.Color = imageColorBlue
				} else if strings.Contains(item.Note, "截止") {
					textObj.Color = imageColorRed
				} else {
					textObj.Color = theme.Color(theme.ColorNameForeground)
				}
			}

			textObj.Refresh()

			btn.OnTapped = func() {
				data.SaveToHistory(appointments[i])

				appointments = append(appointments[:i], appointments[i+1:]...)
				_ = data.SaveData(appointments)

				list.Refresh()
			}
		},
	)

	dateEntry := widget.NewEntry()
	dateEntry.SetText(time.Now().Format("2006-01-02"))

	otherHospitalEntry := widget.NewEntry()
	otherHospitalEntry.SetPlaceHolder("請輸入其他醫院名稱...")
	otherHospitalEntry.Hide()

	hospitalSelect = widget.NewSelect(hospitalOptions, func(selected string) {
		if selected == "其他 (自行輸入)" {
			otherHospitalEntry.Show()
		} else {
			otherHospitalEntry.Hide()
		}
	})
	hospitalSelect.PlaceHolder = "請選擇醫院或藥局"

	clinicEntry := widget.NewEntry()
	clinicEntry.SetPlaceHolder("例如: 心臟內科 / 領藥")

	noteEntry := widget.NewEntry()
	noteEntry.SetPlaceHolder("例如: 記得空腹、拿第2次藥")

	form := widget.NewForm(
		widget.NewFormItem("看診日期", dateEntry),
		widget.NewFormItem("選擇醫院", hospitalSelect),
		widget.NewFormItem("", otherHospitalEntry),
		widget.NewFormItem("看診科別", clinicEntry),
		widget.NewFormItem("小叮嚀", noteEntry),
	)

	addButton := widget.NewButton("➕ 新增看診紀錄", func() {
		if dateEntry.Text == "" || hospitalSelect.Selected == "" {
			return
		}

		finalHospital := hospitalSelect.Selected
		if hospitalSelect.Selected == "其他 (自行輸入)" {
			if otherHospitalEntry.Text == "" {
				return
			}
			finalHospital = otherHospitalEntry.Text
		}

		newAppt := data.Appointment{
			ID:       fmt.Sprintf("%d", len(appointments)+1),
			Date:     dateEntry.Text,
			Hospital: finalHospital,
			Clinic:   clinicEntry.Text,
			Note:     noteEntry.Text,
		}

		appointments = append(appointments, newAppt)
		sortAppointments(appointments)
		_ = data.SaveData(appointments)

		list.Refresh()

		dateEntry.SetText(time.Now().Format("2006-01-02"))
		hospitalSelect.ClearSelected()
		otherHospitalEntry.SetText("")
		otherHospitalEntry.Hide()
		clinicEntry.SetText("")
		noteEntry.SetText("")
	})

	// 頂部選單功能表直接呼叫外部獨立檔案寫好的大寫函式！
	configMenuItem := fyne.NewMenuItem("⚙️ 維護醫院清單", func() { OpenManageWindow(myApp, hospitalSelect) })
	historyMenuItem := fyne.NewMenuItem("📜 查看歷史封存", func() { OpenHistoryWindow(myApp) })

	calendarMenuItem := fyne.NewMenuItem("🗓️ 查看月曆看板", func() { OpenCalendarWindow(myApp) })

	configMenu := fyne.NewMenu("⚙️ 設定", calendarMenuItem, configMenuItem, historyMenuItem)

	mainMenu := fyne.NewMainMenu(configMenu)
	myWindow.SetMainMenu(mainMenu)

	leftSection := container.NewVBox(
		widget.NewLabelWithStyle("📝 新增看診日程", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		form,
		addButton,
	)

	split := container.NewHSplit(leftSection, list)
	split.Offset = 0.38

	now := time.Now()
	dateStr := now.Format("2006-01-02")

	weekdayMap := map[time.Weekday]string{
		time.Sunday:    "星期日",
		time.Monday:    "星期一",
		time.Tuesday:   "星期二",
		time.Wednesday: "星期三",
		time.Thursday:  "星期四",
		time.Friday:    "星期五",
		time.Saturday:  "星期六",
	}
	weekdayStr := weekdayMap[now.Weekday()]

	titleText := fmt.Sprintf("📋 我的健康日曆管理系統 (今天是: %s %s)", dateStr, weekdayStr)
	mainTitle := widget.NewLabelWithStyle(
		titleText,
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	content := container.NewBorder(mainTitle, nil, nil, nil, split)

	myWindow.SetContent(content)
	myWindow.Show()
}
