package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetAPIVersion(t *testing.T) {
	assert.Equal(t, "1",getAPIVersion("1.0.0"), "getting version")
	assert.Equal(t, "0",getAPIVersion("1.0"), "getting version")
}
