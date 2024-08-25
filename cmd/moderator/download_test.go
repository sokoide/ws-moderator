package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_downloadFile(t *testing.T) {
	// filename, err := downloadFile("https://oaidalleapiprodscus.blob.core.windows.net/private/org-IceESB4jC7Mxlyq4ovpjNDkY/user-BFasECWehr9D1pQk7kbSlD4g/img-ywy5kGZkynlHz3jZOgEsxceC.png?st=2024-08-25T00%3A07%3A50Z&se=2024-08-25T02%3A07%3A50Z&sp=r&sv=2024-08-04&sr=b&rscd=inline&rsct=image/png&skoid=d505667d-d6c1-4a0a-bac7-5c84a87759f8&sktid=a48cca56-e6da-484e-a814-9c849652bcb3&skt=2024-08-24T23%3A44%3A16Z&ske=2024-08-25T23%3A44%3A16Z&sks=b&skv=2024-08-04&sig=H9On2h%2BBcp15I8tUZQlFFQpE3gipNfW%2Bmm9bw3mDyAs%3D", "./tmp", "hoge@example.com")
	// assert.Nil(t, err)
	// assert.NotEqual(t, "", filename)
}

func Test_downloadFileDir(t *testing.T) {
	type testCase struct {
		base  string
		input string
		want  string
	}

	for _, c := range []testCase{
		{"./", "hoge#@page#.com", "hoge__page_.com"},
		{"/path/to", "foo@bar.com", "/path/to/foo_bar.com"},
		{"/path/to/", "foo@bar.com", "/path/to/foo_bar.com"}} {
		got := downloadFileDir(c.base, c.input)
		assert.Equal(t, c.want, got)
	}
}

func Test_convertEmailToPath(t *testing.T) {
	type testCase struct {
		input string
		want  string
	}

	for _, c := range []testCase{
		{"foo!@bar.com", "foo__bar.com"},
	} {
		got := convertEmailToPath(c.input)
		assert.Equal(t, c.want, got)
	}

}
