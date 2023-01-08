package docx

import (
	"testing"

	"github.com/saman3d/samdoc/xml"
	. "github.com/smartystreets/goconvey/convey"
)

var paragraph_attrs = [][2]string{
	{"w:rsidR", "00000000"},
	{"w:rsidDel", "00000000"},
	{"w:rsidP", "00000000"},
	{"w:rsidRDefault", "00000000"},
	{"w:rsidRPr", "00000000"},
	{"w14:paraId", "000000107"},
}
var paragraph_properties = &xml.UniversalElement{
	XMLName: "w:pPr",
	Data: `<w:pageBreakBefore w:val="0"/>
<w:bidi w:val="1"/>
<w:rPr/>`,
}

var run_properties = &xml.UniversalElement{
	XMLName: "w:rPr",
	Data: `<w:pageBreakBefore w:val="0"/>
<w:bidi w:val="1"/>
<w:rPr/>`}
var doc = &xml.UniversalElement{
	XMLName: "w:document",
	Children: []*xml.UniversalElement{
		{
			XMLName: "w:body",
			Children: []*xml.UniversalElement{
				{
					XMLName: "w:p",
					Attrs:   paragraph_attrs,
					Children: []*xml.UniversalElement{
						paragraph_properties,
						{
							XMLName: "w:r",
							Children: []*xml.UniversalElement{
								run_properties,
								{
									XMLName: "w:t",
									Data:    "hi ",
								},
							},
						},
						{
							XMLName: "w:r",
							Children: []*xml.UniversalElement{
								run_properties,
								{
									XMLName: "w:t",
									Data:    "I'm ",
								},
							},
						},
						{
							XMLName: "w:r",
							Children: []*xml.UniversalElement{
								run_properties,
								{
									XMLName: "w:t",
									Data:    "your ",
								},
							},
						},
						{
							XMLName: "w:r",
							Children: []*xml.UniversalElement{
								run_properties,
								{
									XMLName: "w:t",
									Data:    "friends ",
								},
							},
						},
						{
							XMLName: "w:r",
							Children: []*xml.UniversalElement{
								run_properties,
								{
									XMLName: "w:t",
									Data:    "doctor ",
								},
							},
						},
					},
				},
				{
					XMLName: "w:p",
					Attrs:   paragraph_attrs,
					Children: []*xml.UniversalElement{
						paragraph_properties,
						{
							XMLName: "w:r",
							Children: []*xml.UniversalElement{
								run_properties,
								{
									XMLName: "w:t",
									Data:    "oh ",
								},
							},
						},
						{
							XMLName: "w:r",
							Children: []*xml.UniversalElement{
								run_properties,
								{
									XMLName: "w:t",
									Data:    "I ",
								},
							},
						},
						{
							XMLName: "w:r",
							Children: []*xml.UniversalElement{
								run_properties,
								{
									XMLName: "w:t",
									Data:    "just ",
								},
							},
						},
						{
							XMLName: "w:r",
							Children: []*xml.UniversalElement{
								run_properties,
								{
									XMLName: "w:t",
									Data:    "lied ",
								},
							},
						},
						{
							XMLName: "w:r",
							Children: []*xml.UniversalElement{
								run_properties,
								{
									XMLName: "w:t",
									Data:    "{{",
								},
							},
						},
						{
							XMLName: "w:r",
							Children: []*xml.UniversalElement{
								run_properties,
								{
									XMLName: "w:t",
									Data:    ".Last",
								},
							},
						},
						{
							XMLName: "w:r",
							Children: []*xml.UniversalElement{
								run_properties,
								{
									XMLName: "w:t",
									Data:    "Name",
								},
							},
						},
						{
							XMLName: "w:r",
							Children: []*xml.UniversalElement{
								run_properties,
								{
									XMLName: "w:t",
									Data:    "}} ",
								},
							},
						},
						{
							XMLName: "w:r",
							Children: []*xml.UniversalElement{
								run_properties,
								{
									XMLName: "w:t",
									Data:    "{{.Gender}} ",
								},
							},
						},
						{
							XMLName: "w:r",
							Children: []*xml.UniversalElement{
								run_properties,
								{
									XMLName: "w:t",
									Data:    "{{.Gender}}{{.Fender}}{{.Enter}}",
								},
							},
						},
						{
							XMLName: "w:r",
							Children: []*xml.UniversalElement{
								run_properties,
								{
									XMLName: "w:t",
									Data:    ":D ",
								},
							},
						},
						{
							XMLName: "w:r",
							Children: []*xml.UniversalElement{
								run_properties,
							},
						},
					},
				},
				{
					XMLName:  "w:e",
					Data:     "salam",
					Children: nil,
				},
			},
		},
	},
}

