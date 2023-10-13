package utils

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyFile(t *testing.T) {
	// Create a temporary source file
	srcFile, err := ioutil.TempFile("", "src")
	assert.NoError(t, err)
	defer os.Remove(srcFile.Name())

	// Write some content to the source file
	srcContent := "test content"
	_, err = srcFile.Write([]byte(srcContent))
	assert.NoError(t, err)
	srcFile.Close()

	// Create a temporary destination file
	dstFile, err := ioutil.TempFile("", "dst")
	assert.NoError(t, err)
	defer os.Remove(dstFile.Name())
	dstFile.Close()

	// Call the CopyFile function
	err = CopyFile(srcFile.Name(), dstFile.Name())
	assert.NoError(t, err)

	// Read the content of the destination file
	dstContent, err := ioutil.ReadFile(dstFile.Name())
	assert.NoError(t, err)

	// Check that the content of the destination file is the same as the source file
	assert.Equal(t, srcContent, string(dstContent))
}
