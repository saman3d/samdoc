package xml

import (
	"fmt"
	"io"
	"regexp"
)

// ---------------------
//     XML Decoder
// ---------------------

var (
	start_tag_reg        = regexp.MustCompile(`^[ \n\t]*?<([^\/?][a-zA-Z:\d]+)((?: {0,}[a-zA-Z:\d\-]+?="[^"]*?")+|)(?:[^\/]|)>`)
	first_start_tag_reg  = regexp.MustCompile(`(?s)^.{0,}?<([^/?][a-zA-Z:\d]+)((?: {0,}[a-zA-Z:\d\-]+?="[^"]*?")+|)(?:[^/]|)>`)
	end_tag_reg          = regexp.MustCompile(`^</([a-zA-Z:\d]+?)>`)
	self_closing_tag_reg = regexp.MustCompile(`(?s)^[ \n\t]*?(<.+?/>)`)
	chardata_reg         = regexp.MustCompile(`(?s)^(.+?)(?:<)`)
	attrs_reg            = regexp.MustCompile(`(?s)(\S+)="(.*?)"`)
	default_indent       = " "
	SearchSize           = 8096
)

type XMLDecoder struct {
	head       int
	b          []byte
	last_token Token
}

func NewXMLDecoder(d []byte) *XMLDecoder {
	return &XMLDecoder{b: d}
}

func Unmarshal(d []byte, model XMLUnmarshaler) error {
	decoder := NewXMLDecoder(d)
	start, err := decoder.forceStartToken()
	if err != nil {
		return err
	}

	return model.XMLUnmarshal(decoder, start)
}

type XMLUnmarshaler interface {
	XMLUnmarshal(d *XMLDecoder, start StartTag) error
}

type Token interface{}

type StartTag struct {
	Tagname string
	Attrs   [][2]string
}

func (t StartTag) String() string {
	attrs := ""
	for _, v := range t.Attrs {
		attrs += fmt.Sprintf(` %s="%s"`, v[0], v[1])
	}
	return fmt.Sprintf("<%s%s>", t.Tagname, attrs)
}

type EndTag struct {
	Tagname string
}

func (t EndTag) String() string {
	return fmt.Sprintf("</%s>", t.Tagname)
}

type CharData string

func (xp *XMLDecoder) Token() (Token, error) {
	var till = SearchSize + xp.head
	if len(xp.b)-1 < xp.head+SearchSize {
		till = len(xp.b) - 1
	}
	fmt.Println(string(xp.b[xp.head:till]))
	matches := start_tag_reg.FindStringSubmatch(string(xp.b[xp.head:till]))
	if len(matches) != 0 {
		xp.head += len(matches[0])
		fmt.Println(matches[1])
		xp.last_token = xp.parseStartTag(matches[1], matches[2])
		return xp.last_token, nil
	}

	if _, ok := xp.last_token.(EndTag); ok {
		for {
			if regexp.MustCompile(`[ \t\n]+`).Match(xp.b[xp.head : xp.head+1]) {
				xp.head++
			} else {
				break
			}
		}
	}

	matches = end_tag_reg.FindStringSubmatch(string(xp.b[xp.head:till]))
	if len(matches) != 0 {
		xp.head += len(matches[0])
		xp.last_token = EndTag{Tagname: matches[1]}
		return xp.last_token, nil
	}

	matches = chardata_reg.FindStringSubmatch(string(xp.b[xp.head:till]))
	if len(matches) != 0 {
		xp.head += len(matches[1])
		xp.last_token = CharData(matches[1])
		return xp.last_token, nil
	}

	return nil, io.EOF
}

func (xp *XMLDecoder) forceStartToken() (StartTag, error) {
	matches := first_start_tag_reg.FindStringSubmatch(string(xp.b[xp.head:]))
	if len(matches) != 0 {
		xp.head += len(matches[0])
		return xp.parseStartTag(matches[1], matches[2]), nil
	}

	return StartTag{}, io.EOF
}

func (xp *XMLDecoder) parseStartTag(tagname string, rawattrs string) StartTag {
	return StartTag{
		Tagname: tagname,
		Attrs:   xp.parseAttrsFromString(rawattrs),
	}
}

