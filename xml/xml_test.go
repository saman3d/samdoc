package xml

import (
	"io"
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

func TestHtml(t *testing.T) {
	var SimpleHtmlTemplate = `
	<html>
	<head>
		<title>Simple HTML Template</title>
	</head>
	<body display="flex" flex-direction="column">
		<div display="flex" flex-direction="row">
			<div flex="1" border="true" display="block"></div>
			<div flex="1" border="true" display="block"></div>
		</div>
					<div flex="3" display="flex" flex-direction="column">
						<div border="true" display="flex" max-height="2" id="table_header">
							<p id="pid">pid</p>
							<p id="name">name</p>
							<p id="cpu">cpu</p>
							<p id="mem">mem</p>
		</div>
						<div border="true" display="block"></div>
						<div border="true" display="block" max-height="2"></div>
					</div>
	</body>
	</html>`
	Convey("Test xml: read xml", t, func() {
		var parser = NewXMLDecoder([]byte(SimpleHtmlTemplate))
		t, err := parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, StartTag{Tagname: "html", Attrs: [][2]string{}})
		tt := t.(StartTag)
		So(tt.Tagname, ShouldEqual, "html")
		So(tt.Attrs, ShouldResemble, [][2]string{})

		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, StartTag{Tagname: "head", Attrs: [][2]string{}})
		tt = t.(StartTag)
		So(tt.Tagname, ShouldEqual, "head")
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, StartTag{Tagname: "title", Attrs: [][2]string{}})
		tt = t.(StartTag)
		So(tt.Tagname, ShouldEqual, "title")
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldHaveSameTypeAs, CharData("Simple HTML Template"))
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, EndTag{Tagname: "title"})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, EndTag{Tagname: "head"})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, StartTag{Tagname: "body", Attrs: [][2]string{
			{"display", "flex"},
			{"flex-direction", "column"},
		}})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, StartTag{Tagname: "div", Attrs: [][2]string{
			{"display", "flex"},
			{"flex-direction", "row"},
		}})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, StartTag{Tagname: "div", Attrs: [][2]string{
			{"flex", "1"},
			{"border", "true"},
			{"display", "block"},
		}})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, EndTag{Tagname: "div"})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, StartTag{Tagname: "div", Attrs: [][2]string{
			{"flex", "1"},
			{"border", "true"},
			{"display", "block"},
		}})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, EndTag{Tagname: "div"})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, EndTag{Tagname: "div"})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, StartTag{Tagname: "div", Attrs: [][2]string{
			{"flex", "3"},
			{"display", "flex"},
			{"flex-direction", "column"},
		}})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, StartTag{Tagname: "div", Attrs: [][2]string{
			{"border", "true"},
			{"display", "flex"},
			{"max-height", "2"},
			{"id", "table_header"},
		}})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, StartTag{Tagname: "p", Attrs: [][2]string{
			{"id", "pid"},
		}})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldHaveSameTypeAs, CharData("pid"))
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, EndTag{Tagname: "p"})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, StartTag{Tagname: "p", Attrs: [][2]string{
			{"id", "name"},
		}})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldHaveSameTypeAs, CharData("name"))
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, EndTag{Tagname: "p"})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, StartTag{Tagname: "p", Attrs: [][2]string{
			{"id", "cpu"},
		}})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldHaveSameTypeAs, CharData("cpu"))
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, EndTag{Tagname: "p"})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, StartTag{Tagname: "p", Attrs: [][2]string{
			{"id", "mem"},
		}})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldHaveSameTypeAs, CharData("mem"))
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, EndTag{Tagname: "p"})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, EndTag{Tagname: "div"})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, StartTag{Tagname: "div", Attrs: [][2]string{
			{"border", "true"},
			{"display", "block"},
		}})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, EndTag{Tagname: "div"})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, StartTag{Tagname: "div", Attrs: [][2]string{
			{"border", "true"},
			{"display", "block"},
			{"max-height", "2"},
		}})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, EndTag{Tagname: "div"})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, EndTag{Tagname: "div"})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, EndTag{Tagname: "body"})
		t, err = parser.Token()
		So(err, ShouldBeNil)
		So(t, ShouldResemble, EndTag{Tagname: "html"})
		t, err = parser.Token()
		So(err, ShouldEqual, io.EOF)
		So(t, ShouldBeNil)

	})

}
