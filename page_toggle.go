package main

import (
	"github.com/mattn/go-runewidth"
	"github.com/thomgray/egg"
)

type PageToggleViewController struct {
	*egg.View
	Page      int
	PageNames []string
}

func MakePageToggleViewController(pages []string) *PageToggleViewController {
	pt := PageToggleViewController{}
	pt.View = egg.MakeView()
	pt.PageNames = pages
	pt.Page = 0

	pt.OnDraw(func(c egg.Canvas) {
		x := 0
		c.DrawString2("< |", x, 0)
		x += 3
		for i, p := range pt.PageNames {
			fg := c.Foreground
			atts := c.Attribute

			if i == pt.Page {
				fg = egg.ColorCyan
				atts = egg.AttrBold | egg.AttrUnderline
			}

			c.DrawString(" "+p+" ", x, 0, fg, c.Background, atts)
			x += runewidth.StringWidth(p) + 2
			c.DrawString2("|", x, 0)
			x++
		}
		c.DrawString2(" >", x, 0)
	})
	return &pt
}

func (pt *PageToggleViewController) SetPage(page int) {
	pt.Page = page
}

func (pt *PageToggleViewController) GetView() *egg.View {
	return pt.View
}
