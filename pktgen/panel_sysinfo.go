// SPDX-License-Identifier: BSD-3-Clause
// Copyright(c) 2022 Intel Corporation

package main

import (
	"fmt"
	"strings"
	"sort"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"

	cz "github.com/KeithWiles/go-pktgen/pkgs/colorize"
	"github.com/KeithWiles/go-pktgen/pkgs/devbind"
	tab "github.com/KeithWiles/go-pktgen/pkgs/taborder"
	tlog "github.com/KeithWiles/go-pktgen/pkgs/ttylog"
)

// Display system device configuration

// tableData for each view
type tableData struct {
	name       string
	classes    []*devbind.DeviceConfig
	align      int
	fixedSize  int
	proportion int
	focus      bool
	key        rune
}

// tableInfo for each Table window
type tableInfo struct {
	changed bool
	length  int
	name    string
	view    *tview.Table
	classes []*devbind.DeviceConfig
	devlist []*devbind.DeviceClass
}

// PageSysInfo - Data for main page information
type PageSysInfo struct {
	topFlex *tview.Flex
	host    *tview.TextView
	mem     *tview.TextView
	hostNet *tview.Table
	devbind *devbind.BindInfo
	tables  []tableData
	tInfos  map[string]*tableInfo
}

const (
	sysinfoPanelName string = "System"
)

func init() {
	tlog.Register("SysInfoLogID")
}

// Printf - send message to the ttylog interface
func (ps *PageSysInfo) Printf(format string, a ...interface{}) {
	tlog.Log("SysInfoLogID", fmt.Sprintf("%T.", ps)+format, a...)
}

// setupSysInfo - setup and init the sysInfo page
func setupSysInfo() *PageSysInfo {

	ps := &PageSysInfo{}

	ps.devbind = devbind.New()

	db := ps.devbind

	ps.tInfos = make(map[string]*tableInfo)

	// Create the set of tables to display each section in a different window
	ps.tables = []tableData{
		{
			name:       "Network",
			classes:    db.Groups[devbind.NetworkGroup],
			align:      tview.AlignLeft,
			fixedSize:  0,
			proportion: 1,
			focus:      true,
			key:        'N',
		}, {
			name:       "Crypto",
			classes:    db.Groups[devbind.CryptoGroup],
			align:      tview.AlignLeft,
			fixedSize:  0,
			proportion: 1,
			focus:      true,
			key:        'c',
		}, {
			name:       "Compression",
			classes:    db.Groups[devbind.CompressGroup],
			align:      tview.AlignLeft,
			fixedSize:  0,
			proportion: 1,
			focus:      true,
			key:        'C',
		}, {
			name:       "DMA",
			classes:    db.Groups[devbind.DMAGroup],
			align:      tview.AlignLeft,
			fixedSize:  0,
			proportion: 1,
			focus:      true,
			key:        'd',
		},
	}

	// Add the table above into the tableInfo slice.
	for _, td := range ps.tables {
		ps.tInfos[td.name] = &tableInfo{classes: td.classes, name: td.name}
	}

	return ps
}

// SysInfoPanelSetup setup the main cpu page
func SysInfoPanelSetup(pages *tview.Pages, nextSlide func()) (pageName string, content tview.Primitive) {

	ps := setupSysInfo()

	to := tab.New(sysinfoPanelName, pktgen.app)

	flex0 := tview.NewFlex().SetDirection(tview.FlexRow)
	flex1 := tview.NewFlex().SetDirection(tview.FlexRow)
	flex2 := tview.NewFlex().SetDirection(tview.FlexColumn)

	TitleBox(flex0)

	ps.host = CreateTextView(flex2, "Host (h)", tview.AlignLeft, 0, 1, true)
	ps.mem = CreateTextView(flex2, "Memory (m)", tview.AlignLeft, 0, 1, false)
	flex1.AddItem(flex2, 0, 1, true)

	ps.hostNet = CreateTableView(flex1, "Host Network Stats (n)", tview.AlignLeft, 0, 1, false).
		SetSelectable(false, false).
		SetFixed(1, 1).
		SetSeparator(tview.Borders.Vertical)

	to.Add("host", ps.host, 'h')
	to.Add("memory", ps.mem, 'm')
	to.Add("hostName", ps.hostNet, 'n')

	ti := ps.tInfos

	// Create each table view for each of the device table entries
	for _, td := range ps.tables {
		s := fmt.Sprintf("%s Devices (%c)", td.name, td.key)

		ti[td.name].view = CreateTableView(flex1, s, td.align, td.fixedSize, td.proportion, td.focus).
			SetFixed(1, 0).
			SetSeparator(tview.Borders.Vertical)

		// Add the single key and define the tab order.
		to.Add(fmt.Sprintf("Table-%v", td.key), ti[td.name].view, td.key)
	}
	flex0.AddItem(flex1, 0, 3, true)

	to.SetInputDone()

	ps.topFlex = flex0

	modal := tview.NewModal().
		SetText("This is the Help Box: sysInfoHelp  Thank you for asking for help! Press Esc to close.").
		AddButtons([]string{"Got it"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.HidePage("sysInfoHelp")
		})
	AddModalPage("sysInfoHelp", modal)

	flex0.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Rune()
		switch k {
        case '?':
			pages.ShowPage("sysInfoHelp")
        case '-':
			pages.ShowPage("sysInfoWindow")
		}
		return event
	})

	// Setup static pages
	ps.displayHost(ps.host)
	ps.displayHostNet(ps.hostNet)
	ps.hostNet.ScrollToBeginning()

	pktgen.timers.Add(sysinfoPanelName, func(step int, ticks uint64) {
		if ps.topFlex.HasFocus() {
			pktgen.app.QueueUpdateDraw(func() {
				ps.displaySysInfo(step, ticks)
			})
		}
	})

	return sysinfoPanelName, ps.topFlex
}

