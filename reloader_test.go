package onlineconf

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidReloader(t *testing.T) {
	assert := assert.New(t)

	_, err := NewModuleReloader(&ReloaderOptions{Name: "NoSuchModule"})
	log.Printf("reloader err: %#v", err)
	assert.NotNil(err)
}
