// SPDX-License-Identifier: BSD-3-Clause
// Copyright(c) 2022 Intel Corporation

package main

import (
	"fmt"
	"sync"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/rivo/tview"

	cz "github.com/KeithWiles/go-pktgen/pkgs/colorize"
	tab "github.com/KeithWiles/go-pktgen/pkgs/taborder"
	tlog "github.com/KeithWiles/go-pktgen/pkgs/ttylog"
)

// PageSingleMode - Data for main page information
type PageSingleMode struct {
	topFlex      *tview.Flex
	singleConfig *tview.Table
	singleStats  *tview.Table
	singleSizes  *tview.Table
	configOnce   sync.Once
	statsOnce    sync.Once
	selectOnce   sync.Once
}

const (
	singlePanelName string = "Single"
)

func init() {
	tlog.Register("SingleModeLogID")
}

// Printf - send message to the ttylog interface
func (ps *PageSingleMode) Printf(format string, a ...interface{}) {
	tlog.Log("SingleModeLogID", fmt.Sprintf("%T.", ps)+format, a...)
}

// setupSingleMode - setup and init the sysInfo page
func setupSingleMode() *PageSingleMode {

	ps := &PageSingleMode{}

	return ps
}

// SingleModePanelSetup setup the main cpu page
func SingleModePanelSetup(nextSlide func()) (pageName string, content tview.Primitive) {

	ps := setupSingleMode()

	to := tab.New(singlePanelName, pktgen.app)

	flex0 := tview.NewFlex().SetDirection(tview.FlexRow)
	flex1 := tview.NewFlex().SetDirection(tview.FlexRow)

	TitleBox(flex0)

	ps.singleStats = CreateTableView(flex1, "Stats (s)", tview.AlignLeft, 12, 4, true).
		SetSelectable(false, false).
		SetFixed(1, 1).
		SetSeparator(tview.Borders.Vertical)

	ps.singleSizes = CreateTableView(flex1, "Size Stats (S)", tview.AlignLeft, 12, 4, true).
		SetSelectable(false, false).
		SetFixed(2, 1).
		SetSeparator(tview.Borders.Vertical)

	ps.singleConfig = CreateTableView(flex1, "Single Configuration (c)", tview.AlignLeft, 11, 1, true).
		SetSelectable(true, false).
		SetFixed(1, 1).
		SetSeparator(tview.Borders.Vertical)

	flex0.AddItem(flex1, 0, 1, true)

	to.Add("singleStats", ps.singleStats, 's')
	to.Add("singleSizes", ps.singleSizes, 'S')
	to.Add("singleConfig", ps.singleConfig, 'c')
	to.SetInputDone()

	ps.topFlex = flex0

	pktgen.timers.Add(singlePanelName, func(step int, ticks uint64) {
		if ps.topFlex.HasFocus() {
			ps.selectOnce.Do(func() {
				to.SetInputFocus('c') // Select the config view after we have focus
			})
			pktgen.app.QueueUpdateDraw(func() {
				ps.displaySingleMode(step, ticks)
			})
		}
	})

	return singlePanelName, ps.topFlex
}

// Callback timer routine to display the cpuinfo panel
func (ps *PageSingleMode) displaySingleMode(step int, ticks uint64) {

	switch step {
	case 0:

	case 2:
		ps.configTable(ps.singleConfig)
		ps.displayStats(ps.singleStats)
		ps.displaySizes(ps.singleSizes)
	}
}

func (ps *PageSingleMode) configTable(table *tview.Table) {

	row := 0
	col := 0

	titles := []string{
		cz.Yellow("Port", 4),
		cz.Yellow("TX Count", 10),
		cz.Yellow("Rate", 5),
		cz.Yellow("Size", 4),
		cz.Yellow("Burst", 5),
		cz.Yellow("TTL", 3),
		cz.Yellow("sport", 5),
		cz.Yellow("dport", 5),
		cz.Yellow("IPType", 5),
		cz.Yellow("Proto", 5),
		cz.Yellow("VLAN", 4),
		cz.Yellow("IP Dst", 15),
		cz.Yellow("IP Src", 18),
		cz.Yellow("MAC Dst", 14),
		cz.Yellow("MAC Src", 14),
		cz.Yellow(" ", 4), // Extra field to allow scrolling horizontal
	}
	row = TableSetHeaders(ps.singleConfig, 0, 0, titles)

	for v := 0; v < pktgen.portCnt; v++{

		rowData := []string{
			cz.Yellow(v),
			cz.CornSilk("Forever"),
			cz.DeepPink(fmt.Sprintf("%3d%%", 100)),
			cz.LightCoral(64),
			cz.LightCoral(128),
			cz.LightCoral(64),
			cz.LightCoral(1234),
			cz.LightCoral(5678),
			cz.LightBlue("IPv4"),
			cz.LightBlue("UDP"),
			cz.Cyan(v+1),
			cz.CornSilk("198.18.0.1"),
			cz.CornSilk("198.18.0.1/24"),
			cz.Green("1234:5678:9000"),
			cz.Green("5678:1234:0001"),
			}
		for i, d := range rowData {
			if i == 0 {
				col = TableCellSelect(table, row, 0, d)
			} else {
				col = TableCellSet(table, row, col, d)
			}
		}
	
		row++
	}
	ps.configOnce.Do(func() {
		ps.singleConfig.ScrollToBeginning()
	})
}