// Callback timer routine to display the sysinfo panel
func (ps *PageSysInfo) displaySysInfo(step int, ticks uint64) {

	switch step {
	case 0:
		ps.displayMem(ps.mem)
		for _, t := range ps.tInfos {
			ps.collectData(t)
		}

	case 1:

	case 2:
		ps.displayHostNet(ps.hostNet)
		ps.displayPageSysInfo(step)
	}
}

// Display the Host information
func (ps *PageSysInfo) displayHost(view *tview.TextView) {

	str := ""
	info, _ := host.Info()
	str += fmt.Sprintf("Hostname: %s\n", cz.Yellow(info.Hostname))
	str += fmt.Sprintf("Host ID : %s\n", cz.Green(info.HostID))

	c := cases.Title(language.AmericanEnglish)
	str += fmt.Sprintf("OS      : %s-%s\n",
		cz.GoldenRod(c.String(info.OS)), cz.Orange(c.String(info.KernelVersion)))
	str += fmt.Sprintf("Platform: %s %s\nFamily  : %s\n",
		cz.MediumSpringGreen(c.String(info.Platform)),
		cz.LightSkyBlue(c.String(info.PlatformVersion)),
		cz.Green(c.String(info.PlatformFamily)))

	days := info.Uptime / (60 * 60 * 24)
	hours := (info.Uptime - (days * 60 * 60 * 24)) / (60 * 60)
	minutes := ((info.Uptime - (days * 60 * 60 * 24)) - (hours * 60 * 60)) / 60
	s := fmt.Sprintf("%d days, %d hours, %d minutes", days, hours, minutes)
	str += fmt.Sprintf("Uptime  : %s\n", cz.DeepPink(s))

	role := info.VirtualizationRole
	if len(role) == 0 {
		role = "unknown"
	}
	vsys := info.VirtualizationSystem
	if len(vsys) == 0 {
		vsys = "unknown"
	}
	str += fmt.Sprintf("Virtual Role: %s, System: %s", cz.Yellow(role), cz.Yellow(vsys))

	view.SetText(str)
}

// Display the information about the memory in the system
func (ps *PageSysInfo) displayMem(view *tview.TextView) {

	str := ""

	v, _ := mem.VirtualMemory()

	p := message.NewPrinter(language.English)
	str += fmt.Sprintf("Memory  Total: %s MiB\n", cz.Green(p.Sprintf("%d", v.Total/MegaBytes), 6))
	str += fmt.Sprintf("         Free: %s MiB\n", cz.Green(p.Sprintf("%d", v.Free/MegaBytes), 6))
	str += fmt.Sprintf("         Used: %s Percent\n\n", cz.Orange(v.UsedPercent, 6, 1))

	str += fmt.Sprintf("%s:\n", cz.MediumSpringGreen("Total Hugepage Info"))
	str += fmt.Sprintf("   Free/Total: %s/%s pages\n", cz.LightBlue(p.Sprintf("%d", v.HugePagesFree), 6),
		cz.LightBlue(p.Sprintf("%d", v.HugePagesTotal), 6))
	str += fmt.Sprintf("Hugepage Size: %s Kb", cz.LightBlue(p.Sprintf("%d", v.HugePageSize/KiloBytes), 6))

	view.SetText(str)
}

