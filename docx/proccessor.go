package docx

import (
	"github.com/saman3d/samdoc/xml"
)

type Processor struct {
	Document *xml.UniversalElement
}

func (p *Processor) LoadElement(doc *xml.UniversalElement) {
	p.Document = doc
}

func (p *Processor) LoadAndReplace(inp []byte, f ReplacerFunc) ([]byte, error) {
	var contentRoot xml.UniversalElement
	err := xml.Unmarshal(inp, &contentRoot)
	if err != nil {
		return nil, err
	}
	p.Document = &contentRoot
	return p.Replace(f)
}

func (p *Processor) Replace(repfunc ReplacerFunc) ([]byte, error) {
	err := p.WalkAndReplace(p.Document, repfunc)
	if err != nil {
		return nil, err
	}

	return xml.Marshal(p.Document)
}

func (p *Processor) WalkAndReplace(start *xml.UniversalElement, repf ReplacerFunc) error {
	for _, child := range start.Children {
		if len(child.Children) == 0 {
			continue
		}
		if e := child.GetElementByIndx(0); p != nil && e.XMLName == "w:p" {
			err := p.ProccessReplace(child, repf)
			if err != nil {
				return err
			}
			continue
		}
		err := p.WalkAndReplace(child, repf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Processor) ProccessReplace(con *xml.UniversalElement, repf ReplacerFunc) error {
	list := new(CharList)
	list.LoadFromElement(con)

	if list.Len() == 0 {
		return nil
	}

	err := list.Replace(repf)
	if err != nil {
		return err
	}

	pl := list.ToParagraphList()
	con.Children = append(pl, con.Children[len(pl):]...)
	return nil
}
