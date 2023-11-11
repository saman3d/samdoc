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

	doc, err := NewDocxFromStream(bytes.NewReader(d), int64(len(d)))
	if err != nil {
		return nil, err
	}

	return &Template{File: doc}, nil
}

func (t *Template) rawExecute(model interface{}, exts ...TemplateExecuteExtension) error {
	repfunc, err := NewStructReplacerFunc(model)
	if err != nil {
		return err
	}

	err = t.File.Replace(repfunc)
	if err != nil {
		return err
	}

	for _, ext := range exts {
		err = ext(t)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Template) ExecuteToWriter(model interface{}, writer io.Writer, exts ...TemplateExecuteExtension) error {
	err := t.rawExecute(model, exts...)
	if err != nil {
		return err
	}

	return t.File.Save(writer)
}

func (t *Template) ExecuteToPDF(model interface{}, exts ...TemplateExecuteExtension) ([]byte, error) {
	err := t.rawExecute(model, exts...)
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

type TemplateExecuteExtension func(*Template) error

func WithImageReplaceByFingerprint(ims map[string]io.Reader) TemplateExecuteExtension {
	return func(t *Template) error {
		for fingerprint, image := range ims {
			err := t.File.ReplaceImageByFingerPrint(fingerprint, image)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func WithImageReplaceByName(ims map[string]io.Reader) TemplateExecuteExtension {
	return func(t *Template) error {
		for name, image := range ims {
			err := t.File.ReplaceImageByFingerPrint(name, image)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