// Display the Host network information
func (ps *PageSysInfo) displayHostNet(view *tview.Table) {

	row := 0
	col := 0

	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("network interfaces: %s\n", err)
		return
	}

	titles := []string{
		cz.Yellow("Name"),
		cz.Yellow("IP Address"),
		cz.Yellow("MTU"),
		cz.Yellow("RX Pkts"),
		cz.Yellow("TX Pkts"),
		cz.Yellow("RX Err"),
		cz.Yellow("TX Err"),
		cz.Yellow("RX Drop"),
		cz.Yellow("Tx Drop"),
		cz.Yellow("Flags"),
		cz.Yellow("MAC"),
		cz.Yellow(" ", 20),
	}
	row = TableSetHeaders(view, 0, 0, titles)

	setCell := func(row, col int, value string, left bool) int {
		align := tview.AlignRight
		if left {
			align = tview.AlignLeft
		}
		tableCell := tview.NewTableCell(value).
			SetAlign(align).
			SetSelectable(false)
		ps.hostNet.SetCell(row, col, tableCell)
		col++

		return col
	}

	ioCount, err := net.IOCounters(true)
	if err != nil {
		ps.Printf("network IO Count: %s\n", err)
		return
	}

	p := message.NewPrinter(language.English)

	for _, f := range ifaces {
		if f.Name == "lo" {
			continue
		}

		col = setCell(row, 0, cz.LightBlue(f.Name), true)
		if len(f.Addrs) > 0 {
			col = setCell(row, col, cz.Orange(f.Addrs[0].Addr), false)
		} else {
			col = setCell(row, col, " ", false)
		}
		col = setCell(row, col, cz.MediumSpringGreen(f.MTU), false)

		for _, k := range ioCount {
			if k.Name != f.Name {
				continue
			}
			rowData := []string{
				cz.Wheat(p.Sprintf("%d", k.PacketsRecv)),
				cz.Wheat(p.Sprintf("%d", k.PacketsSent)),
				cz.Red(p.Sprintf("%d", k.Errin)),
				cz.Red(p.Sprintf("%d", k.Errout)),
				cz.Red(p.Sprintf("%d", k.Dropin)),
				cz.Red(p.Sprintf("%d", k.Dropout)),
			}
			for _, v := range rowData {
				col = TableCellSet(ps.hostNet, row, col, v)
			}
			break
		}
		col = setCell(row, col, cz.LightSkyBlue(f.Flags), false)
		setCell(row, col, cz.Cyan(f.HardwareAddr), false)

		row++
	}

	for ; row < view.GetRowCount(); row++ {
		view.RemoveRow(row)
	}
}

// Display the given devbind data panel for each window
func (ps *PageSysInfo) displayPageSysInfo(step int) {
	for _, ti := range ps.tInfos {
		if ti.changed {
			ti.changed = false
			ps.displayView(ti)
		}
	}
}

// Collect the data to be displayed in the different device windows
func (ps *PageSysInfo) collectData(ti *tableInfo) {

	deviceList := make([]*devbind.DeviceClass, 0)

	tlog.DebugPrintf("Name: %s, Classes: %+v\n", ti.name, ti.classes)

	// Convert the map into a slice to be able to sort it
	for _, l := range ps.devbind.FindDevicesByDeviceClass(ti.name, ti.classes) {
		deviceList = append(deviceList, l)
	}
	tlog.DebugPrintf("Name: %s, deviceList: %+v\n", ti.name, deviceList)

	sort.Slice(deviceList, func(i, j int) bool {
		return deviceList[j].Slot > deviceList[i].Slot
	})

	// Set the device list and set the changed flag to force update of window
	ti.devlist = deviceList
	if ti.length != len(deviceList) {
		ti.changed = true
		ti.length = len(deviceList)
	}
}

// Display the formation into the given table, all windows use this routine
func (ps *PageSysInfo) displayView(ti *tableInfo) {

	view := ti.view

	titles := []string{
		cz.CornSilk("Slot"),
		cz.CornSilk("Vendor ID"),
		cz.CornSilk("Vendor Name"),
		cz.CornSilk("Device Description"),
		cz.CornSilk("Interface"),
		cz.CornSilk("Driver"),
		cz.CornSilk("Active"),
		cz.CornSilk("NUMA"),
	}
	row := TableSetHeaders(view, 0, 0, titles)

	for _, d := range ti.devlist {
		col := 0

		SetCell(view, row, col, cz.DeepPink(d.Slot), tview.AlignLeft)
		col++

		s := fmt.Sprintf("%s:%s", cz.SkyBlue(d.Vendor.ID), cz.SkyBlue(d.Device.ID))
		SetCell(view, row, col, s, tview.AlignLeft)
		col++

		str := d.Vendor.Str
		idx := strings.Index(str, "[")
		if idx != -1 {
			str = str[:idx-1]
		}
		SetCell(view, row, col, cz.SkyBlue(str), tview.AlignLeft)
		col++

		str = d.SDevice.Str
		idx = strings.Index(str, "[")
		if idx != -1 {
			str = str[:idx-1]
		}
		SetCell(view, row, col, cz.LightGreen(str), tview.AlignLeft)
		col++

		str = d.Interface
		SetCell(view, row, col, cz.ColorWithName("Tomato", str), tview.AlignLeft)
		col++

		str = d.Driver
		SetCell(view, row, col, cz.LightYellow(str), tview.AlignLeft)
		col++

		str = ""
		if d.Active {
			str = cz.Orange("*Active*")
		}
		SetCell(view, row, col, str, tview.AlignLeft)
		col++

		str = d.NumaNode
		idx = strings.Index(str, "[")
		if idx != -1 {
			str = str[:idx-1]
		}
		SetCell(view, row, col, cz.MistyRose(str), tview.AlignLeft)
		col++

		row++
	}

	ti.view.ScrollToBeginning()
}