func (xp *XMLDecoder) parseAttrsFromString(attrs string) [][2]string {
	var ars = make([][2]string, 0)
	matches := attrs_reg.FindAllStringSubmatch(attrs, -1)
	for _, m := range matches {
		ars = append(ars, [2]string{m[1], m[2]})
	}
	return ars
}

// ---------------------
//     XML Encoder
// ---------------------

type XMLEncoder struct {
	w          []byte
	baseindent int
}

func NewXMLEncoder() *XMLEncoder {
	b := make([]byte, 0)
	header := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	header = append(header, 13)
	b = append(b, header...)
	return &XMLEncoder{
		w:          b,
		baseindent: -1,
	}
}

func Marshal(model XMLMarshaler) ([]byte, error) {
	encoder := NewXMLEncoder()
	return encoder.w, model.XMLMarshal(encoder)
}

type XMLMarshaler interface {
	XMLMarshal(e *XMLEncoder) error
}

func (xp *XMLEncoder) EncodeToken(tkn Token) error {
	switch t := tkn.(type) {
	case StartTag:
		return xp.formatStartTag(t)
	case EndTag:
		return xp.formatEndTag(t)
	case CharData:
		return xp.formatCharData(t)
	}
	return nil
}

func (xp *XMLEncoder) indent() {
	xp.baseindent++
}

func (xp *XMLEncoder) unindent() {
	xp.baseindent--
}

func (xp *XMLEncoder) formatStartTag(t StartTag) error {
	xp.w = append(xp.w, '\n')
	xp.w = append(xp.w, []byte(t.String())...)
	return nil
}

func (xp *XMLEncoder) formatCharData(d CharData) error {
	xp.w = append(xp.w, []byte(d)...)
	return nil
}

func (xp *XMLEncoder) formatEndTag(t EndTag) error {
	xp.w = append(xp.w, []byte(t.String())...)
	return nil
}

// ---------------------
//   Universal Element
// ---------------------

type UniversalElement struct {
	XMLName  string
	Attrs    [][2]string
	Data     string `xml:",chardata"`
	Children []*UniversalElement
}

func (u *UniversalElement) XMLUnmarshal(e *XMLDecoder, start StartTag) error {
	u.Attrs = start.Attrs
	u.XMLName = start.Tagname
	for {
		t, err := e.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		switch tt := t.(type) {
		case StartTag:
			var uniel = &UniversalElement{}
			err = uniel.XMLUnmarshal(e, tt)
			if err != nil {
				return err
			}
			if u.Children == nil {
				u.Children = make([]*UniversalElement, 0)
			}
			u.Children = append(u.Children, uniel)
		case CharData:
			u.Data += string(tt)
		case EndTag:
			return nil
		}
	}

}

func (u *UniversalElement) XMLMarshal(e *XMLEncoder) error {
	t := StartTag{
		Tagname: u.XMLName,
		Attrs:   u.Attrs,
	}
	err := e.EncodeToken(t)
	if err != nil {
		return err
	}
	if u.Data != "" {
		var t = CharData([]byte(u.Data))
		err = e.EncodeToken(t)
		if err != nil {
			return err
		}
	}
	if u.Children != nil {
		for _, c := range u.Children {
			err = c.XMLMarshal(e)
			if err != nil {
				return err
			}
		}
	}
	te := EndTag{
		Tagname: u.XMLName,
	}
	err = e.EncodeToken(te)
	if err != nil {
		return err
	}
	return nil
}

func (u *UniversalElement) LenChildren() int {
	return len(u.Children)
}

func (u *UniversalElement) GetElementByIndx(indx int) *UniversalElement {
	if indx < 0 && indx > u.LenChildren()-1 {
		return nil
	}

	return u.Children[indx]
}

func (u *UniversalElement) GetElementByName(name string) *UniversalElement {
	for i := 0; i < u.LenChildren(); i++ {
		if e := u.GetElementByIndx(i); e.XMLName == name {
			return e
		}
	}
	return nil
}
