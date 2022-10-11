package logging

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestCreateLogFileSyncer(t *testing.T) {
	defaultOut := &bytes.Buffer{}
	tmpFile, err := ioutil.TempFile("", "test-*")
	got, err := createLogFileSyncer(tmpFile.Name())
	assert.NoError(t, err)
	assert.NotEqual(t, got, zapcore.AddSync(defaultOut))
}