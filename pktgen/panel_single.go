// SPDX-License-Identifier: BSD-3-Clause
// Copyright(c) 2022 Intel Corporation

package main

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"sync"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	cz "github.com/KeithWiles/go-pktgen/pkgs/colorize"
	"github.com/KeithWiles/go-pktgen/pkgs/meter"
	tab "github.com/KeithWiles/go-pktgen/pkgs/taborder"
	tlog "github.com/KeithWiles/go-pktgen/pkgs/ttylog"
)

// PageSingleMode - Data for main page information
type PageSingleMode struct {
	topFlex       *tview.Flex
	singleConfig  *tview.Table
	singleStats   *tview.Table
	singleSizes   *tview.Table
	singlePerf    *tview.TextView
	configOnce    sync.Once
	statsOnce     sync.Once
	sizesOnce     sync.Once
	rxPercentRate []float64
	txPercentRate []float64
	configForms   []*tview.Flex
	currentPort   int
	to            *tab.Tab
	meter         *meter.Meter
}

const (
	singlePanelName  string = "Single"
	singleInfoHelp   string = "singleInfoHelp"
	singlePortConfig string = "singlePortConfig"
)

func init() {
	tlog.Register("SingleModeLogID")

	for i := 0; i < pktgen.portCnt; i++ {
		pktgen.single[i] = &SinglePacketConfig{
			PortIndex:   i,
			TxCount:     0,
			PercentRate: 100.0,
			PktSize:     64,
			BurstCount:  128,
			TimeToLive:  64,
			SrcPort:     1245,
			DstPort:     5678,
			PType:       "IPv4",
			ProtoType:   "UDP",
			VlanId:      1,
			DstIP:       net.IPNet{IP: net.IPv4(198, 18, 1, 1), Mask: net.CIDRMask(0, 32)},
			SrcIP:       net.IPNet{IP: net.IPv4(198, 18, 0, 1), Mask: net.CIDRMask(24, 32)},
			DstMAC:      []byte{0x12, 0x34, 0x45, 0x67, 0x89, 00},
			SrcMAC:      []byte{0x12, 0x34, 0x45, 0x67, 0x89, 01},
			TxState:     false,
		}
	}
}

// setupSingleMode - setup and init the sysInfo page
func setupSingleMode() *PageSingleMode {

	ps := &PageSingleMode{}

	return ps
}

