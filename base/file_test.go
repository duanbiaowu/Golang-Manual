package base

import (
	"bufio"
	"compress/gzip"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompressFile(t *testing.T) {
	logFile := "/tmp/index.html"
	file, err := os.Open(logFile)
	if err != nil {
		t.Skip()
	}

	reader := bufio.NewReader(file)
	data, err := ioutil.ReadAll(reader)
	assert.Nil(t, err)

	logFile = logFile + ".gz"
	file, err = os.Create(logFile)
	writer := gzip.NewWriter(file)
	_, err = writer.Write(data)
	assert.Nil(t, err)

	_ = writer.Close()
}
