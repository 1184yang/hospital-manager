package ui

import (
	"hospital-manager/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func OpenManageWindow(myApp fyne.App, hospitalSelect *widget.Select) {
	manageWindow := myApp.NewWindow("⚙️ 常用醫院清單管理")
	manageWindow.Resize(fyne.NewSize(500, 400))

	currentHospitals := data.LoadHospitalOptions()

	var hList *widget.List
	hList = widget.NewList(
		func() int { return len(currentHospitals) },
		func() fyne.CanvasObject {
			lbl := widget.NewLabel("")
			btn := widget.NewButton("🗑️ 刪除", nil)
			return container.NewBorder(nil, nil, nil, btn, lbl)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			box := o.(*fyne.Container)
			lbl := box.Objects[0].(*widget.Label)
			btn := box.Objects[1].(*widget.Button)

			lbl.SetText(currentHospitals[i])

			if currentHospitals[i] == "其他 (自行輸入)" {
				btn.Disable()
			} else {
				btn.Enable()
				btn.OnTapped = func() {
					currentHospitals = append(currentHospitals[:i], currentHospitals[i+1:]...)
					_ = data.SaveHospitalOptions(currentHospitals)
					hList.Refresh()

					hospitalSelect.Options = data.LoadHospitalOptions()
					hospitalSelect.Refresh()
				}
			}
		},
	)
	newHospitalEntry := widget.NewEntry()
	newHospitalEntry.SetPlaceHolder("輸入想增加的醫院/診所名稱")

	addHBtn := widget.NewButton("➕ 點我加入常用清單", func() {
		if newHospitalEntry.Text == "" {
			return
		}

		insertIndex := len(currentHospitals) - 1
		currentHospitals = append(currentHospitals[:insertIndex], append([]string{newHospitalEntry.Text}, currentHospitals[insertIndex:]...)...)

		_ = data.SaveHospitalOptions(currentHospitals)
		hList.Refresh()
		newHospitalEntry.SetText("")

		hospitalSelect.Options = data.LoadHospitalOptions()
		hospitalSelect.Refresh()
	})

	inputArea := container.NewVBox(newHospitalEntry, addHBtn)
	manageWindow.SetContent(container.NewBorder(widget.NewLabel("目前的常用名單："), inputArea, nil, nil, hList))
	manageWindow.Show()
}
