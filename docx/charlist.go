package docx

import (
	"errors"
	"io"
	"strings"

	"github.com/saman3d/samdoc"
	"github.com/saman3d/samdoc/xml"
)

var (
	StartPlace = "{{"
	EndPlace   = "}}"
)

type Char struct {
	Rune rune
	T    *xml.UniversalElement
	R    *xml.UniversalElement
	P    *xml.UniversalElement
}

type CharNode struct {
	Char *Char
	Prev *CharNode
	Next *CharNode
}

type CharList struct {
	Head  *CharNode
	Tail  *CharNode
	First *CharNode
}

type ReplacerFunc func(placeholder string) (string, bool)

func NewReplacerFunc(model interface{}) (ReplacerFunc, error) {
	strct, err := samdoc.NewStructure(model)
	if err != nil {
		return nil, err
	}
	return func(f string) (string, bool) {
		spl := strings.Split(strings.TrimSpace(f), ".")

		val, err := strct.Get(samdoc.FieldQuery(spl))
		if err != nil {
			return "", false
		}

		return val, true
	}, nil
}

func (l *CharList) LoadFromElement(con *xml.UniversalElement) {
	nump := 0
	for _, p := range con.Children {
		if p.XMLName == "w:p" {
			for _, r := range p.Children {
				if r.XMLName == "w:r" {
					nump++
					t := r.GetElementByName("w:t")
					if t == nil {
						char := &Char{
							Rune: 0,
							T:    nil,
							R:    r,
							P:    p,
						}
						l.Insert(char)
					} else {
						for _, ch := range t.Data {
							char := &Char{
								Rune: ch,
								T:    t,
								R:    r,
								P:    p,
							}
							l.Insert(char)
						}
					}

				}
			}

		} else {
			break
		}
	}
}

func (l *CharList) Insert(r *Char) {
	node := &CharNode{Char: r, Prev: nil, Next: nil}
	if l.Head == nil {
		l.Head = node
		l.Tail = node
		l.First = node
	} else {
		p := l.Head
		for p.Next != nil {
			p = p.Next
		}
		node.Prev = p
		p.Next = node
		l.Tail = node
	}
}

func (l *CharList) InsertBefore(r *Char) {
	node := &CharNode{Char: r, Prev: nil, Next: nil}
	if l.Head == nil {
		l.Head = node
		l.Tail = node
		l.First = node
	} else {
		node.Next = l.Head
		node.Prev = l.Head.Prev
		if l.Head.Prev != nil {
			l.Head.Prev.Next = node
			l.Head.Prev = node
		}
	}
}

func (l *CharList) GoToFirst() {
	l.Head = l.First
}

func (l *CharList) Len() int {
	c := 0
	if l.Head == nil {
		return c
	}

	h := l.Tail
	for h != nil {
		c++
		h = h.Prev
	}
	return c
}

func (l *CharList) Current() *Char {
	return l.Head.Char
}

func (l *CharList) Next() *Char {
	if l.Head.Next == nil {
		return nil
	}
	l.Head = l.Head.Next
	return l.Head.Char
}

func (l *CharList) Prev() *Char {
	if l.Head == nil || l.Head.Prev == nil {
		return nil
	}
	l.Head = l.Head.Prev
	return l.Head.Char
}

func (l *CharList) ToParagraphList() []*xml.UniversalElement {

	var ps = make([]*xml.UniversalElement, 0)
	l.GoToFirst()

	p := l.Current().P

	hp := NewParagraph(l.Current().P, l.Current().R, l.Current().T)
	hp.Insert(l.Current())
	for c := l.Next(); c != nil; c = l.Next() {
		if p != c.P {
			p = c.P
			ps = append(ps, hp.ToUniversal())
			hp = NewParagraph(l.Current().P, l.Current().R, l.Current().T)
		}
		hp.Insert(c)
	}
	ps = append(ps, hp.ToUniversal())

	return ps
}

func (l *CharList) LookAhead(t int) string {
	b := ""
	i := 0
	for c := l.Head.Next; c != nil && i < t; c = c.Next {
		b += string(c.Char.Rune)
		i++
	}
	return b
}

