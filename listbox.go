// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.
//
// Portions copyright 2017 Patrick Devine <patrick@immense.ly>

package termui

// List displays []Item as its items (items are pairs of text and values),
// it has a Overflow option (default is "hidden"), when set to "hidden",
// the item exceeding List's width is truncated, but when set to "wrap",
// the overflowed text breaks into next line.

type Item struct {
	ItemVal string
	Text    string
}

type ListBox struct {
	Block
	Items       []Item
	ItemFgColor Attribute
	ItemBgColor Attribute
	Selected    int
	lowerBound  int
}

// NewList returns a new *List with current theme.
func NewListBox() *ListBox {
	l := &ListBox{Block: *NewBlock()}
	l.ItemFgColor = ThemeAttr("list.item.fg")
	l.ItemBgColor = ThemeAttr("list.item.bg")
	l.Selected = 0
	l.lowerBound = 0
	return l
}

func (l *ListBox) ResetSelection() {
	l.Selected = 0
	l.lowerBound = 0
}

// Buffer implements Bufferer interface.
func (l *ListBox) Buffer() Buffer {
	buf := l.Block.Buffer()

	trimItems := l.GetItemsStrs()
	totalItems := len(l.GetItemsStrs())
	if len(trimItems) > l.innerArea.Dy() {
		trimItems = trimItems[l.lowerBound : l.innerArea.Dy()+l.lowerBound]
	}
	for i, v := range trimItems {
		var cs []Cell
		if i+l.lowerBound == l.Selected {
			cs = DTrimTxCls(DefaultTxBuilder.Build(v, l.ItemFgColor, AttrReverse), l.innerArea.Dx())
		} else {
			cs = DTrimTxCls(DefaultTxBuilder.Build(v, l.ItemFgColor, l.ItemBgColor), l.innerArea.Dx())
		}
		j := 0
		for _, vv := range cs {
			w := vv.Width()
			buf.Set(l.innerArea.Min.X+j, l.innerArea.Min.Y+i, vv)
			j += w
		}
	}
	// display scroll arrows
	if l.lowerBound > 0 {
		buf.Set(l.innerArea.Dx(), 1, Cell{Ch: '^'})
	}
	if totalItems > l.lowerBound+l.innerArea.Dy() {
		buf.Set(l.innerArea.Dx(), l.innerArea.Dy(), Cell{Ch: 'v'})
	}
	return buf
}

func (l *ListBox) GetItemsStrs() []string {
	var strs []string
	for _, item := range l.Items {
		strs = append(strs, item.Text)
	}
	return strs
}

func (l *ListBox) SelectExtY(y int) {
	y = y - l.Y + l.lowerBound - 1
	if y >= len(l.Items) {
		y = len(l.Items) - 1
	}
	l.Selected = y
}

func (l *ListBox) CentreSelection() {
	if len(l.Items) < l.innerArea.Dy() {
		return
	}
	l.lowerBound = l.Selected - l.innerArea.Dy()/2
	if l.lowerBound+l.innerArea.Dy() >= len(l.Items) {
		l.lowerBound = len(l.Items) - l.innerArea.Dy() - 1
	}
	if l.lowerBound < 0 {
		l.lowerBound = 0
	}
}

func (l *ListBox) Up() {
	if l.Selected > 0 {
		l.Selected -= 1
		if l.Selected < l.lowerBound {
			l.lowerBound -= 1
		}
	}
}

func (l *ListBox) Down() {
	if l.Selected < len(l.Items)-1 {
		l.Selected += 1
		if l.Selected >= l.innerArea.Dy()+l.lowerBound {
			l.lowerBound += 1
		}
	}
}

func (l *ListBox) PageUp() {
	for n := 0; n < l.innerArea.Dy(); n++ {
		l.Up()
	}
}

func (l *ListBox) PageDown() {
	for n := 0; n < l.innerArea.Dy(); n++ {
		l.Down()
	}
}

func (l *ListBox) Current() Item {
	return l.Items[l.Selected]
}
