// Copyright The OWASP Coraza contributors
// SPDX-License-Identifier: Apache-2.0

package wasmplugin

import (
	"bytes"
	"io/fs"
	"time"
)

// This package implements fs.FS, but based on a map[string][]byte, that is
// returned from a cache server.

// It is used by a Coraza instance to load ruledata on .WithRootFS(root)

// ruledataFS implements a read-only fs.FS to be used by Coraza
// We do not implement a mutex on it because it is instantiated once, and is just
// readable. Any attempt to write to it is an error.
type ruledataFS struct {
	content map[string][]byte
}

// NewRuleDataFS creates a new instance of an in-memory filesystem that implements fs.FS
func NewRuleDataFS(data map[string][]byte) *ruledataFS {
	internalContentData := make(map[string][]byte)
	for file, value := range data {
		internalContentData[file] = bytes.Clone(value)
	}

	return &ruledataFS{content: internalContentData}
}

// Open implements the fs.FS interface (Read-Only access)
func (m *ruledataFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}
	data, ok := m.content[name]
	if !ok {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}

	return &ruleDataFile{
		name:   name,
		Reader: bytes.NewReader(data),
		size:   int64(len(data)),
	}, nil
}

// ruleDataFile wraps bytes.Reader to satisfy the fs.File interface
type ruleDataFile struct {
	name string
	*bytes.Reader
	size int64
}

// Stat returns a FileInfo from a file
func (f *ruleDataFile) Stat() (fs.FileInfo, error) { return f, nil }

// Close closes the reader of a file
func (f *ruleDataFile) Close() error { return nil }

// Name returns a file name
func (f *ruleDataFile) Name() string { return f.name }

// Mode returns a file access mode
func (f *ruleDataFile) Mode() fs.FileMode { return 0600 }

// ModTime returns a file modification time
func (f *ruleDataFile) ModTime() time.Time { return time.Now() }

// IsDir returns a boolean that represents if a file entry is a directory
func (f *ruleDataFile) IsDir() bool { return false }

// Sys does not return anything
func (f *ruleDataFile) Sys() any { return nil }

// Size returns the size of a file content
func (f *ruleDataFile) Size() int64 { return f.size }
