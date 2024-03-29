package docx

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testFile = "./TestDocument.docx"
const testFileResult = "./TestDocumentResult.docx"
const testOldImage = "word/media/image1.png"
const newTestImage = "./NewTestImage.png"
const oldTestImage = "./OldTestImage.png"

func ReadFile(path string) (io.ReadCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func TestReplaceImageByImageName(t *testing.T) {
	readerTestFile, err := ReadFile(testFile)
	assert.Nil(t, err)
	assert.NotNil(t, readerTestFile)
	testFileTemplate, err := NewTemplate(readerTestFile)
	assert.Nil(t, err)
	assert.NotNil(t, testFileTemplate)
	result, err := os.Create(testFileResult)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	err = testFileTemplate.ExecuteToWriter(nil, result, WithImageReplaceByName(map[string]io.Reader{
		testOldImage: filePathToReader(newTestImage),
	}))
	assert.Nil(t, err)
	result.Close()

	readerTestFileResult, err := ReadFile(testFileResult)
	assert.Nil(t, err)
	assert.NotNil(t, readerTestFileResult)
	testFileResultTemplate, err := NewTemplate(readerTestFileResult)
	assert.Nil(t, err)
	assert.NotNil(t, testFileResultTemplate)

	newTestImageFingerprint := filePathToFingerprint(newTestImage)
	assert.NotEqual(t, "", newTestImageFingerprint)
	assert.Equal(t, true, testFileResultTemplate.File.images.Has(newTestImageFingerprint))
}

func TestReplaceImageByFingerPrint(t *testing.T) {
	reader, err := ReadFile(testFile)
	assert.Nil(t, err)
	assert.NotNil(t, reader)
	tmp, err := NewTemplate(reader)
	assert.Nil(t, err)
	assert.NotNil(t, tmp)

	type S struct {
		ContractDate string
	}

	result, err := os.Create(testFileResult)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	err = tmp.ExecuteToWriter(&S{ContractDate: "saman"}, result, WithImageReplaceByFingerprint(map[string]io.Reader{
		filePathToFingerprint(oldTestImage): filePathToReader(newTestImage),
	}))
	assert.Nil(t, err)
	result.Close()

	readerTestFileResult, err := ReadFile(testFileResult)
	assert.Nil(t, err)
	assert.NotNil(t, readerTestFileResult)
	testFileResultTemplate, err := NewTemplate(readerTestFileResult)
	assert.Nil(t, err)
	assert.NotNil(t, testFileResultTemplate)

	newTestImageFingerprint := filePathToFingerprint(newTestImage)
	assert.NotEqual(t, "", newTestImageFingerprint)
	assert.Equal(t, true, testFileResultTemplate.File.images.Has(newTestImageFingerprint))
}

func filePathToFingerprint(path string) string {
	file, err := os.Open(path)
	if nil == err {
		data := streamToByte(file)
		if nil == err {
			h := sha256.New()
			if _, err := h.Write(data); err == nil {
				return fmt.Sprintf("%x", h.Sum(nil))
			}
		}
	}
	return ""
}

func filePathToReader(path string) io.Reader {
	file, err := os.Open(path)
	if nil == err {
		return file
	}
	return nil
}
