package docx

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testFile = "./TestDocument.docx"
const testFileResult = "./TestDocumentResult.docx"
const testOldImage = "word/media/image1.png"
const testNewImage = "./TestImage.png"

func ReadFile(path string) (io.ReadCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func TestReplaceImage(t *testing.T) {
	reader, err := ReadFile(testFile)
	assert.Nil(t, err)
	assert.NotNil(t, reader)
	tmp, err := NewTemplate(reader)
	assert.Nil(t, err)
	assert.NotNil(t, tmp)
	tmp.File.ReplaceImage(testOldImage, testNewImage)
	tmp.File.WriteToFile(testFileResult)

	reader, err = ReadFile(testFileResult)
	assert.Nil(t, err)
	assert.NotNil(t, reader)
	tmp, err = NewTemplate(reader)
	assert.Nil(t, err)
	assert.NotNil(t, tmp)

	_, ok := tmp.File.images[testOldImage]
	assert.Equal(t, true, ok)
}
