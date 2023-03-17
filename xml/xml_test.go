package xml

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestXML(t *testing.T) {
	Convey("Test xml: attrs", t, func() {
		var x = &XMLDecoder{}
		attrs := x.parseAttrsFromString(`asf="sad"  rr="ff" w:fsdf="fff"`)
		So(attrs, ShouldResemble, [][2]string{{"asf", "sad"}, {"rr", "ff"}, {"w:fsdf", "fff"}})

		attrs = x.parseAttrsFromString(`          asf="sad"            rr="ff"           w:fsdf="fff"          `)
		So(attrs, ShouldResemble, [][2]string{{"asf", "sad"}, {"rr", "ff"}, {"w:fsdf", "fff"}})
	})

	Convey("Test xml: read xml", t, func() {
		data := `
		<w:p w:rsidR="00000000" w:rsidDel="00000000" w:rsidP="00000000" w:rsidRDefault="00000000" w:rsidRPr="00000000" w14:paraId="0000000D" flex-direction="row">
<w:pPr>
<w:bidi w:val="1"/>
<w:rPr/>
</w:pPr>
<w:r w:rsidDel="00000000" w:rsidR="00000000" w:rsidRPr="00000000">
<w:rPr>
<w:rtl w:val="0"/>
</w:rPr>
</w:r>
</w:p>
`
		var parser = NewXMLDecoder([]byte(data))
		t, err := parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldHaveSameTypeAs, StartTag{})
		tt := t.(StartTag)
		So(tt.Tagname, ShouldEqual, "w:p")
		So(tt.Attrs, ShouldResemble, [][2]string{
			{"w:rsidR", "00000000"},
			{"w:rsidDel", "00000000"},
			{"w:rsidP", "00000000"},
			{"w:rsidRDefault", "00000000"},
			{"w:rsidRPr", "00000000"},
			{"w14:paraId", "0000000D"},
			{"flex-direction", "row"},
		})

		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldHaveSameTypeAs, StartTag{})
		tt = t.(StartTag)
		So(tt.Tagname, ShouldEqual, "w:pPr")
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldEqual, CharData("\n"))
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldEqual, CharData("<w:bidi w:val=\"1\"/>\n"))
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldEqual, CharData("<w:rPr/>\n"))
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldHaveSameTypeAs, EndTag{})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldHaveSameTypeAs, StartTag{})
		tt = t.(StartTag)
		So(tt.Tagname, ShouldEqual, "w:r")
		parser.Token()
		parser.Token()
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldHaveSameTypeAs, CharData(""))
		So(t, ShouldEqual, `<w:rtl w:val="0"/>
`)
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldHaveSameTypeAs, EndTag{})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldHaveSameTypeAs, EndTag{})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldHaveSameTypeAs, EndTag{})
	})

}
