// Apache-2.0 License
// Copyright [2020] [guonaihong]
package brouter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testSplit struct {
	path string
	need []string
}

func Test_SplitPath(t *testing.T) {
	ts := []testSplit{
		{path: "/a/b/c", need: []string{"/a", "/b", "/c"}},
		{path: "/a/b/:name", need: []string{"/a", "/b", "/:name"}},
		{path: "/a/b/c/:name/c/d/:hello", need: []string{"/a", "/b", "/c", "/:name", "/c", "/d", "/:hello"}},
		{path: "/teams/:id/repos", need: []string{"/teams", "/:id", "/repos"}},
	}

	for _, test := range ts {
		got := splitPath(test.path)
		assert.Equal(t, got, test.need)
	}
}

type testGenPath struct {
	path     string
	segments []segment
}

func Test_GenPath(t *testing.T) {
	ts := []testGenPath{
		{
			path: "/teams/:id/repos",
			segments: []segment{
				{
					path:     "/teams/",
					nodeType: ordinary,
				},
				{
					path:      "",
					nodeType:  param,
					paramName: "id",
				},
				{
					path:     "/repos",
					nodeType: ordinary,
				},
			},
		},
		{
			path: "/a/b/c/:name/c/d/:hello",
			segments: []segment{
				{
					path:     "/a/b/c/",
					nodeType: ordinary,
				},
				{
					path:      "",
					nodeType:  param,
					paramName: "name",
				},
				{
					path:     "/c/d/",
					nodeType: ordinary,
				},
				{
					path:      "",
					nodeType:  param,
					paramName: "hello",
				},
			},
		},
		{
			path: "/repos/:owner/:repo/pulls/:number/files",
			segments: []segment{
				{
					path:     "/repos/",
					nodeType: ordinary,
				},
				{
					path:      "",
					nodeType:  param,
					paramName: "owner",
				},
				{
					path:      "/",
					nodeType:  param,
					paramName: "repo",
				},
				{
					path:     "/pulls/",
					nodeType: ordinary,
				},
				{
					path:      "",
					nodeType:  param,
					paramName: "number",
				},
				{
					path:     "/files",
					nodeType: ordinary,
				},
			},
		},
	}

	for _, test := range ts {
		got := genPath(test.path, nil)
		assert.Equal(t, test.segments, got.segments)
	}
}

func Test_GenPath1(t *testing.T) {
	ts := []testGenPath{
		{
			path: "/applications/:client_id/tokens/:access_token",
			segments: []segment{
				{
					path:     "/applications/",
					nodeType: ordinary,
				},
				{
					path:      "",
					nodeType:  param,
					paramName: "client_id",
				},
				{
					path:     "/tokens/",
					nodeType: ordinary,
				},
				{
					path:      "",
					nodeType:  param,
					paramName: "access_token",
				},
			},
		},
	}

	for _, test := range ts {
		got := genPath(test.path, nil)
		assert.Equal(t, test.segments, got.segments)
	}
}
