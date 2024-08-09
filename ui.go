package main

import (
	"fmt"
	"log"
	"math"

	"github.com/pwiecz/go-fltk"
)

// If the screen is portrait or landscape, the window will be scaled
// accordingly.
const (
	WIDTH_PORTRAIT   = 100
	HEIGHT_PORTRAIT  = 150
	WIDTH_LANDSCAPE  = 150
	HEIGHT_LANDSCAPE = 100

	PAGE_MAIN     uint8 = 0
	PAGE_SETTINGS uint8 = 1
)

// Positioning (x,y,w,h) for fltk elements
type Pos struct {
	X int
	Y int
	W int
	H int
}

// isPortrait returns true if the screen is taller than it is wide. It returns
// false otherwise, including for square screens.
func isPortrait() (bool, error) {
	_, _, width, height := fltk.ScreenWorkArea(int(fltk.SCREEN))

	if width == 0 || height == 0 {
		return false, fmt.Errorf("received 0 for one of screen height or width")
	}

	if width > height {
		return false, nil
	}

	return true, nil
}

// Translates the widget's width/height from the original 100 or 150px base
// width/height to the window's current width/height
func tr(i int, winW int, winH int, useHeight bool) int {
	if portrait {
		if useHeight {
			return int(math.Round((float64(i) / float64(HEIGHT_PORTRAIT)) * float64(winH)))
		} else {
			return int(math.Round((float64(i) / float64(WIDTH_PORTRAIT)) * float64(winW)))
		}
	} else {
		if useHeight {
			return int(math.Round((float64(i) / float64(HEIGHT_LANDSCAPE)) * float64(winH)))
		} else {
			return int(math.Round((float64(i) / float64(WIDTH_LANDSCAPE)) * float64(winW)))
		}
	}
}

// Translate converts a predefined position into a scaled position based on
// the latest width & height of the window.
func (p *Pos) Translate(winW, winH int) {
	p.X = tr(p.X, winW, winH, false)
	p.Y = tr(p.Y, winW, winH, true)
	p.W = tr(p.W, winW, winH, false)
	p.H = tr(p.H, winW, winH, true)
}

func switchPage(p uint8) {
	currentPage = p
	switch p {
	case PAGE_MAIN:
		// hide settings page content
		backBtn.Hide()
		saveBtn.Hide()
		maxEntriesInput.Hide()
		captureIntervalMsInput.Hide()
		darkModeBtn.Hide()
		backBtn.Deactivate()
		saveBtn.Deactivate()
		maxEntriesInput.Deactivate()
		captureIntervalMsInput.Deactivate()
		darkModeBtn.Deactivate()

		// show main page content
		settingsBtn.Activate()
		deleteBtn.Activate()
		copyBtn.Activate()
		logBrowser.Activate()
		settingsBtn.Show()
		deleteBtn.Show()
		copyBtn.Show()
		logBrowser.Show()
	case PAGE_SETTINGS:
		// hide main page content
		settingsBtn.Hide()
		deleteBtn.Hide()
		copyBtn.Hide()
		logBrowser.Hide()
		settingsBtn.Deactivate()
		deleteBtn.Deactivate()
		copyBtn.Deactivate()
		logBrowser.Deactivate()

		// show settings page content
		backBtn.Activate()
		saveBtn.Activate()
		maxEntriesInput.Activate()
		captureIntervalMsInput.Activate()
		darkModeBtn.Activate()
		backBtn.Show()
		saveBtn.Show()
		maxEntriesInput.Show()
		captureIntervalMsInput.Show()
		darkModeBtn.Show()
	}
}

