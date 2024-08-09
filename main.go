package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"slices"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/adrg/xdg"
	"github.com/atotto/clipboard"
	"github.com/pwiecz/go-fltk"
)

const (
	DEFAULT_MAX_ENTRIES         = 100
	DEFAULT_CAPTURE_INTERVAL_MS = 1000
)

var (
	forcePortrait  bool
	forceLandscape bool
	portrait       bool

	maxEntries        int
	captureIntervalMs int

	// Data is stored between runs of this application in this yml config file.
	configFilePath string
	appConf        AppConfig

	currentPage uint8 = PAGE_MAIN
)

type ClipboardEntry struct {
	Value    string
	Selected bool
}

type AppConfig struct {
	// Each value is the value stored in the clipboard, and if it is selected
	// it will be true or false.
	// Log map[string]ClipboardEntry
	Log               []ClipboardEntry `json:"log"`
	CaptureIntervalMS int              `json:"captureIntervalMs"`
	MaxEntries        int              `json:"maxEntries"`
	DarkMode          bool             `json:"darkMode"`
	// A list of secrets and the values to mask them with.
	// Can only be supplied by directly editing the config.
	Secrets map[string]string `json:"secrets"`
}

// Buttons, inputs, widgets, etc that need to be repositioned in a
// responsive manner.
var (
	win *fltk.Window
	// For switching to the settings pane.
	settingsBtn *fltk.Button
	// For deleting the currently selected entries.
	deleteBtn *fltk.Button
	// For copying the currently selected item.
	copyBtn *fltk.Button
	// For saving settings - only shown on the settings page.
	// saveBtn *fltk.Button
	// Each clipboard entry will go into here.
	logBrowser *fltk.MultiBrowser

	// Settings page items
	maxEntriesInput        *fltk.Input
	captureIntervalMsInput *fltk.Input
	backBtn                *fltk.Button
	saveBtn                *fltk.Button
	darkModeBtn            *fltk.CheckButton
)

func parseFlags() {
	flag.BoolVar(&forcePortrait, "portrait", false, "force portrait orientation for the interface")
	flag.BoolVar(&forceLandscape, "landscape", false, "force landscape orientation for the interface")
	flag.StringVar(&configFilePath, "f", "", "the config file to write to, instead of the default provided by XDG config directories")
	flag.IntVar(&captureIntervalMs, "ms", DEFAULT_CAPTURE_INTERVAL_MS, "interval between each attempt to read the clipboard")
	flag.IntVar(&maxEntries, "entries", DEFAULT_MAX_ENTRIES, "interval between each attempt to read the clipboard")
	flag.Parse()
}

