// SPDX-License-Identifier: BSD-3-Clause
// Copyright(c) 2022 Intel Corporation

package main

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/shirou/gopsutil/cpu"

	cz "github.com/KeithWiles/go-pktgen/pkgs/colorize"
	"github.com/KeithWiles/go-pktgen/pkgs/meter"
	tab "github.com/KeithWiles/go-pktgen/pkgs/taborder"
	tlog "github.com/KeithWiles/go-pktgen/pkgs/ttylog"
)

// PageCPULoad - Data for main page information
type PageCPULoad struct {
	topFlex   *tview.Flex
	cpuInfo   *tview.TextView
	cpuLayout *tview.Table
	cpuInfo1  *tview.TextView
	cpuInfo2  *tview.TextView
	cpuInfo3  *tview.TextView
	tabOrder  *tab.Tab
	percent   []float64
	meter     *meter.Meter
}

const (
	cpuPanelName string = "CPU"
	labelWidth int = 14
)

func init() {
	tlog.Register("CPULoadLogID")
}

// Printf - send message to the ttylog interface
func (pg *PageCPULoad) Printf(format string, a ...interface{}) {
	tlog.Log("CPULoadLogID", fmt.Sprintf("%T.", pg)+format, a...)
}

// setupCPULoad - setup and init the sysInfo page
func setupCPULoad() *PageCPULoad {

	pg := &PageCPULoad{}

	return pg
}

// CPULoadPanelSetup setup
func CPULoadPanelSetup(pages *tview.Pages, nextSlide func()) (pageName string, content tview.Primitive) {

	pg := setupCPULoad()

	to := tab.New(cpuPanelName, pktgen.app)
	pg.tabOrder = to

	flex0 := tview.NewFlex().SetDirection(tview.FlexRow)
	flex1 := tview.NewFlex().SetDirection(tview.FlexColumn)
	flex2 := tview.NewFlex().SetDirection(tview.FlexColumn)

	TitleBox(flex0)

	pg.cpuInfo = CreateTextView(flex1, "CPU (c)", tview.AlignLeft, 0, 2, true)
	pg.cpuLayout = CreateTableView(flex1, "CPU Layout (l)", tview.AlignLeft, 0, 1, false)
	flex0.AddItem(flex1, 0, 1, true)

	pg.cpuInfo1 = CreateTextView(flex2, "CPU Load (1)", tview.AlignLeft, 0, 1, true)
	pg.cpuInfo2 = CreateTextView(flex2, "CPU Load (2)", tview.AlignLeft, 0, 1, false)
	pg.cpuInfo3 = CreateTextView(flex2, "CPU Load (3)", tview.AlignLeft, 0, 1, false)
	flex0.AddItem(flex2, 0, 4, true)

	to.Add("cpuInfo", pg.cpuInfo, 'c')
	to.Add("cpuLayout", pg.cpuLayout, 'l')

	to.Add("cpuInfo1", pg.cpuInfo1, '1')
	to.Add("cpuInfo2", pg.cpuInfo2, '2')
	to.Add("cpuInfo3", pg.cpuInfo3, '3')

	to.SetInputDone()

	pg.topFlex = flex0

	// Setup static pages
	pg.displayCPU(pg.cpuInfo)
	pg.displayLayout(pg.cpuLayout)

	percent, err := cpu.Percent(0, true)
	if err != nil {
		tlog.DoPrintf("Percent: %v\n", err)
	}
	pg.percent = percent

	pktgen.timers.Add(cpuPanelName, func(step int, ticks uint64) {
		if pg.topFlex.HasFocus() {
			pktgen.app.QueueUpdateDraw(func() {
				pg.displayCPULoad(step, ticks)
			})
		}
	})

	modal := tview.NewModal().
		SetText("This is the Help Box: cpuInfoHelp Thank you for asking for help! Press Esc to close.").
		AddButtons([]string{"Got it"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.HidePage("cpuInfoHelp")
		})
	AddModalPage("cpuInfoHelp", modal)

	flex0.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Rune()
		switch k {
        case '?':
            tlog.DoPrintf("Question Mark! HasPage(%v)\n", pages.HasPage("cpuInfoHelp"))
			pages.ShowPage("cpuInfoHelp")
		}
		return event
	})

	pg.meter = meter.New().
		SetWidth(func() int {
			_, _, width, _ := pg.cpuInfo1.GetInnerRect()

			return width
		}).
		SetDraw(func(mi *meter.Info) string {
			var str string = ""

			for _, l := range mi.Labels {

				if l.Fn == nil {
					l.Fn = cz.Default
				}
				str += l.Fn(l.Val)
			}
			str += fmt.Sprintf("[%s]\n", mi.Bar.Fn(mi.Bar.Val))
			return str
		}).
		SetRateLimits(0.0, 100.0)


	return cpuPanelName, pg.topFlex
}