func (ps *PageSingleMode) displayStats(table *tview.Table) {

	row := 0
	col := 0

	titles := []string{
		cz.Yellow(""),
		cz.Yellow(""),
		cz.Yellow(""),
		cz.Yellow(""),
		cz.Yellow("Rx/Tx"),
		cz.Yellow(""),
		cz.Yellow(""),
		cz.Yellow("Rx/Tx"),
		cz.Yellow("Total"),
		cz.Yellow("Total"),
		cz.Yellow("Total"),
		cz.Yellow("Total"),
	}
	row = TableSetHeaders(table, 0, 0, titles)

	titles = []string{
		cz.Yellow("Port", 4),
		cz.Yellow("Link State", 12),
		cz.Yellow("Rx pps", 12),
		cz.Yellow("Tx pps", 12),
		cz.Yellow("Mbits", 12),
		cz.Yellow("Rx Max", 12),
		cz.Yellow("Tx Max", 12),
		cz.Yellow("Errors", 12),
		cz.Yellow("Rx Pkts", 14),
		cz.Yellow("Tx Pkts", 14),
		cz.Yellow("Rx Mbits", 14),
		cz.Yellow("Tx Mbits", 14),
		cz.Yellow(" ", 2), // Extra field to allow scrolling horizontal
	}
	row = TableSetHeaders(table, row, 0, titles)

	p := message.NewPrinter(language.English)

	comma := func(n interface{}) string {
		return p.Sprintf("%d", n)
	}

	for v := 0; v < pktgen.portCnt; v++{

		rowData := []string{
			cz.Yellow(v),
			cz.LightYellow("UP-40000-FD"),
			cz.Cyan(comma(0)),
			cz.Cyan(comma(0)),
			cz.Wheat("0/0"),
			cz.Cyan(comma(0)),
			cz.Cyan(comma(0)),
			cz.Red("0/0"),
			cz.Cyan(comma(0)),
			cz.Cyan(comma(0)),
			cz.Cyan(comma(0)),
			cz.Cyan(comma(0)),
		}
		for i, d := range rowData {
			if i == 0 {
				col = TableCellSelect(table, row, 0, d)
			} else {
				col = TableCellSet(table, row, col, d)
			}
		}
		row++
	}

	ps.statsOnce.Do(func() {
		ps.singleStats.ScrollToBeginning()
	})
}

func (ps *PageSingleMode) displaySizes(table *tview.Table) {

	row := 0
	col := 0

	titles := []string{
		cz.Yellow(""),
		cz.Yellow(""),
		cz.Yellow(""),
		cz.Yellow("Pkt Sizes"),
		cz.Yellow("Pkt Sizes"),
		cz.Yellow("Pkt Sizes"),
		cz.Yellow("Pkt Sizes"),
		cz.Yellow("Pkt Sizes"),
		cz.Yellow(""),
		cz.Yellow(""),
	}
	row = TableSetHeaders(table, 0, 0, titles)

	titles = []string{
		cz.Yellow("Port", 4),
		cz.Yellow("Broadcast", 12),
		cz.Yellow("Multicast", 12),
		cz.Yellow("64", 12),
		cz.Yellow("128-255", 12),
		cz.Yellow("256-511", 12),
		cz.Yellow("512-1023", 12),
		cz.Yellow("1024-1518", 12),
		cz.Yellow("Runts/Jumbos", 14),
		cz.Yellow("ARPs/ICMPs", 14),
		cz.Yellow(" ", 6), // Extra field to allow scrolling horizontal
	}
	row = TableSetHeaders(table, row, 0, titles)

	for v := 0; v < pktgen.portCnt; v++ {

		rowData := []string{
			cz.Yellow(v),
			cz.Wheat(0),
			cz.GoldenRod(0),
			cz.Cyan(0),
			cz.Cyan(0),
			cz.Cyan(0),
			cz.Cyan(0),
			cz.Cyan(0),
			cz.DeepPink("0/0"),
			cz.Wheat("0/0"),
		}
		for i, d := range rowData {
			if i == 0 {
				col = TableCellSelect(table, row, 0, d)
			} else {
				col = TableCellSet(table, row, col, d)
			}
		}
		row++
	}

	ps.statsOnce.Do(func() {
		ps.singleStats.ScrollToBeginning()
	})
}
