package docx

import (
	"bytes"
	"io"
	"os"
	"os/exec"

	"github.com/saman3d/samdoc"
)

type Template struct {
	File *Docx
}

func NewTemplate(reader io.Reader) (*Template, error) {
	d, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	doc, err := Open(bytes.NewReader(d), int64(len(d)))
	if err != nil {
		return nil, err
	}

	return &Template{File: doc}, nil
}

func (t *Template) rawExecute(model interface{}) error {
	repfunc, err := NewReplacerFunc(model)
	if err != nil {
		return err
	}

	return t.File.Replace(repfunc)
}

func (t *Template) ExecuteToWriter(model interface{}, writer io.Writer) error {
	err := t.rawExecute(model)
	if err != nil {
		return err
	}

	return t.File.Save(writer)
}

func (t *Template) ExecuteToPDF(model interface{}) ([]byte, error) {
	err := t.rawExecute(model)
	if err != nil {
		return nil, err
	}

	docx, err := samdoc.NewTempFile()
	if err != nil {
		return nil, err
	}
	defer closeAndDelete(docx)

	err = t.File.Save(docx)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("lowriter", "--convert-to", "pdf", docx.Name())
	cmd.Dir = samdoc.TEMPDIR
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	pdfile, err := os.Open(docx.Name() + ".pdf")
	if err != nil {
		return nil, err
	}
	defer closeAndDelete(pdfile)

	return io.ReadAll(pdfile)
}