func TestCharlist(t *testing.T) {
	Convey("Test Charlist: Initiation Test", t, func() {
		pl := new(CharList)
		pl.LoadFromElement(doc.Children[0])
		So(pl.Len(), ShouldEqual, 104)

		pl.GoToFirst()
		So(pl.Current().Rune, ShouldEqual, 'h')
		So(pl.Next().Rune, ShouldEqual, 'i')
		So(pl.Next().Rune, ShouldEqual, ' ')
		So(pl.Next().Rune, ShouldEqual, 'I')
		So(pl.Next().Rune, ShouldEqual, '\'')
		So(pl.Next().Rune, ShouldEqual, 'm')

	})

	Convey("Test Charlist: Converting To Paragraph List", t, func() {
		pl := new(CharList)
		pl.LoadFromElement(doc.Children[0])
		So(pl.Len(), ShouldEqual, 104)
		pll := pl.ToParagraphList()
		So(pll, ShouldHaveLength, 2)
		So(pll[0].Children, ShouldHaveLength, 6)
		So(pll[0], ShouldResemble, doc.GetElementByName("w:body").Children[0])
		So(pll[0].Children, ShouldHaveLength, 6)
		So(pll[1], ShouldResemble, doc.GetElementByName("w:body").Children[1])
	})

	Convey("Test Charlist: Testing Look Ahead", t, func() {
		pl := new(CharList)
		pl.LoadFromElement(doc.Children[0])
		So(pl.Current().Rune, ShouldEqual, 'h')
		So(pl.LookAhead(1), ShouldEqual, "i")
		So(pl.LookAhead(5), ShouldEqual, "i I'm")
		pl.GoToFirst()
		lt, err := pl.LookAheadTill("I'm")
		So(err, ShouldBeNil)
		So(lt, ShouldEqual, "i ")
		lt, err = pl.LookAheadTill("your")
		So(err, ShouldBeNil)
		So(lt, ShouldEqual, "i I'm ")
		lt, err = pl.LookAheadTill("doctor")
		So(err, ShouldBeNil)
		So(lt, ShouldEqual, "i I'm your friends ")
	})

	Convey("Test Charlist: Testing Seek", t, func() {
		pl := new(CharList)
		pl.LoadFromElement(doc.Children[0])
		So(pl.Current().Rune, ShouldEqual, 'h')
		pl.Seek(5)
		So(pl.Current().Rune, ShouldEqual, 'm')
		pl.Seek(2)
		So(pl.Current().Rune, ShouldEqual, 'y')
		pl.Seek(3)
		So(pl.Current().Rune, ShouldEqual, 'r')
		pl.Seek(-3)
		So(pl.Current().Rune, ShouldEqual, 'y')
		pl.Seek(-2)
		So(pl.Current().Rune, ShouldEqual, 'm')

		pl.GoToFirst()
		err := pl.SeekTill("ends")
		So(err, ShouldBeNil)
		So(pl.Current().Rune, ShouldEqual, 'i')
		err = pl.SeekTill("s")
		So(err, ShouldBeNil)
		So(pl.Current().Rune, ShouldEqual, 'd')

		pl.GoToFirst()
		err = pl.SeekTill(" friends")
		So(err, ShouldBeNil)
		So(pl.Current().Rune, ShouldEqual, 'r')
	})

	Convey("Test Charlist: Remove", t, func() {
		pl := new(CharList)
		pl.LoadFromElement(doc.Children[0])
		pl.Seek(2)
		err := pl.Remove(4)
		So(err, ShouldBeNil)
		lok := pl.LookAhead(4)
		So(err, ShouldBeNil)
		So(lok, ShouldEqual, "your")
	})

	Convey("Test Charlist: Replacing", t, func() {
		pl := new(CharList)
		pl.LoadFromElement(doc.Children[0])
		pl.Replace(func(i string) (string, bool) { return "saman koushki", true })
		So(pl.Len(), ShouldEqual, 111)
	})

	Convey("Test Charlist: New Replacement Method", t, func() {
		var wdocument = Processor{Document: doc}
		_, err := wdocument.Replace(func(i string) (string, bool) { return "{{" + i + "}}", false })
		So(err, ShouldBeNil)
		So(wdocument.Document, ShouldResemble, doc)
	})

}
