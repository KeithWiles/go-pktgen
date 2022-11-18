// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2019-2022 Intel Corporation

package taborder

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	defaultBorderColor   tcell.Color = tcell.ColorGreen
	highlightBorderColor tcell.Color = tcell.ColorBlue
)

// TabInfo for windows on the current panel
type TabInfo struct {
	Index     int
	View      interface{}
	EKey      *tcell.EventKey
	Name      string
}

// Tab for all windows in a panel
type Tab struct {
	Name          string
	TabList       []*TabInfo
	CurrentIndex  int
	PreviousIndex int
	Application   *tview.Application
}

// New information object
func New(name string, application *tview.Application) *Tab {
	return &Tab{Name: name, Application: application}
}

// Add to the given list of windows
func (to *Tab) Add(name string, w interface{}, key interface{}) (*TabInfo, error) {
	if to == nil {
		return nil, fmt.Errorf("invalid tabOrder pointer")
	}

	tabInfo := &TabInfo{View: w, Name: name}

	if key != nil {
		switch k := key.(type) {
		case tcell.Key:
			tabInfo.EKey = tcell.NewEventKey(k, 0, tcell.ModNone)
		case rune:
			tabInfo.EKey = tcell.NewEventKey(tcell.KeyRune, k, tcell.ModNone)
		}
	}

	tabInfo.Index = len(to.TabList)

	to.TabList = append(to.TabList, tabInfo)

	if len(to.TabList) == 1 {
		to.colorBorder(tabInfo.View, highlightBorderColor)
	} else {
		to.colorBorder(tabInfo.View, defaultBorderColor)
	}

	return tabInfo, nil
}

// SetDefaultBorderColor to the normal non-selected border color
func (to *Tab) SetDefaultBorderColor(color tcell.Color) {
	defaultBorderColor = color
}

// SetHighlightBorderColor to the normal non-selected border color
func (to *Tab) SetHighlightBorderColor(color tcell.Color) {
	highlightBorderColor = color
}

// setFocus to the tview primitive
func (to *Tab) setFocus(a interface{}) {

	switch t := a.(type) {
	case *tview.TextView:
		to.Application.SetFocus(t)
	case *tview.Table:
		to.Application.SetFocus(t)
	case *tview.Form:
		to.Application.SetFocus(t)
	}
}

// colorBorder to the tview
func (to *Tab) colorBorder(a interface{}, color tcell.Color) {

	switch t := a.(type) {
	case *tview.TextView:
		t.Box.SetBorderColor(color)
	case *tview.Table:
		t.Box.SetBorderColor(color)
	case *tview.Form:
		t.Box.SetBorderColor(color)
	}
}

func (to *Tab) findKey(ek *tcell.EventKey) *TabInfo {

	for _, tab := range to.TabList {
		if tab.EKey.Name() == ek.Name() {
			return tab
		}
	}
	return nil
}

// inputCapture for taborder
func (to *Tab) inputCapture(ek *tcell.EventKey) *tcell.EventKey {

	if ek.Key() != tcell.KeyBacktab && ek.Key() != tcell.KeyTab {
		if tab := to.findKey(ek); tab != nil {
			to.colorBorder(to.TabList[to.CurrentIndex].View, defaultBorderColor)
			to.setFocus(tab.View)
			to.colorBorder(tab.View, highlightBorderColor)
			to.PreviousIndex, to.CurrentIndex = to.CurrentIndex, tab.Index
		}
	}
	return ek
}

// SetInputFocus sets the focus to the given event
func (to *Tab) SetInputFocus(key interface{}) {

	var eKey *tcell.EventKey

	switch k := key.(type) {
	case tcell.Key:
		if k == 0 {
			return
		}
		eKey = tcell.NewEventKey(k, 0, tcell.ModNone)
	case rune:
		eKey = tcell.NewEventKey(tcell.KeyRune, k, tcell.ModNone)
	}

	if tab := to.findKey(eKey); tab != nil {
		to.colorBorder(to.TabList[to.CurrentIndex].View, defaultBorderColor)
		to.setFocus(tab.View)
		to.colorBorder(tab.View, highlightBorderColor)
		to.PreviousIndex, to.CurrentIndex = to.CurrentIndex, tab.Index
	}
}

func (to *Tab) SetCurrentInputFocus() {
	to.SetInputFocus(to.TabList[to.CurrentIndex].EKey)
}

// doDone key handling for Tab and Backtab
func (to *Tab) doDone(key tcell.Key) {

	p := to.TabList[to.CurrentIndex]
	to.colorBorder(p.View, defaultBorderColor)

	if key == tcell.KeyBacktab {
		if to.CurrentIndex == 0 {
			p = to.TabList[len(to.TabList)-1]
		} else {
			p = to.TabList[to.CurrentIndex-1]
		}
	} else if key == tcell.KeyTab {
		if to.CurrentIndex < (len(to.TabList) - 1) {
			p = to.TabList[to.CurrentIndex+1]
		} else {
			p = to.TabList[0]
		}
	}

	to.setFocus(p.View)
	to.colorBorder(p.View, highlightBorderColor)

	to.PreviousIndex, to.CurrentIndex = to.CurrentIndex, p.Index
}

// setInput for tview
func (to *Tab) setInput(a interface{}, inputFunc func(ev *tcell.EventKey) *tcell.EventKey) {

	switch t := a.(type) {
	case *tview.TextView:
		t.SetInputCapture(inputFunc)
	case *tview.Table:
		t.SetInputCapture(inputFunc)
	case *tview.Form:
		t.SetInputCapture(inputFunc)
	}
}

// setDone function for tabs
func (to *Tab) setDone(a interface{}, doneFunc func(key tcell.Key)) {

	switch t := a.(type) {
	case *tview.TextView:
		t.SetDoneFunc(doneFunc)
	case *tview.Table:
		t.SetDoneFunc(doneFunc)
	case *tview.Form:
	}
}

// SetInputDone functions and data
func (to *Tab) SetInputDone() error {
	if to.TabList == nil {
		return fmt.Errorf("tab list is nil")
	}

	for _, tab := range to.TabList {
		to.setInput(tab.View, func(ek *tcell.EventKey) *tcell.EventKey {
			return to.inputCapture(ek)
		})
		to.setDone(tab.View, func(key tcell.Key) {
			to.doDone(key)
		})
	}

	return nil
}
