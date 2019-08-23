package options

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	data, err := ioutil.ReadFile("../../examples/config.yaml")
	assert.Nil(t, err)

	opt, err := New(data)
	assert.Nil(t, err)

	assert.NotNil(t, opt)
}
