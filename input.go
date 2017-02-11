package termui

type Input struct {
	Block
	Text        []rune
	TextFgColor Attribute
	TextBgColor Attribute
	Hi          Attribute
	Password    bool
	WrapLength  int // words wrap limit. Note it may not work properly with multi-width char
	CursorPos   int
	OnEnter     func(str string)
}

func (p *Input) Input() string {
	return string(p.Text)
}

// NewPar returns a new *Par with given text as its content.
func NewInput(s string) *Input {
	a := []rune(s)
	return &Input{
		Block:       *NewBlock(),
		Text:        a,
		TextFgColor: ThemeAttr("par.text.fg"),
		Hi:          ColorYellow,
		TextBgColor: ThemeAttr("par.text.bg"),
		WrapLength:  0,
		CursorPos:   len(a),
	}
}

func InputMouseHandler(e Event, w Widget) {
	t := e.Data.(EvtMouse)
	x, _ := t.X, t.Y
	input := DefaultWgtMgr[w.Id()].Data.(*Input)
	input.CursorPos = x - input.X - 1
	if input.CursorPos<0 {
		input.CursorPos=0
	}
	if input.CursorPos>len(input.Text) {
		input.CursorPos=len(input.Text)
	}
	Render(w)
}

func InputKbdHandler(e Event, w Widget) {
	if ActiveWgtId!=w.Id() {
		return
	}
	event := e.Data.(EvtKbd)
	input := DefaultWgtMgr[w.Id()].Data.(*Input)
	char  := event.KeyStr
	if len(char)>1 && char!="<space>" {
		switch(char) {
		case "<left>":
			input.CursorPos--
		case "<right>":
			input.CursorPos++
		case "C-8":
			i := input.CursorPos-1
			if input.CursorPos!=0 && len(input.Text)>0 && i<=len(input.Text) {
				input.Text = append(input.Text[:i], input.Text[i+1:]...)
				input.CursorPos--
			}
		case "<enter>":
			if input.OnEnter!=nil {
				input.OnEnter(input.Input())
			}
		case "<tab>":
			if input.OnEnter!=nil {
				input.OnEnter(input.Input())
			}
		case "<delete>":
			i := input.CursorPos
			if len(input.Text)>0 && i<len(input.Text) {
				input.Text = append(input.Text[:i], input.Text[i+1:]...)
			}
		}
	} else {
		if char=="<space>" {
			char = " "
		}
		var res []rune
		a := []rune(char)
		if len(input.Text)==input.CursorPos {
			input.Text = append(input.Text, a[0])
			input.CursorPos++
		} else {
			for i, r := range input.Text {
				if i==input.CursorPos {
					res = append(res, a[0])
				}
				res = append(res, r)
			}
			input.Text = res
			input.CursorPos++
		}
	}
	if input.CursorPos<0 {
		input.CursorPos=0
	}
	if input.CursorPos>len(input.Text) {
		input.CursorPos=len(input.Text)
	}

	//ks := []string{"<insert>", "<delete>", "<home>", "<end>", "<previous>", "<next>", "<up>", "<down>", "<left>", "<right>"}
	//w.BlockRef().BorderLabel=fmt.Sprintf("%s", event.KeyStr)
	Render(w)
}

// Buffer implements Bufferer interface.
func (p *Input) Buffer() Buffer {
	buf := p.Block.Buffer()

	fg, bg := p.TextFgColor, p.TextBgColor

	text := string(p.Text)

	if p.Password {
		text = ""
		for i:=0; i<len(p.Text); i++ {
			text += "*"
		}
	}

    active := ActiveWgtId==p.Id()
	if active && p.CursorPos>=len(p.Text) {
		text += " "
	}

	cs := DefaultTxBuilder.Build(text, fg, bg)
	hi := DefaultTxBuilder.Build(text, fg, p.Hi)

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

		if active && x==p.CursorPos {
			buf.Set(p.innerArea.Min.X+x, p.innerArea.Min.Y+y, hi[n])
		} else {
			buf.Set(p.innerArea.Min.X+x, p.innerArea.Min.Y+y, cs[n])
		}

		n++
		x += w
	}

	return buf
}

