package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_storeCompletion(t *testing.T) {
	err := storeCompletion("test_title", "test_user", "test_employee", "test_email", "66cae2ef28219394bf3ebce6", "66cae32128219394bf3ebce8")
	assert.Nil(t, err)
}
