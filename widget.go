// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"fmt"
	"sync"
)

// event mixins
type WgtMgr map[string]WgtInfo

var ActiveWgtId string

type WgtInfo struct {
	Handlers map[string]func(Event, Widget)
	WgtRef   Widget
	BlockRef *Block
	Id       string
	Data interface{}
}

type Widget interface {
	Id() string
	BlockRef() *Block
	Bufferer
}

func (w WgtInfo) IncludePoint(X int, Y int) bool {
	b := w.BlockRef
	return b.X<X && b.X+b.Width>X && b.Y<Y && b.Y+b.Height>Y
}

func NewWgtInfo(wgt Widget) WgtInfo {
	return WgtInfo{
		Handlers: make(map[string]func(Event, Widget)),
		WgtRef:   wgt,
		BlockRef: wgt.BlockRef(),
		Id:       wgt.Id(),
		Data:     wgt,
	}
}

func NewWgtMgr() WgtMgr {
	wm := WgtMgr(make(map[string]WgtInfo))
	return wm

}

func (wm WgtMgr) AddWgt(wgts ...Widget) {
	for _, wgt := range wgts {
		wm[wgt.Id()] = NewWgtInfo(wgt)
	}
}

func (wm WgtMgr) RmWgt(wgts ...Widget) {
	for _, wgt := range wgts {
		wm.RmWgtById(wgt.Id())
	}
}

func (wm WgtMgr) RmWgtById(id string) {
	delete(wm, id)
}

func (wm WgtMgr) AddWgtHandler(id, path string, h func(Event, Widget)) {
	if w, ok := wm[id]; ok {
		w.Handlers[path] = h
	}
}

func (wm WgtMgr) RmWgtHandler(id, path string) {
	if w, ok := wm[id]; ok {
		delete(w.Handlers, path)
	}
}

var counter struct {
	sync.RWMutex
	count int
}

func GenId() string {
	counter.Lock()
	defer counter.Unlock()

	counter.count += 1
	return fmt.Sprintf("%d", counter.count)
}

func (wm WgtMgr) WgtHandlersHook() func(Event) {
	return func(e Event) {
		for _, v := range wm {
			if k := findMatchEx(v.Handlers, e.Path); k != "" {
				if e.Path=="/sys/mouse" {
					m_e := e.Data.(EvtMouse)
					if v.IncludePoint(m_e.X, m_e.Y) {
						v.Handlers[k](e, v.WgtRef)
					}
				} else if e.Path[0:8]=="/sys/kbd" {
					if v.Id==ActiveWgtId {
						v.Handlers[k](e, v.WgtRef)
					}
				} else {
					v.Handlers[k](e, v.WgtRef)
				}
			}
		}
	}
}

var DefaultWgtMgr WgtMgr

//func (b *Block) Handle(path string, handler func(Event)) {
//	if _, ok := DefaultWgtMgr[b.Id()]; !ok {
//		DefaultWgtMgr.AddWgt(b)
//	}
//
//	DefaultWgtMgr.AddWgtHandler(b.Id(), path, handler)
//}