func main() {
	parseFlags()

	var err error

	if configFilePath == "" {
		configFilePath, err = xdg.SearchConfigFile("go-fltk-clipboard/config.json")
		if err != nil {
			log.Printf("failed to get xdg config dir: %v", err.Error())
		}
	}

	if configFilePath != "" {
		bac, err := os.ReadFile(configFilePath)
		if err != nil {
			log.Printf("config file not readable at %v", configFilePath)
		}

		err = json.Unmarshal(bac, &appConf)
		if err != nil {
			log.Printf("config file %v failed to parse: %v", configFilePath, err.Error())
		}

		log.Printf("loaded config from %v", configFilePath)
	} else {
		if xdg.ConfigHome != "" {
			configFilePath = path.Join(xdg.ConfigHome, "go-fltk-clipboard", "config.json")
			log.Printf("using %v for config file path", configFilePath)
		} else {
			log.Println("unable to automatically identify any suitable config dirs; configuration will not be saved")
		}
	}

	// propagate values to the config if unset previously
	if appConf.CaptureIntervalMS == 0 {
		appConf.CaptureIntervalMS = captureIntervalMs
	}
	if appConf.MaxEntries == 0 {
		appConf.MaxEntries = maxEntries
	}
	if appConf.Secrets == nil {
		appConf.Secrets = make(map[string]string)
	}

	portrait, err = isPortrait()
	if err != nil {
		log.Fatalf("failed to determine screen size: %v", err.Error())
	}

	// probably could write this more intelligently later
	windowWidth := WIDTH_LANDSCAPE
	windowHeight := HEIGHT_LANDSCAPE
	if portrait || forcePortrait {
		windowWidth = WIDTH_PORTRAIT
		windowHeight = HEIGHT_PORTRAIT
		portrait = true
	}
	if forceLandscape {
		windowWidth = WIDTH_LANDSCAPE
		windowHeight = HEIGHT_LANDSCAPE
		portrait = false
	}

	win = fltk.NewWindow(windowWidth, windowHeight, "Clipboard Manager FLTK")
	// fltk.SetScheme("gtk+")
	fltk.InitStyles()
	fltk.SetTooltipDelay(0.1)
	win.Resizable(win)

	// main page widgets
	settingsBtn = fltk.NewButton(0, 0, 0, 0, "&Settings")
	deleteBtn = fltk.NewButton(0, 0, 0, 0, "&Delete")
	copyBtn = fltk.NewButton(0, 0, 0, 0, "&Copy")
	logBrowser = fltk.NewMultiBrowser(0, 0, 0, 0)
	logBrowser.SetLabelSize(10)
	logBrowser.SetLabelFont(fltk.HELVETICA)

	// settings page widgets
	backBtn = fltk.NewButton(0, 0, 0, 0, "&Back")
	saveBtn = fltk.NewButton(0, 0, 0, 0, "&Save")
	maxEntriesInput = fltk.NewInput(0, 0, 0, 0, "&Max Items")
	captureIntervalMsInput = fltk.NewInput(0, 0, 0, 0, "&Capture Interval (ms)")
	darkModeBtn = fltk.NewCheckButton(0, 0, 0, 0, "&Dark Mode")

	maxEntriesInput.SetValue(fmt.Sprint(appConf.MaxEntries))
	captureIntervalMsInput.SetValue(fmt.Sprint(appConf.CaptureIntervalMS))

	maxEntriesInput.SetTooltip(fmt.Sprintf("This can be a large number, but performance may suffer. Default=%v", DEFAULT_MAX_ENTRIES))
	captureIntervalMsInput.SetTooltip(fmt.Sprintf("The smaller the interval, the sooner clipboard events will show up in this app. Avoid setting this number too small. Default=%v", DEFAULT_CAPTURE_INTERVAL_MS))
	darkModeBtn.SetTooltip("Toggling the UI mode requires a restart, and this setting will persist to settings between app restarts.")

	maxEntriesInput.SetAlign(fltk.ALIGN_TOP_LEFT)
	captureIntervalMsInput.SetAlign(fltk.ALIGN_TOP_LEFT)
	logBrowser.SetAlign(fltk.ALIGN_BOTTOM_LEFT)

	// hide the settings page widgets on first load
	backBtn.Hide()
	saveBtn.Hide()
	maxEntriesInput.Hide()
	captureIntervalMsInput.Hide()
	darkModeBtn.Hide()

	darkModeChanged := false
	darkModeBtn.SetValue(appConf.DarkMode)

	darkModeBtn.SetCallback(func() {
		appConf.DarkMode = !appConf.DarkMode
		if !darkModeChanged {
			fltk.MessageBox("App Restart Required", "In order for this to take effect, the theme change won't take effect until the application is restarted.")
			darkModeChanged = true
		}
		darkModeBtn.SetValue(appConf.DarkMode)
	})

	captureIntervalMsInput.SetCallback(func() {
		interval, err := strconv.ParseInt(captureIntervalMsInput.Value(), 10, 64)
		if err != nil {
			fltk.MessageBox("Invalid", fmt.Sprintf("Failed to validate your input: %v", err.Error()))
			return
		}

		if interval < 30 {
			fltk.MessageBox("Too small", "To avoid excessive CPU usage, the capture interval must be above 30.")
			return
		}

		appConf.CaptureIntervalMS = int(interval)
	})

	maxEntriesInput.SetCallback(func() {
		mi, err := strconv.ParseInt(maxEntriesInput.Value(), 10, 64)
		if err != nil {
			fltk.MessageBox("Invalid", fmt.Sprintf("Failed to validate your input: %v", err.Error()))
			return
		}

		if mi <= 0 {
			fltk.MessageBox("Too small", "You must set a value greater than 0.")
			return
		}

		appConf.MaxEntries = int(mi)
	})

	settingsBtn.SetCallback(func() {
		switchPage(PAGE_SETTINGS)
		responsive(win)
	})

	backBtn.SetCallback(func() {
		switchPage(PAGE_MAIN)
		responsive(win)
	})

	saveBtn.SetCallback(func() {
		err := saveConfig(configFilePath, &appConf)
		if err != nil {
			fltk.MessageBox("Error", fmt.Sprintf("failed to save config: %v", err.Error()))
		}

		fltk.MessageBox("Saved", "Saved successfully.")
	})

	reconstruct := func() {
		logBrowser.Clear()
		l := len(appConf.Log)
		if l > appConf.MaxEntries {
			appConf.Log = appConf.Log[1 : appConf.MaxEntries+1]
			l = len(appConf.Log)
		}
		// initialize with the previously stored entries
		if l > 0 {
			i := 1
			for j := l - 1; j >= 0; j-- {
				v := strings.ReplaceAll(appConf.Log[j].Value, "\n", "\\n")
				// v = fmt.Sprintf("%v.  %v", j+1, v[:minz(len(v)-1, 200)])
				v = v[0:minz(len(v), 200)]
				v = obscure(v, appConf.Secrets)
				v = fmt.Sprintf("%v.  %v", i, v)
				logBrowser.Add(v)
				_ = logBrowser.SetSelected(j+1, appConf.Log[j].Selected)
				i++
			}
		}
	}

	reconstruct()

	addEntry := func(entry string) {
		l := len(appConf.Log)

		if l > 0 && appConf.Log[l-1].Value == entry {
			return
		}

		appConf.Log = append(appConf.Log, ClipboardEntry{
			Value:    entry,
			Selected: false,
		})

		reconstruct()
	}

	captureClipboard := func() {
		latest, err := clipboard.ReadAll()
		if err != nil {
			// Logf("failed to read clipboard: %v, ", err.Error())
			return
		}

		addEntry(latest)
	}

	logBrowser.SetCallback(func() {
		i := logBrowser.Value()
		j := len(appConf.Log) - i
		logBrowser.SetTooltip(appConf.Log[j].Value)
	})

	copyAction := func() {
		total := 0
		copyStr := new(strings.Builder)
		l := len(appConf.Log)
		itemsCopied := 0
		for i := 1; i <= appConf.MaxEntries; i++ {
			if i > l {
				break
			}

			// note that entries in the browser are reversed, relative to
			// appConf.Log's items.

			j := l - i
			if j < 0 {
				break
			}

			if logBrowser.IsSelected(i) {
				copyStr.WriteString(fmt.Sprintf("%v\n", appConf.Log[j].Value))
				itemsCopied++
			}

			appConf.Log[j].Selected = false
			logBrowser.SetSelected(i, false)
			total += len(appConf.Log[j].Value)
		}

		result := copyStr.String()
		if result == "" {
			log.Println("nothing to copy")
			return
		}

		// remove the final trailing newline
		result = strings.TrimSuffix(result, "\n")

		msg := fmt.Sprintf("%v/%v items copied, %v bytes (%v bytes in history)", itemsCopied, l, len(result), total)
		logBrowser.SetLabel(msg)
		log.Println(msg)

		err := clipboard.WriteAll(result)
		if err != nil {
			fltk.MessageBox("Error", fmt.Sprintf("Failed to write to clipboard: %v", err.Error()))
			return
		}
	}

	delAction := func() {
		l := len(appConf.Log)
		toDel := []int{}
		for i := 1; i <= appConf.MaxEntries; i++ {
			if i > l {
				break
			}

			// note that entries in the browser are reversed, relative to
			// appConf.Log's items.

			j := l - i
			if j < 0 {
				break
			}

			if logBrowser.IsSelected(i) {
				toDel = append(toDel, j)
			}
		}

		sort.Ints(toDel)
		slices.Reverse(toDel)

		msg := fmt.Sprintf("%v/%v items deleted", len(toDel), l)
		logBrowser.SetLabel(msg)
		logBrowser.Redraw()

		// i = 40, len = 65
		// [:40], [41:]...
		for _, i := range toDel {
			log.Printf("deleting %v (%v left)", i, len(appConf.Log))
			appConf.Log = append(appConf.Log[:i], appConf.Log[i+1:]...)
		}

		reconstruct()
	}

	saveAction := func() {
		err := saveConfig(configFilePath, &appConf)
		if err != nil {
			fltk.MessageBox("Error", fmt.Sprintf("Failed to save config: %v", err.Error()))
			return
		}
		fltk.MessageBox("Success", "Saved configuration successfully.")
	}

	selectAllAction := func() {
		scrollPos := logBrowser.TopLine()
		for i := range appConf.Log {
			appConf.Log[i].Selected = true
			logBrowser.SetSelected(i+1, true)
		}
		// reconstruct()
		_ = logBrowser.SetTopLine(scrollPos)
	}

	// homeAction := func() {
	// 	_ = logBrowser.SetTopLine(0)
	// 	logBrowser.SetSelected(1, true)
	// }

	// endAction := func() {
	// 	_ = logBrowser.SetBottomLine(0)
	// 	logBrowser.SetSelected(len(appConf.Log), true)
	// }

	gracefulExit := func() {
		Log("closing app and saving config, please wait a moment...")
		err := saveConfig(configFilePath, &appConf)
		if err != nil {
			log.Printf("failed to save config: %v", err.Error())
		}

		Log("done, exiting now.")
		os.Exit(0)
	}
	// invisible menu that receives keyboard shortcuts
	topMenu := fltk.NewMenuBar(0, 0, 0, 0)
	topMenu.AddEx("Copy", fltk.CTRL+'c', copyAction, 0)
	topMenu.AddEx("Delete", fltk.DELETE, delAction, 0)
	topMenu.AddEx("Save", fltk.CTRL+'s', saveAction, 0)
	topMenu.AddEx("Select All", fltk.CTRL+'a', selectAllAction, 0)
	topMenu.AddEx("Quit", fltk.CTRL+'q', gracefulExit, 0)
	// topMenu.AddEx("Home", fltk.HOME, homeAction, 0)
	// topMenu.AddEx("End", fltk.END, endAction, 0)
	copyBtn.SetCallback(copyAction)
	deleteBtn.SetCallback(delAction)

	go func() {
		for {
			captureClipboard()

			time.Sleep(time.Duration(appConf.CaptureIntervalMS) * time.Millisecond)
		}
	}()

	win.SetCallback(gracefulExit)

	fltk.EnableTooltips()
	fltk.SetTooltipDelay(0.1)

	if portrait {
		win.Resize(0, 0, WIDTH_PORTRAIT*3, HEIGHT_PORTRAIT*3)
	} else {
		win.Resize(0, 0, WIDTH_LANDSCAPE*3, HEIGHT_LANDSCAPE*3)
	}

	win.SetResizeHandler(func() {
		responsive(win)
	})

	theme(appConf.DarkMode)

	responsive(win)

	win.SetXClass("gfltkclip")

	win.End()
	win.Show()

	go fltk.Run()

	// Create a channel to receive OS signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-signalChan

	gracefulExit()
}
