package termui

import "fmt"

type Input struct {
	Block
	Text        string
	TextFgColor Attribute
	TextBgColor Attribute
	WrapLength  int // words wrap limit. Note it may not work properly with multi-width char
	CursorPos   int
	OnEnter     func(str string)
}

// NewPar returns a new *Par with given text as its content.
func NewInput(s string) *Input {
	return &Input{
		Block:       *NewBlock(),
		Text:        s,
		TextFgColor: ThemeAttr("par.text.fg"),
		TextBgColor: ThemeAttr("par.text.bg"),
		WrapLength:  0,
		CursorPos:   0,
	}
}

func InputKbdHandler(e Event, w Widget) {
	if ActiveWgtId!=w.Id() {
		return
	}
	event := e.Data.(EvtKbd)
	input := DefaultWgtMgr[w.Id()].Data.(*Input)
	//char  := event.KeyStr
	//if char.siI

			//ks := []string{"<insert>", "<delete>", "<home>", "<end>", "<previous>", "<next>", "<up>", "<down>", "<left>", "<right>"}
	input.Text += event.KeyStr
	w.BlockRef().BorderLabel=fmt.Sprintf("%s", event.KeyStr)
	Render(w)
}

// Buffer implements Bufferer interface.
func (p *Input) Buffer() Buffer {
	buf := p.Block.Buffer()

	fg, bg := p.TextFgColor, p.TextBgColor

	var text string

	if ActiveWgtId==p.Id() {
		text = p.Text + "█"
	} else {
		text = p.Text
	}

	cs := DefaultTxBuilder.Build(text, fg, bg)

	// wrap if WrapLength set
	if p.WrapLength < 0 {
		cs = wrapTx(cs, p.Width-2)
	} else if p.WrapLength > 0 {
		cs = wrapTx(cs, p.WrapLength)
	}

	y, x, n := 0, 0, 0
	for y < p.innerArea.Dy() && n < len(cs) {
		w := cs[n].Width()
		if cs[n].Ch == '\n' || x+w > p.innerArea.Dx() {
			y++
			x = 0 // set x = 0
			if cs[n].Ch == '\n' {
				n++
			}

			if y >= p.innerArea.Dy() {
				buf.Set(p.innerArea.Min.X+p.innerArea.Dx()-1,
					p.innerArea.Min.Y+p.innerArea.Dy()-1,
					Cell{Ch: '…', Fg: p.TextFgColor, Bg: p.TextBgColor})
				break
			}
			continue
		}

		buf.Set(p.innerArea.Min.X+x, p.innerArea.Min.Y+y, cs[n])

		n++
		x += w
	}

	return buf
}

