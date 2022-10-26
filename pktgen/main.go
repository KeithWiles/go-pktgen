// SPDX-License-Identifier: BSD-3-Clause
// Copyright(c) 2022 Intel Corporation

package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	cz "github.com/KeithWiles/go-pktgen/pkgs/colorize"
	"github.com/KeithWiles/go-pktgen/pkgs/cpudata"
	tlog "github.com/KeithWiles/go-pktgen/pkgs/ttylog"
	flags "github.com/jessevdk/go-flags"

	"github.com/KeithWiles/go-pktgen/pkgs/etimers"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	// pmeVersion string
	pktgenVersion = "22.10.0"
)

// PanelInfo for title and primitive
type PanelInfo struct {
	title     string
	primitive tview.Primitive
}

// Panels is a function which returns the feature's main primitive and its title.
// It receives a "nextFeature" function which can be called to advance the
// presentation to the next slide.
type Panels func(nextPanel func()) (title string, content tview.Primitive)

// Pktgen for monitoring and system performance data
type Pktgen struct {
	version string             // Version of Pktgen
	app     *tview.Application // Application or top level application
	timers  *etimers.EventTimers
	cpuData *cpudata.CPUData
	panels  []PanelInfo
	portCnt int
}

// Options command line options
type Options struct {
	Ptty        string `short:"p" long:"ptty" description:"path to ptty /dev/pts/X"`
	Dbg         bool   `short:"D" long:"debug" description:"Wait 15 seconds (default) to connect debugger"`
	WaitTime    uint   `short:"W" long:"wait-time" description:"N seconds before startup" default:"15"`
	ShowVersion bool   `short:"V" long:"version" description:"Print out version and exit"`
	Verbose     bool   `short:"v" long:"Verbose output for debugging"`
}

// Global to the main package for the tool
var pktgen Pktgen
var options Options
var parser = flags.NewParser(&options, flags.Default)

const (
	mainLog = "MainLogID"
)

func buildPanelString(idx int) string {
	// Build the panel selection string at the bottom of the xterm and
	// highlight the selected tab/panel item.
	s := ""
	for index, p := range pktgen.panels {
		if index == idx {
			s += fmt.Sprintf("F%d:[orange::r]%s[white::-]", index+1, p.title)
		} else {
			s += fmt.Sprintf("F%d:[orange::-]%s[white::-]", index+1, p.title)
		}
		if (index + 1) < len(pktgen.panels) {
			s += " "
		}
	}
	return s
}

// Setup the tool's global information and startup the process info connection
func init() {
	tlog.Register(mainLog, true)

	pktgen = Pktgen{}
	pktgen.version = pktgenVersion

	// Create the main tview application.
	pktgen.app = tview.NewApplication()

	cd, err := cpudata.New()
	if err != nil {
		fmt.Printf("New CPU data failed: %s\n", err)
		return
	}
	pktgen.cpuData = cd
	pktgen.portCnt = 8
}

// Version number string
func Version() string {
	return pktgen.version
}

func main() {

	cz.SetDefault("ivory", "", 0, 2, "")

	_, err := parser.Parse()
	if err != nil {
		fmt.Printf("*** invalid arguments %v\n", err)
		os.Exit(1)
	}

	if len(options.Ptty) > 0 {
		err = tlog.Open(options.Ptty)
		if err != nil {
			fmt.Printf("ttylog open failed: %s\n", err)
			os.Exit(1)
		}
	}
	if options.ShowVersion {
		fmt.Printf("PME Version: %s\n", pktgen.version)
		return
	}

	tlog.Log(mainLog, "\n===== %s =====\n", PktgenInfo(false))
	fmt.Printf("\n===== %s =====\n", PktgenInfo(false))

	app := pktgen.app

	pktgen.timers = etimers.New(time.Second/4, 4)
	pktgen.timers.Start()

	panels := []Panels{
		SysInfoPanelSetup,
		SingleModePanelSetup,
		CPULoadPanelSetup,
	}

	// The bottom row has some info on where we are.
	info := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false)

	currentPanel := 0
	info.Highlight(strconv.Itoa(currentPanel))

	pages := tview.NewPages()

	previousPanel := func() {
		currentPanel = (currentPanel - 1 + len(panels)) % len(panels)
		info.Highlight(strconv.Itoa(currentPanel)).
			ScrollToHighlight()
		pages.SwitchToPage(strconv.Itoa(currentPanel))
		info.SetText(buildPanelString(currentPanel))
	}

	nextPanel := func() {
		currentPanel = (currentPanel + 1) % len(panels)
		info.Highlight(strconv.Itoa(currentPanel)).
			ScrollToHighlight()
		pages.SwitchToPage(strconv.Itoa(currentPanel))
		info.SetText(buildPanelString(currentPanel))
	}

	modal := tview.NewModal().
		SetText("This is the Help Box: it is a pop-up model that will display detailed information about the current panel's data. Thank you for asking for help! Press Esc to close.").
		AddButtons([]string{"Got it", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Got it" {
				app.SetFocus(pages.SwitchToPage(strconv.Itoa(currentPanel)))
			} else if buttonLabel == "Cancel" {
				/*if err := pages.SwitchToPage(strconv.Itoa(currentPanel)); err != nil {
					panic(err)
				}*/
			}
		})

	panelHelp := func() {
		if err := app.SetRoot(modal, false).SetFocus(modal).Run(); err != nil {
			panic(err)
		}
	}

	for index, f := range panels {
		title, primitive := f(nextPanel)
		pages.AddPage(strconv.Itoa(index), primitive, true, index == currentPanel)
		pktgen.panels = append(pktgen.panels, PanelInfo{title: title, primitive: primitive})
	}
	info.SetText(buildPanelString(0))

	// Create the main panel.
	panel := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(pages, 0, 1, true).
		AddItem(info, 1, 1, false)

	// Shortcuts to navigate the panels.
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlN {
			nextPanel()
		} else if event.Key() == tcell.KeyCtrlP {
			previousPanel()
		} else if event.Key() == tcell.KeyCtrlQ {
			app.Stop()
		} else if event.Key() == tcell.KeyCtrlH {
			panelHelp()
		} else {
			var idx int

			switch {
			case tcell.KeyF1 <= event.Key() && event.Key() <= tcell.KeyF19:
				idx = int(event.Key() - tcell.KeyF1)
			case event.Rune() == 'q':
				app.Stop()
			default:
				idx = -1
			}
			if idx != -1 {
				if idx < len(panels) {
					currentPanel = idx
					info.Highlight(strconv.Itoa(currentPanel)).ScrollToHighlight()
					pages.SwitchToPage(strconv.Itoa(currentPanel))
				}
				info.SetText(buildPanelString(idx))
			}
		}
		return event
	})

	setupSignals(syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)

	if options.Dbg {
		fmt.Printf("Waiting %d seconds for dlv to attach\n", options.WaitTime)
		time.Sleep(time.Second * time.Duration(options.WaitTime))
	}

	// Start the application.
	if err := app.SetRoot(panel, true).Run(); err != nil {
		panic(err)
	}

	tlog.Log(mainLog, "===== Done =====\n")
}

func setupSignals(signals ...os.Signal) {
	app := pktgen.app

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, signals...)
	go func() {
		sig := <-sigs

		tlog.Log(mainLog, "Signal: %v\n", sig)
		time.Sleep(time.Second)

		app.Stop()
		os.Exit(1)
	}()
}