func (ps *PageSingleMode) setupConfigForm(pages *tview.Pages, port int) *tview.Flex {

	pg := fmt.Sprintf("%v-%v", singlePortConfig, port)

	form := tview.NewForm().
		SetItemPadding(1).
		SetHorizontal(false).
		SetFieldTextColor(tcell.ColorBlack).
		SetFieldBackgroundColor(tcell.ColorBlue).
		SetItemPadding(0).
		SetCancelFunc(func() {
			pages.HidePage(pg)
			ps.to.SetInputFocus('c')
		})

	form.SetTitleAlign(tview.AlignLeft).SetRect(0, 0, 35, 21)

	sc := *pktgen.single[port]

	form.AddInputField("Port ID  :", strconv.Itoa(int(sc.PortIndex)), 2,
		func(textToCheck string, lastChar rune) bool {
			return false
		}, nil)

	form.AddInputField("TxCount  :", strconv.Itoa(int(sc.TxCount)), 15,
		func(textToCheck string, lastChar rune) bool {
			return len(textToCheck) <= 15 && acceptNumber(textToCheck, lastChar)
		}, func(text string) {
			parseNumberUint64(text, &sc.TxCount)
		})

	form.AddInputField("Rate     :", strconv.FormatFloat(sc.PercentRate, 'f', 2, 64), 6,
		func(textToCheck string, lastChar rune) bool {
			return len(textToCheck) <= 6 && acceptFloat(textToCheck, lastChar)
		}, func(text string) {
			if err := parseNumberFloat64(text, &sc.PercentRate); err == nil {
				if sc.PercentRate == 0 || sc.PercentRate > 100.00 {
					sc.PercentRate = 100.00
				}
			}
		})

	form.AddInputField("PktSize  :", strconv.Itoa(int(sc.PktSize)), 5,
		func(textToCheck string, lastChar rune) bool {
			return len(textToCheck) <= 5 && acceptNumber(textToCheck, lastChar)
		}, func(text string) {
			if err := parseNumberUint16(text, &sc.PktSize); err == nil {
				if sc.PktSize < 64 {
					sc.PktSize = 64
				} else if sc.PktSize > 1522 {
					sc.PktSize = 1522
				}
			}
		})

	form.AddInputField("Burst    :", strconv.Itoa(int(sc.BurstCount)), 3,
		func(textToCheck string, lastChar rune) bool {
			return len(textToCheck) <= 3 && acceptNumber(textToCheck, lastChar)
		}, func(text string) {
			if err := parseNumberUint16(text, &sc.BurstCount); err == nil {
				if sc.BurstCount < 32 {
					sc.BurstCount = 32
				} else if sc.BurstCount > 256 {
					sc.BurstCount = 256
				}
			}
		})

	form.AddInputField("TTL      :", strconv.Itoa(int(sc.TimeToLive)), 3,
		func(textToCheck string, lastChar rune) bool {
			return len(textToCheck) <= 3 && acceptNumber(textToCheck, lastChar)
		}, func(text string) {
			if err := parseNumberUint16(text, &sc.TimeToLive); err == nil {
				if sc.TimeToLive > 255 {
					sc.TimeToLive = 64
				}
			}
		})

	form.AddInputField("SrcPort  :", strconv.Itoa(int(sc.SrcPort)), 5,
		func(textToCheck string, lastChar rune) bool {
			return len(textToCheck) <= 5 && acceptNumber(textToCheck, lastChar)
		}, func(text string) {
			parseNumberUint16(text, &sc.SrcPort)
		})

	form.AddInputField("DstPort  :", strconv.Itoa(int(sc.DstPort)), 5,
		func(textToCheck string, lastChar rune) bool {
			return len(textToCheck) <= 5 && acceptNumber(textToCheck, lastChar)
		}, func(text string) {
			parseNumberUint16(text, &sc.DstPort)
		})

	form.AddDropDown("PType    :", []string{"IPv4", "IPv6", "ICMP"}, 0,
		func(option string, optionIndex int) {
			sc.PType = option
		})

	form.AddDropDown("Protocol :", []string{"UDP", "TCP"}, 0,
		func(option string, optionIndex int) {
			sc.ProtoType = option
		})

	form.AddInputField("VlanID   :", strconv.Itoa(int(sc.VlanId)), 4,
		func(textToCheck string, lastChar rune) bool {
			return len(textToCheck) <= 4 && acceptNumber(textToCheck, lastChar)
		}, func(text string) {
			if err := parseNumberUint16(text, &sc.VlanId); err == nil {
				if sc.VlanId == 0 {
					sc.VlanId = 1
				} else if sc.VlanId > 4095 {
					sc.VlanId = 4095
				}
			}
		})

	form.AddInputField("DstIP    :", sc.DstIP.String(), 15,
		func(textToCheck string, lastChar rune) bool {
			return len(textToCheck) <= 15 && acceptIPv4(textToCheck, lastChar)
		}, func(text string) {
			ip := net.ParseIP(text)
			sc.DstIP.IP = ip
			sc.DstIP.Mask = ip.DefaultMask()
		})

	form.AddInputField("SrcIP    :", sc.SrcIP.String(), 18,
		func(textToCheck string, lastChar rune) bool {
			return len(textToCheck) <= 18 && acceptIPv4CiDR(textToCheck, lastChar)
		}, func(text string) {
			ip := net.ParseIP(text)
			sc.SrcIP.IP = ip
			sc.SrcIP.Mask = ip.DefaultMask()
		})

	form.AddInputField("DstMAC   :", sc.DstMAC.String(), 18,
		func(textToCheck string, lastChar rune) bool {
			return len(textToCheck) <= 18 && acceptMac(textToCheck, lastChar)
		}, func(text string) {
			mac, err := net.ParseMAC(text)
			if err == nil {
				sc.DstMAC = mac
			}
		})

	form.AddInputField("SrcMAC   :", sc.SrcMAC.String(), 18,
		func(textToCheck string, lastChar rune) bool {
			return len(textToCheck) <= 18 && acceptMac(textToCheck, lastChar)
		}, func(text string) {
			mac, err := net.ParseMAC(text)
			if err == nil {
				sc.SrcMAC = mac
			}
		})

	form.AddButton("Save", func() {
		pktgen.single[port] = &sc
		pages.HidePage(pg)
		ps.to.SetInputFocus('c')
	}).SetButtonTextColor(tcell.ColorBlack)

	form.AddButton("Cancel", func() {
		pages.HidePage(pg)
		ps.to.SetInputFocus('c')
	}).SetButtonTextColor(tcell.ColorBlack)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(form, 0, 1, true)

	flex.SetTitle(TitleColor("Edit Port")).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetRect(20, 3, 35, 21)

	AddModalPage(pg, flex)

	return flex
}