// Callback timer routine to display the cpuinfo panel
func (pg *PageCPULoad) displayCPULoad(step int, ticks uint64) {

	switch step {
	case 0:
		percent, err := cpu.Percent(0, true)
		if err != nil {
			tlog.DoPrintf("Percent: %v\n", err)
		}
		pg.percent = percent

	case 2:
		pg.displayLoadData(pg.cpuInfo1, 1)
		pg.displayLoadData(pg.cpuInfo2, 2)
		pg.displayLoadData(pg.cpuInfo3, 3)
	}
}

// Display the CPU information
func (pg *PageCPULoad) displayCPU(view *tview.TextView) {
	str := ""

	cd := pktgen.cpuData
	str += fmt.Sprintf("CPU   Vendor   : %s\n", cz.GoldenRod(cd.CpuInfo(0).VendorID, -14))
	str += fmt.Sprintf("      Model    : %s\n\n", cz.MediumSpringGreen(cd.CpuInfo(0).ModelName))
	str += fmt.Sprintf("Cores Logical  : %s\n", cz.Yellow(cd.NumLogicalCores(), -6))
	str += fmt.Sprintf("      Physical : %s\n", cz.Yellow(cd.NumPhysicalCores(), -6))
	str += fmt.Sprintf("      Threads  : %s\n", cz.Yellow(cd.NumHyperThreads(), -6))
	str += fmt.Sprintf("      Sockets  : %s\n", cz.Yellow(cd.NumSockets()))

	view.SetText(str)
	view.ScrollToBeginning()
}

// Build up a string for displaying the CPU layout window
func buildStr(a []uint16, width int) string {

	str := "{"

	for k, v := range a {
		str += cz.Green(v, width)
		if k < (len(a) - 1) {
			str += " /"
		}
	}

	return str + " }"
}

// Display the CPU layout data
func (pg *PageCPULoad) displayLayout(view *tview.Table) {

	cd := pktgen.cpuData

	str := cz.LightBlue(" Core", -5)
	tableCell := tview.NewTableCell(cz.YellowGreen(str)).
		SetAlign(tview.AlignLeft).
		SetSelectable(false)
	view.SetCell(0, 0, tableCell)

	for k, s := range cd.Sockets() {
		str = cz.LightBlue(fmt.Sprintf("Socket %d", s))
		tableCell := tview.NewTableCell(cz.YellowGreen(str)).
			SetAlign(tview.AlignCenter).
			SetSelectable(false)
		view.SetCell(0, k+1, tableCell)
	}

	row := int16(1)

	pg.Printf("numPhysical %d, numSockets %d\n", cd.NumPhysicalCores(), cd.NumSockets())
	pg.Printf("cd.Cores = %v\n", cd.Cores())
	for _, cid := range cd.Cores() {
		col := int16(0)

		tableCell := tview.NewTableCell(cz.Red(cid, 4)).
			SetAlign(tview.AlignLeft).
			SetSelectable(false)
		view.SetCell(int(row), int(col), tableCell)

		pg.Printf("cid %d\n", cid)
		for sid := int16(0); sid < cd.NumSockets(); sid++ {
			pg.Printf("  sid %d\n", sid)
			key := uint16(sid<<uint16(8)) | cid
			v, ok := cd.CoreMapItem(key)
			if ok {
				str = fmt.Sprintf(" %s", buildStr(v, 3))
			} else {
				str = fmt.Sprintf(" %s", strings.Repeat(".", 10))
			}
			tableCell := tview.NewTableCell(cz.YellowGreen(str)).
				SetAlign(tview.AlignLeft).
				SetSelectable(false)
			view.SetCell(int(row), int(col+1), tableCell)
			col++
		}
		row++
	}
	view.ScrollToBeginning()
}

// Grab the percent load data and display the meters
func (pg *PageCPULoad) displayLoadData(view *tview.TextView, flg int) {

	cd := pktgen.cpuData
	num := int16(cd.NumLogicalCores()/3) + 1

	switch flg {
	case 1:
		pg.displayLoad(pg.percent, 0, num, view)
	case 2:
		pg.displayLoad(pg.percent, num, num*int16(2), view)
	case 3:
		pg.displayLoad(pg.percent, num*int16(2), cd.NumLogicalCores(), view)
	}
}

// Display the load meters
func (pg *PageCPULoad) displayLoad(percent []float64, start, end int16, view *tview.TextView) {

	_, _, width, _ := view.GetInnerRect()

	width -= labelWidth
	if width <= 0 {
		return
	}
	str := ""

	str += fmt.Sprintf("%s\n", cz.Orange("Core Percent          Load Meter"))

	for i := start; i < end; i++ {
		str += pg.meter.Draw(percent[i], &meter.Info{
			Labels: []*meter.LabelInfo{
				{Val: fmt.Sprintf("%3d", i), Fn: nil},
				{Val: ":", Fn: nil},
				{Val: fmt.Sprintf("%5.1f", percent[i]), Fn: cz.Red},
				{Val: "%", Fn: nil},
			},
			Bar: &meter.LabelInfo{Val: "", Fn: cz.MediumSpringGreen},
		})
	}
	str = str[:len(str) - 1]	// Strip the last newline character

	view.SetText(str)
	view.ScrollToBeginning()
}
