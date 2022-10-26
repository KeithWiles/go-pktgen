// SPDX-License-Identifier: BSD-3-Clause
// Copyright(c) 2022 Intel Corporation

package main

import (
	"fmt"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"

	cz "github.com/KeithWiles/go-pktgen/pkgs/colorize"
	tab "github.com/KeithWiles/go-pktgen/pkgs/taborder"
	tlog "github.com/KeithWiles/go-pktgen/pkgs/ttylog"
)

// PageSysInfo - Data for main page information
type PageSysInfo struct {
	topFlex *tview.Flex
	host    *tview.TextView
	mem     *tview.TextView
	hostNet *tview.Table
}

const (
	sysinfoPanelName string = "System"
)

func init() {
	tlog.Register("SysInfoLogID")
}

// Printf - send message to the ttylog interface
func (pg *PageSysInfo) Printf(format string, a ...interface{}) {
	tlog.Log("SysInfoLogID", fmt.Sprintf("%T.", pg)+format, a...)
}

// setupSysInfo - setup and init the sysInfo page
func setupSysInfo() *PageSysInfo {

	pg := &PageSysInfo{}

	return pg
}

// SysInfoPanelSetup setup the main cpu page
func SysInfoPanelSetup(nextSlide func()) (pageName string, content tview.Primitive) {

	pg := setupSysInfo()

	to := tab.New(sysinfoPanelName, pktgen.app)

	flex0 := tview.NewFlex().SetDirection(tview.FlexRow)
	flex1 := tview.NewFlex().SetDirection(tview.FlexRow)
	flex2 := tview.NewFlex().SetDirection(tview.FlexColumn)

	TitleBox(flex0)

	pg.host = CreateTextView(flex2, "Host (h)", tview.AlignLeft, 0, 1, true)
	pg.mem = CreateTextView(flex2, "Memory (m)", tview.AlignLeft, 0, 1, false)
	flex1.AddItem(flex2, 10, 1, true)

	pg.hostNet = CreateTableView(flex1, "Host Network Stats (n)", tview.AlignLeft, 0, 1, false).
		SetSelectable(false, false).
		SetFixed(1, 1).
		SetSeparator(tview.Borders.Vertical)
	flex0.AddItem(flex1, 0, 3, true)

	to.Add("host", pg.host, 'h')
	to.Add("memory", pg.mem, 'm')
	to.Add("hostName", pg.hostNet, 'n')

	to.SetInputDone()

	pg.topFlex = flex0

	// Setup static pages
	pg.displayHost(pg.host)
	pg.displayHostNet(pg.hostNet)
	pg.hostNet.ScrollToBeginning()

	pktgen.timers.Add(sysinfoPanelName, func(step int, ticks uint64) {
		if pg.topFlex.HasFocus() {
			pktgen.app.QueueUpdateDraw(func() {
				pg.displaySysInfo(step, ticks)
			})
		}
	})

	return sysinfoPanelName, pg.topFlex
}

// Callback timer routine to display the sysinfo panel
func (pg *PageSysInfo) displaySysInfo(step int, ticks uint64) {

	switch step {
	case 0:
		pg.displayMem(pg.mem)

	case 1:

	case 2:
		pg.displayHostNet(pg.hostNet)
	}
}

// Display the Host information
func (pg *PageSysInfo) displayHost(view *tview.TextView) {

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
func (pg *PageSysInfo) displayMem(view *tview.TextView) {

	str := ""

	v, _ := mem.VirtualMemory()

	p := message.NewPrinter(language.English)
	str += fmt.Sprintf("Memory  Total: %s MiB\n", cz.Green(p.Sprintf("%d", v.Total/MegaBytes), 6))
	str += fmt.Sprintf("         Free: %s MiB\n", cz.Green(p.Sprintf("%d", v.Free/MegaBytes), 6))
	str += fmt.Sprintf("         Used: %s Percent\n\n", cz.Orange(v.UsedPercent, 6, 1))

	str += fmt.Sprintf("%s:\n", cz.MediumSpringGreen("Total Hugepage Info"))
	str += fmt.Sprintf("   Free/Total: %s/%s pages\n", cz.LightBlue(p.Sprintf("%d", v.HugePagesFree), 6),
		cz.LightBlue(p.Sprintf("%d", v.HugePagesTotal), 6))
	str += fmt.Sprintf("Hugepage Size: %s Kb\n", cz.LightBlue(p.Sprintf("%d", v.HugePageSize/KiloBytes), 6))

	view.SetText(str)
}

// Display the Host network information
func (pg *PageSysInfo) displayHostNet(view *tview.Table) {

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
		cz.Yellow(" ", 10),
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
		pg.hostNet.SetCell(row, col, tableCell)
		col++

		return col
	}

	ioCount, err := net.IOCounters(true)
	if err != nil {
		pg.Printf("network IO Count: %s\n", err)
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
				col = TableCellSet(pg.hostNet, row, col, v)
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