func (l *CharList) LookAheadTill(delim string) (string, error) {
	b := ""
	for c := l.Head.Next; c != nil; c = c.Next {
		if chs := c; chs.Char.Rune == rune(delim[0]) {
			for i, ch := range delim[1:] {
				chs = chs.Next
				if chs.Char.Rune != rune(ch) {
					break
				} else if i == len(delim)-2 {
					return b, nil
				}
			}
		}
		b += string(c.Char.Rune)
	}
	return "", errors.New("didn't find a match")
}

func (l *CharList) String() string {
	b := ""
	l.GoToFirst()
	b += string(l.Current().Rune)
	for c := l.Next(); c != nil; c = l.Next() {
		b += string(c.Rune)
	}
	return b
}

func (l *CharList) Seek(t int) {
	if t < 0 {
		for i := 0; i > t && l.Head.Prev != nil; i-- {
			l.Prev()
		}
	} else {
		for i := 0; i < t && l.Head.Next != nil; i++ {
			l.Next()
		}
	}
}

func (l *CharList) SeekTill(delim string) error {
	for ; l.Head.Next != nil; l.Next() {
		if l.LookAhead(len(delim)) == delim {
			return nil
		}
	}
	return io.EOF
}

func (l *CharList) Remove(t int) error {
	var cur = l.Head
	for i := 0; i <= t; i++ {
		c := l.Next()
		if c == nil {
			return io.EOF
		}
	}
	l.Head.Prev = cur
	cur.Next = l.Head
	l.Prev()
	return nil
}

func (l *CharList) Replace(rf ReplacerFunc) error {
	l.GoToFirst()
	var err error
	var res string
	for err = l.SeekTill(StartPlace); err == nil; l.SeekTill(StartPlace) {
		l.Seek(2)
		res, err = l.LookAheadTill(EndPlace)
		if err != nil {
			continue
		}
		repstring, ok := rf(res)
		if !ok {
			l.Seek(2 + len(res))
			continue
		}
		char := l.Next()
		l.Seek(-3)
		l.Remove(4 + len(res))
		l.Next()
		for _, rchar := range repstring {
			nchar := &Char{
				Rune: rune(rchar),
				T:    char.T,
				R:    char.R,
				P:    char.P,
			}
			l.InsertBefore(nchar)
		}
	}
	return nil
}

type Paragraph struct {
	ControlR *xml.UniversalElement
	ControlT *xml.UniversalElement
	xml.UniversalElement
}

func NewParagraph(p, r, t *xml.UniversalElement) *Paragraph {
	return &Paragraph{
		ControlR: r,
		ControlT: t,
		UniversalElement: xml.UniversalElement{
			XMLName: "w:p",
			Attrs:   p.Attrs,
			Children: []*xml.UniversalElement{
				p.GetElementByName("w:pPr"),
				{
					XMLName: "w:r",
					Attrs:   r.Attrs,
					Children: []*xml.UniversalElement{
						r.GetElementByName("w:rPr"),
					},
				},
			},
		},
	}
}

func (p *Paragraph) Insert(char *Char) {
	if char.R != p.ControlR {
		p.ControlR = char.R
		p.Children = append(p.Children, &xml.UniversalElement{
			XMLName: "w:r",
			Attrs:   p.ControlR.Attrs,
			Children: []*xml.UniversalElement{
				p.ControlR.GetElementByName("w:rPr"),
			},
		})
	}
	if char.Rune == 0 {
		return
	}
	if p.Children[p.LenChildren()-1].GetElementByName("w:t") == nil {
		p.ControlT = char.T
		p.Children[p.LenChildren()-1].Children = append(p.Children[p.LenChildren()-1].Children, &xml.UniversalElement{
			XMLName: "w:t",
			Attrs:   p.ControlT.Attrs,
		})
	}

	p.Children[p.LenChildren()-1].GetElementByName("w:t").Data += string(char.Rune)
}

func (p *Paragraph) ToUniversal() *xml.UniversalElement {
	return &p.UniversalElement
}
