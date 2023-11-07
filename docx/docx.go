package docx

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
)

var (
	ErrCouldntFindWordDoc = errors.New("invalid docx file, couldn't find word document xml")
	ErrImageNotFound      = errors.New("image not found")
)

type ZipInMemory struct {
	data *zip.Reader
}

type Docx struct {
	zipReader  *zip.Reader
	proccessor *Processor
	content    []byte
	headers    map[string][]byte
	footers    map[string][]byte
	images     map[string]string
}

func newDocx(reader *zip.Reader) (*Docx, error) {
	docx := &Docx{
		zipReader:  reader,
		proccessor: new(Processor),
	}
	return docx, docx.load()
}

func (d *Docx) load() error {
	err := d.loadContent()
	if err != nil {
		return err
	}
	d.loadImageFilenames()
	err = d.loadHeadersAndFooters()
	return err
}

func (d *Docx) loadContent() error {
	var err error
	d.content, err = readWordDoc(d.zipReader)
	if err != nil {
		return err
	}
	return nil
}

func (d *Docx) loadHeadersAndFooters() error {
	d.headers = make(map[string][]byte)
	d.footers = make(map[string][]byte)
	for _, f := range d.zipReader.File {
		if strings.Contains(f.Name, "header") {
			fo, err := f.Open()
			if err != nil {
				return err
			}
			h, _ := io.ReadAll(fo)
			d.headers[f.Name] = h
		}
		if strings.Contains(f.Name, "footer") {
			fo, err := f.Open()
			if err != nil {
				return err
			}
			h, _ := io.ReadAll(fo)
			d.footers[f.Name] = h
		}
	}
	return nil
}

func (d *Docx) loadImageFilenames() {
	d.images = make(map[string]string)
	for _, f := range d.zipReader.File {
		if strings.HasPrefix(f.Name, "word/media/") {
			d.images[f.Name] = ""
		}
	}
}

func (d *Docx) WriteToFile(path string) (err error) {
	var target *os.File
	target, err = os.Create(path)
	if err != nil {
		return
	}
	defer target.Close()
	err = d.Save(target)
	return
}

func (d *Docx) Save(ioWriter io.Writer) (err error) {
	w := zip.NewWriter(ioWriter)
	for _, file := range d.zipReader.File {
		var writer io.Writer
		var readCloser io.ReadCloser

		writer, err = w.Create(file.Name)
		if err != nil {
			return err
		}
		readCloser, err = file.Open()
		if err != nil {
			return err
		}
		if file.Name == "word/document.xml" {
			writer.Write([]byte(d.content))
		} else if strings.Contains(file.Name, "header") && len(d.headers[file.Name]) != 0 {
			writer.Write([]byte(d.headers[file.Name]))
		} else if strings.Contains(file.Name, "footer") && len(d.footers[file.Name]) != 0 {
			writer.Write([]byte(d.footers[file.Name]))
		} else if strings.HasPrefix(file.Name, "word/media/") && d.images[file.Name] != "" {
			newImage, err := os.Open(d.images[file.Name])
			if err != nil {
				return err
			}
			writer.Write(streamToByte(newImage))
			newImage.Close()
		} else {
			writer.Write(streamToByte(readCloser))
		}
	}
	w.Close()
	return
}

func (d *Docx) Replace(f ReplacerFunc) error {
	var err error
	d.content, err = d.proccessor.LoadAndReplace(d.content, f)
	if err != nil {
		return err
	}

	for h := range d.headers {
		d.headers[h], err = d.proccessor.LoadAndReplace(d.headers[h], f)
		if err != nil {
			return err
		}
	}

	for foo := range d.footers {
		d.footers[foo], err = d.proccessor.LoadAndReplace(d.footers[foo], f)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Docx) ReplaceImage(oldImage string, newImage string) (err error) {
	if _, ok := d.images[oldImage]; ok {
		d.images[oldImage] = newImage
		return nil
	}
	return ErrImageNotFound
}

func Open(data io.ReaderAt, size int64) (*Docx, error) {
	reader, err := zip.NewReader(data, size)
	if err != nil {
		return nil, err
	}
	return newDocx(reader)
}

func streamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func readWordDoc(r *zip.Reader) ([]byte, error) {
	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			fo, err := f.Open()
			if err != nil {
				return nil, err
			}
			return io.ReadAll(fo)
		}
	}
	return nil, ErrCouldntFindWordDoc
}

func closeAndDelete(f *os.File) error {
	err := f.Close()
	if err != nil {
		return err
	}

	return os.Remove(f.Name())
}
