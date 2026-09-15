package data

import (
	"encoding/json"
	"os"
)

type Appointment struct {
	ID       string `json:"id"`
	Date     string `json:"date"`
	Hospital string `json:"hospital"`
	Clinic   string `json:"clinic"`
	Note     string `json:"note"`
}

const filename = "appointments.json"
const hospitalFilename = "hospitals.json"
const historyFilename = "history.json"

func LoadData() []Appointment {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return GetDummyData()
	}

	fileData, err := os.ReadFile(filename)
	if err != nil {
		return []Appointment{}
	}

	var appointments []Appointment
	err = json.Unmarshal(fileData, &appointments)
	if err != nil {
		return []Appointment{}
	}

	return appointments
}

func SaveData(appointments []Appointment) error {
	fileData, err := json.MarshalIndent(appointments, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, fileData, 0644)
}

func LoadHistoryData() []Appointment {
	if _, err := os.Stat(historyFilename); os.IsNotExist(err) {
		return []Appointment{}
	}
	fileData, err := os.ReadFile(historyFilename)
	if err != nil {
		return []Appointment{}
	}
	var history []Appointment
	_ = json.Unmarshal(fileData, &history)
	return history
}

func SaveToHistory(item Appointment) {
	history := LoadHistoryData()

	history = append([]Appointment{item}, history...)

	fileData, _ := json.MarshalIndent(history, "", " ")
	_ = os.WriteFile(historyFilename, fileData, 0644)
}

func LoadHospitalOptions() []string {
	if _, err := os.Stat(hospitalFilename); os.IsNotExist(err) {
		defaultHospitals := []string{"台大醫院", "榮民總醫院", "長庚醫院", "奇美醫院", "成大醫院", "社區健保藥局"}
		fileData, _ := json.MarshalIndent(defaultHospitals, "", " ")
		_ = os.WriteFile(hospitalFilename, fileData, 0644)

		return append(defaultHospitals, "其他 (自行輸入)")
	}

	fileData, err := os.ReadFile(hospitalFilename)
	if err != nil {
		return []string{"其他 (自行輸入)"}
	}

	var hospitals []string
	err = json.Unmarshal(fileData, &hospitals)
	if err != nil {
		return []string{"其他 (自行輸入)"}
	}

	return append(hospitals, "其他 (自行輸入)")
}

func SaveHospitalOptions(hostitals []string) error {
	var cleanList []string
	for _, h := range hostitals {
		if h != "其他 (自行輸入)" && h != "" {
			cleanList = append(cleanList, h)
		}
	}

	fileData, err := json.MarshalIndent(cleanList, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(hospitalFilename, fileData, 0644)
}

func GetDummyData() []Appointment {
	return []Appointment{
		{
			ID:       "1",
			Date:     "2026-10-15",
			Hospital: "成大醫院",
			Clinic:   "心臟內科",
			Note:     "定期回診，需抽血",
		},
		{
			ID:       "2",
			Date:     "2026-11-05",
			Hospital: "奇美醫院",
			Clinic:   "泌尿科",
			Note:     "拿第2次慢性病處方箋",
		},
		{
			ID:       "3",
			Date:     "2026-11-20",
			Hospital: "社區健保藥局",
			Clinic:   "領藥",
			Note:     "處方箋第3次領藥截止日！",
		},
	}
}
