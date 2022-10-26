// SPDX-License-Identifier: BSD-3-Clause
// Copyright(c) 2022 Intel Corporation

package main

import (
	"fmt"
	"math"
	"net"
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
	singlePerf   *tview.TextView
	configOnce   sync.Once
	statsOnce    sync.Once
	selectOnce   sync.Once
	ppsPerPort      []float64
}

const (
	singlePanelName string = "Single"
)

func init() {
	tlog.Register("SingleModeLogID")

	for i := 0; i < pktgen.portCnt; i++ {
		pktgen.single[i] = SinglePacketConfig{
			TxCount:     0,
			PercentRate: 100.0,
			PktSize:     64,
			BurstCount:  128,
			TimeToLive:  64,
			SrcPort:     1245,
			DstPort:     5678,
			IPType:      "IPv4",
			ProtoType:   "UDP",
			VlanId:      1,
			DstIP:       net.IPNet{IP: net.IPv4(198, 18, 1, 1), Mask: net.CIDRMask(0, 32)},
			SrcIP:       net.IPNet{IP: net.IPv4(198, 18, 0, 1), Mask: net.CIDRMask(24, 32)},
			DstMAC:      []byte{0x12, 0x34, 0x45, 0x67, 0x89, 00},
			SrcMAC:      []byte{0x12, 0x34, 0x45, 0x67, 0x89, 01},
		}
	}
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

	ps.singlePerf = CreateTextView(flex1, "Performance (p)", tview.AlignLeft, 0, 1, true)
	ps.ppsPerPort = make([]float64, 8)

	flex0.AddItem(flex1, 0, 1, true)

	to.Add("singleStats", ps.singleStats, 's')
	to.Add("singleSizes", ps.singleSizes, 'S')
	to.Add("singleConfig", ps.singleConfig, 'c')
	to.Add("singlePerf", ps.singlePerf, 'p')
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
		ps.displayPerf(ps.singlePerf)
	}
}

func (ps *PageSingleMode) configTable(table *tview.Table) {

	row := 0
	col := 0

	single := pktgen.single

	titles := []string{
		cz.Yellow("Port", 4),
		cz.Yellow("TX Count", 10),
		cz.Yellow("% Rate", 6),
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

	for v := 0; v < pktgen.portCnt; v++ {

		rate := func() string {
			if pktgen.single[v].TxCount == 0 {
				return "Forever"
			}
			return fmt.Sprintf("%v", single[v].TxCount)
		}
		rowData := []string{
			cz.Yellow(v),
			cz.CornSilk(rate()),
			cz.DeepPink(fmt.Sprintf("%v", single[v].PercentRate)),
			cz.LightCoral(single[v].PktSize),
			cz.LightCoral(single[v].BurstCount),
			cz.LightCoral(single[v].TimeToLive),
			cz.LightCoral(single[v].SrcPort),
			cz.LightCoral(single[v].DstPort),
			cz.LightBlue(single[v].IPType),
			cz.LightBlue(single[v].ProtoType),
			cz.Cyan(single[v].VlanId + 1),
			cz.CornSilk(single[v].DstIP.IP.String()),
			cz.CornSilk(single[v].SrcIP.String()),
			cz.Green(single[v].DstMAC.String()),
			cz.Green(single[v].SrcMAC.String()),
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

	for v := 0; v < pktgen.portCnt; v++ {

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

// Grab the load data and display the meters
func (ps *PageSingleMode) displayPerf(view *tview.TextView) {

	ps.displayLoad(ps.ppsPerPort, 0, int16(pktgen.portCnt), view)
}

// Display the load meters
func (ps *PageSingleMode) displayLoad(pps []float64, start, end int16, view *tview.TextView) {

	_, _, width, _ := view.GetInnerRect()

	width -= 14
	if width <= 0 {
		return
	}
	str := ""

	str += fmt.Sprintf("%s\n", cz.Orange("Port   PPS          Load Meter"))

	for i := start; i < end; i++ {
		str += ps.drawMeter(i, pps[i], width)
	}

	view.SetText(str)
	view.ScrollToBeginning()
}

// clamp the data to a fixed set of ranges
func clampPerf(x, low, high float64) float64 {

	if x > high {
		return high
	}
	if x < low {
		return low
	}
	return x
}

// Draw the meter for the load
func (ps *PageSingleMode) drawMeter(id int16, pps float64, width int) string {

	var total uint64 = 100

	p := clampPerf(float64(pps), 0, float64(total))
	if p > 0 {
		p = math.Ceil((float64(p) / float64(total)) * float64(width))
	}

	bar := make([]byte, width)

	for i := 0; i < width; i++ {
		if i <= int(p) {
			bar[i] = '|'
		} else {
			bar[i] = ' '
		}
	}
	str := fmt.Sprintf("%3d:%s [%s]\n",
		id, cz.Red(pps, 5, 1), cz.Yellow(string(bar)))

	return str
}