// Resizes and repositions all components based on the window's size.
func responsive(win *fltk.Window) {
	if forceLandscape || forcePortrait {
		return
	}

	winW := win.W()
	winH := win.H()

	if winW > winH {
		portrait = false
	} else {
		portrait = true
	}

	switch currentPage {
	case PAGE_MAIN:
		logBrowserPos := Pos{X: 5, Y: 5, W: 140, H: 75}
		settingsBtnPos := Pos{X: 5, Y: 85, W: 35, H: 10}
		deleteBtnPos := Pos{X: 45, Y: 85, W: 35, H: 10}
		copyBtnPos := Pos{X: 85, Y: 85, W: 60, H: 10}

		if portrait {
			logBrowserPos = Pos{X: 5, Y: 5, W: 90, H: 95}
			settingsBtnPos = Pos{X: 5, Y: 105, W: 90, H: 10}
			deleteBtnPos = Pos{X: 5, Y: 120, W: 90, H: 10}
			copyBtnPos = Pos{X: 5, Y: 135, W: 90, H: 10}
		}

		settingsBtnPos.Translate(winW, winH)
		deleteBtnPos.Translate(winW, winH)
		copyBtnPos.Translate(winW, winH)
		logBrowserPos.Translate(winW, winH)

		settingsBtn.Resize(settingsBtnPos.X, settingsBtnPos.Y, settingsBtnPos.W, settingsBtnPos.H)
		deleteBtn.Resize(deleteBtnPos.X, deleteBtnPos.Y, deleteBtnPos.W, deleteBtnPos.H)
		copyBtn.Resize(copyBtnPos.X, copyBtnPos.Y, copyBtnPos.W, copyBtnPos.H)
		logBrowser.Resize(logBrowserPos.X, logBrowserPos.Y, logBrowserPos.W, logBrowserPos.H)
	// settings page
	case PAGE_SETTINGS:
		back := Pos{X: 5, Y: 85, W: 35, H: 10}
		save := Pos{X: 45, Y: 85, W: 100, H: 10}
		entries := Pos{X: 5, Y: 15, W: 60, H: 10}
		capture := Pos{X: 85, Y: 15, W: 60, H: 10}
		dark := Pos{X: 5, Y: 30, W: 60, H: 10}

		if portrait {
			back = Pos{X: 5, Y: 135, W: 90, H: 10}
			save = Pos{X: 5, Y: 120, W: 90, H: 10}
			entries = Pos{X: 5, Y: 15, W: 90, H: 10}
			capture = Pos{X: 5, Y: 40, W: 90, H: 10}
			dark = Pos{X: 5, Y: 55, W: 90, H: 10}
		}

		back.Translate(winW, winH)
		save.Translate(winW, winH)
		entries.Translate(winW, winH)
		capture.Translate(winW, winH)
		dark.Translate(winW, winH)

		backBtn.Resize(back.X, back.Y, back.W, back.H)
		saveBtn.Resize(save.X, save.Y, save.W, save.H)
		maxEntriesInput.Resize(entries.X, entries.Y, entries.W, entries.H)
		captureIntervalMsInput.Resize(capture.X, capture.Y, capture.W, capture.H)
		darkModeBtn.Resize(dark.X, dark.Y, dark.W, dark.H)
	}
}

const (
	DARK_COLOR_TEXT              fltk.Color = 0x9f9f9f00
	DARK_COLOR_INPUT_BG          fltk.Color = 0x20202000
	DARK_COLOR_INPUT_SELECTED_BG fltk.Color = 0xafafaf00
	// DARK_COLOR_TEXT              fltk.Color = 0x20202000
	// DARK_COLOR_INPUT_BG          fltk.Color = 0x9f9f9f00
	// DARK_COLOR_INPUT_SELECTED_BG fltk.Color = 0xafafaf00

	LIGHT_COLOR_TEXT              fltk.Color = 0x20030500
	LIGHT_COLOR_INPUT_BG          fltk.Color = 0xFFFFFF00
	LIGHT_COLOR_INPUT_SELECTED_BG fltk.Color = 0x00008000
)

var (
	COLOR_TEXT              fltk.Color = LIGHT_COLOR_TEXT
	COLOR_INPUT_BG          fltk.Color = LIGHT_COLOR_INPUT_BG
	COLOR_INPUT_SELECTED_BG fltk.Color = LIGHT_COLOR_INPUT_SELECTED_BG
)

func theme(dark bool) {
	if dark {
		log.Println("dark mode activated")
		COLOR_TEXT = DARK_COLOR_TEXT
		COLOR_INPUT_BG = DARK_COLOR_INPUT_BG
		COLOR_INPUT_SELECTED_BG = DARK_COLOR_INPUT_SELECTED_BG
		fltk.SetForegroundColor(230, 230, 230)
		fltk.SetBackgroundColor(40, 40, 40)
	} else {
		log.Println("light mode activated")
		COLOR_TEXT = LIGHT_COLOR_TEXT
		COLOR_INPUT_BG = LIGHT_COLOR_INPUT_BG
		COLOR_INPUT_SELECTED_BG = LIGHT_COLOR_INPUT_SELECTED_BG
		fltk.SetBackgroundColor(192, 192, 192)
		fltk.SetForegroundColor(0, 0, 0)
		return
	}

	settingsBtn.SetLabelColor(COLOR_TEXT)
	deleteBtn.SetLabelColor(COLOR_TEXT)
	copyBtn.SetLabelColor(COLOR_TEXT)
	logBrowser.SetLabelColor(COLOR_TEXT)
	maxEntriesInput.SetLabelColor(COLOR_TEXT)
	captureIntervalMsInput.SetLabelColor(COLOR_TEXT)
	backBtn.SetLabelColor(COLOR_TEXT)
	saveBtn.SetLabelColor(COLOR_TEXT)
	darkModeBtn.SetLabelColor(COLOR_TEXT)

	settingsBtn.SetColor(COLOR_INPUT_BG)
	deleteBtn.SetColor(COLOR_INPUT_BG)
	copyBtn.SetColor(COLOR_INPUT_BG)
	logBrowser.SetColor(COLOR_INPUT_BG)
	maxEntriesInput.SetColor(COLOR_INPUT_BG)
	captureIntervalMsInput.SetColor(COLOR_INPUT_BG)
	backBtn.SetColor(COLOR_INPUT_BG)
	saveBtn.SetColor(COLOR_INPUT_BG)
	darkModeBtn.SetColor(COLOR_INPUT_BG)

	settingsBtn.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	deleteBtn.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	copyBtn.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	logBrowser.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	maxEntriesInput.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	captureIntervalMsInput.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	backBtn.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	saveBtn.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	darkModeBtn.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
}