// SingleModePanelSetup setup
func SingleModePanelSetup(pages *tview.Pages, nextSlide func()) (pageName string, content tview.Primitive) {

	ps := setupSingleMode()

	ps.to = tab.New(singlePanelName, pktgen.app)

	flex0 := tview.NewFlex().SetDirection(tview.FlexRow)
	flex1 := tview.NewFlex().SetDirection(tview.FlexRow)

	TitleBox(flex0)

	ps.singleConfig = CreateTableView(flex1, "Configuration (c) Start/Stop-r/s, Start/Stop All-R/S, Edit-e",
		tview.AlignLeft, 11, 0, true).
		SetSelectable(true, false).
		SetFixed(1, 1).
		SetSelectionChangedFunc(func(row, col int) {
			// Adjust for not selecting the first column in a row. table.go will select the next row.
			// on a mouse click
			if col > 0 {
				ps.singleConfig.Select(row, 0)
			}
		}).
		SetSeparator(tview.Borders.Vertical)

	ps.singleStats = CreateTableView(flex1, "Stats (1)", tview.AlignLeft, 0, 1, true).
		SetSelectable(false, false).
		SetFixed(1, 1).
		SetSeparator(tview.Borders.Vertical)

	ps.singleSizes = CreateTableView(flex1, "Size Stats (2)", tview.AlignLeft, 0, 1, true).
		SetSelectable(false, false).
		SetFixed(1, 1).
		SetSeparator(tview.Borders.Vertical)

	ps.singlePerf = CreateTextView(flex1, "Performance (p)", tview.AlignLeft, 18, 0, true)

	ps.rxPercentRate = make([]float64, 8)
	ps.txPercentRate = make([]float64, 8)

	flex0.AddItem(flex1, 0, 1, true)

	ps.to.Add("singleConfig", ps.singleConfig, 'c')
	ps.to.Add("singleStats", ps.singleStats, '1')
	ps.to.Add("singleSizes", ps.singleSizes, '2')
	ps.to.Add("singlePerf", ps.singlePerf, 'p')
	ps.to.SetInputDone()

	ps.topFlex = flex0

	pktgen.timers.Add(singlePanelName, func(step int, ticks uint64) {
		if ps.topFlex.HasFocus() {
			pktgen.app.QueueUpdateDraw(func() {
				ps.displaySingleMode(step, ticks)
			})
		}
	})

	modal := tview.NewModal().
		SetText("This is the Help Box: singleInfoHelp Thank you for asking for help! Press Esc to close.").
		AddButtons([]string{"Got it"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.HidePage(singleInfoHelp)
		})
	AddModalPage(singleInfoHelp, modal)

	for port := 0; port < pktgen.portCnt; port++ {
		f := ps.setupConfigForm(pages, port)
		ps.configForms = append(ps.configForms, f)
	}

	ps.singleConfig.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		ps.currentPort, _ = ps.singleConfig.GetSelection()
		ps.currentPort--

		sc := pktgen.single[ps.currentPort]

		k := event.Rune()
		switch k {
		case 'e':
			pages.ShowPage(fmt.Sprintf("%v-%v", singlePortConfig, ps.currentPort))
		case 'r':
			sc.TxState = true
		case 'R':
			for i := 0; i < pktgen.portCnt; i++ {
				pktgen.single[i].TxState = true
			}
		case 's':
			sc.TxState = false
		case 'S':
			for i := 0; i < pktgen.portCnt; i++ {
				pktgen.single[i].TxState = false
			}
		default:
			ps.to.SetInputFocus(k)
		}
		return event
	})
	flex0.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Rune()
		switch k {
		case '?':
			pages.ShowPage(singleInfoHelp)
		default:
		}
		return event
	})

	ps.meter = meter.New().
		SetWidth(func() int {
			_, _, width, _ := ps.singlePerf.GetInnerRect()

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

	return singlePanelName, ps.topFlex
}

// Callback timer routine to display the panels
func (ps *PageSingleMode) displaySingleMode(step int, ticks uint64) {

	switch step {
	case 0:
		ps.pullStats()

	case 2:
		ps.configTable()
		ps.displayStats()
		ps.displaySizes()
		ps.displayPerf()
	}
}

func (ps *PageSingleMode) pullStats() {

	for port := 0; port < pktgen.portCnt; port++ {
		ps.rxPercentRate[port] = float64(rand.Intn(101))
		ps.txPercentRate[port] = float64(rand.Intn(101))
	}

}

func (ps *PageSingleMode) configTable() {

	table := ps.singleConfig
	row := 0
	col := 0

	titles := []string{
		cz.Yellow("Port", 5),
		cz.Yellow("TX Count", 8),
		cz.Yellow("% Rate", 7),
		cz.Yellow("Size", 4),
		cz.Yellow("Burst", 5),
		cz.Yellow("TTL", 4),
		cz.Yellow("sport", 5),
		cz.Yellow("dport", 5),
		cz.Yellow("PType", 5),
		cz.Yellow("Proto", 5),
		cz.Yellow("VLAN", 4),
		cz.Yellow("IP Dst"),
		cz.Yellow("IP Src"),
		cz.Yellow("MAC Dst", 14),
		cz.Yellow("MAC Src", 14),
		cz.Yellow(" ", 16), // Extra field to allow scrolling horizontal
	}
	row = TableSetHeaders(ps.singleConfig, 0, 0, titles)

	state := func(port int, state bool) string {
		var s string = "   "

		if state {
			s = ">> "
		}
		return fmt.Sprintf("%s%s", cz.DeepPink(s), cz.Yellow(port, 2))
	}

	txCount := func(c uint64) string {
		if c == 0 {
			return "Forever"
		}
		p := message.NewPrinter(language.English)
		return p.Sprintf("%v", c)
	}

	for v := 0; v < pktgen.portCnt; v++ {
		single := pktgen.single[v]

		rowData := []string{
			state(single.PortIndex, single.TxState),
			cz.CornSilk(txCount(single.TxCount)),
			cz.DeepPink(strconv.FormatFloat(single.PercentRate, 'f', 2, 64)),
			cz.LightCoral(single.PktSize),
			cz.LightCoral(single.BurstCount),
			cz.LightCoral(single.TimeToLive),
			cz.LightCoral(single.SrcPort),
			cz.LightCoral(single.DstPort),
			cz.LightBlue(single.PType),
			cz.LightBlue(single.ProtoType),
			cz.Cyan(single.VlanId),
			cz.CornSilk(single.DstIP.IP.String()),
			cz.CornSilk(single.SrcIP.String()),
			cz.Green(single.DstMAC.String()),
			cz.Green(single.SrcMAC.String()),
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

func (ps *PageSingleMode) displayStats() {

	table := ps.singleStats
	row := 0
	col := 0

	titles := []string{
		cz.Yellow("Port", 4),
		cz.Yellow("Link State", 12),
		cz.Yellow("Rx pps", 8),
		cz.Yellow("Tx pps", 8),
		cz.Yellow("Rx/Tx Mbits", 12),
		cz.Yellow("Rx Max", 8),
		cz.Yellow("Tx Max", 8),
		cz.Yellow("Rx/Tx Errors", 12),
		cz.Yellow("Tot Rx Pkts", 12),
		cz.Yellow("Tot Tx Pkts", 12),
		cz.Yellow("Tot Rx Mbits", 12),
		cz.Yellow("Tot Tx Mbits", 12),
		cz.Yellow(" ", 6), // Extra field to allow scrolling horizontal
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
			cz.Wheat(comma(0) + "/" + comma(0)),
			cz.Cyan(comma(0)),
			cz.Cyan(comma(0)),
			cz.Red(comma(0) + "/" + comma(0)),
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

func (ps *PageSingleMode) displaySizes() {

	table := ps.singleSizes
	row := 0
	col := 0

	titles := []string{
		cz.Yellow("Port", 4),
		cz.Yellow("Broadcast", 12),
		cz.Yellow("Multicast", 12),
		cz.Yellow("Sizes 64", 12),
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

	ps.sizesOnce.Do(func() {
		ps.singleSizes.ScrollToBeginning()
	})
}

// Grab the load data and display the meters
func (ps *PageSingleMode) displayPerf() {

	view := ps.singlePerf

	str := ""

	for i := 0; i < pktgen.portCnt; i++ {
		str += ps.meter.Draw(ps.rxPercentRate[i], &meter.Info{
			Labels: []*meter.LabelInfo{
				{Val: fmt.Sprintf("%v", i), Fn: cz.Cyan},
				{Val: ": ", Fn: nil},
				{Val: "Rx ", Fn: cz.Yellow},
				{Val: fmt.Sprintf("%6.2f ", ps.rxPercentRate[i]), Fn: cz.DeepPink},
			},
			Bar: &meter.LabelInfo{Val: "", Fn: cz.MediumSpringGreen},
		})
		str += ps.meter.Draw(ps.txPercentRate[i], &meter.Info{
			Labels: []*meter.LabelInfo{
				{Val: "  ", Fn: nil},
				{Val: " ", Fn: nil},
				{Val: "Tx ", Fn: cz.Yellow},
				{Val: fmt.Sprintf("%6.2f ", ps.txPercentRate[i]), Fn: cz.DeepPink},
			},
			Bar: &meter.LabelInfo{Val: "", Fn: cz.Blue},
		})
	}
	str = str[:len(str)-1] // Strip the last newline character

	view.SetText(str)
	view.ScrollToBeginning()
}
