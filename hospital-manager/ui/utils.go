package ui

import (
	"time"
)

// 顏色與結構定義
var imageColorRed = &colorRGB{R: 230, G: 50, B: 50, A: 255}
var imageColorBlue = &colorRGB{R: 50, G: 120, B: 240, A: 255}  // 舒適的亮藍色
var imageColorGray = &colorRGB{R: 140, G: 140, B: 140, A: 255} // 不搶眼的暗灰色

type colorRGB struct{ R, G, B, A uint8 }

func (c *colorRGB) RGBA() (r, g, b, a uint32) {
	return uint32(c.R) * 0x101, uint32(c.G) * 0x101, uint32(c.B) * 0x101, uint32(c.A) * 0x101
}

// 取得星期縮寫公用函式
func getWeekdayStr(dateStr string) string {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return ""
	}
	weeks := map[time.Weekday]string{
		time.Sunday: "日", time.Monday: "一", time.Tuesday: "二",
		time.Wednesday: "三", time.Thursday: "四", time.Friday: "五", time.Saturday: "六",
	}
	return "(" + weeks[t.Weekday()] + ")"
}
