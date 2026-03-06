// Copyright The OWASP Coraza contributors
// SPDX-License-Identifier: Apache-2.0

package wasmplugin_test

import (
	"io"
	"io/fs"
	"testing"
	"time"

	"github.com/corazawaf/coraza-proxy-wasm/wasmplugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuledataFS(t *testing.T) {
	t.Run("passing a nil map should return an empty filesystem", func(t *testing.T) {
		testfs := wasmplugin.NewRuleDataFS(nil)
		file, err := testfs.Open("nothing.data")
		assert.Nil(t, file)
		require.ErrorIs(t, err, fs.ErrNotExist)
		fileInvalid, err := testfs.Open("../../file")
		assert.Nil(t, fileInvalid)
		require.ErrorIs(t, err, fs.ErrInvalid)
	})

	t.Run("passing a valid map should be usable as a filesystem", func(t *testing.T) {
		originalData := map[string][]byte{
			"file1.data": []byte("somedata1"),
			"file2.data": []byte("somedata2"),
		}
		testfs := wasmplugin.NewRuleDataFS(originalData)
		file, err := testfs.Open("file1.data")
		require.NoError(t, err)
		assert.NotNil(t, file)
		t.Log("reading content of file1.data")
		filecontent, err := io.ReadAll(file)
		require.NoError(t, err)
		assert.Equal(t, []byte("somedata1"), filecontent)
		now := time.Now()
		info, err := file.Stat()
		require.NoError(t, err)
		assert.False(t, info.IsDir())
		assert.Nil(t, info.Sys())
		assert.WithinDuration(t, now, info.ModTime(), time.Second)
		assert.Equal(t, fs.FileMode(0600), info.Mode())
		assert.Equal(t, "file1.data", info.Name())
		assert.Nil(t, file.Close())
		assert.Equal(t, len(originalData["file1.data"]), int(info.Size()))

		t.Log("mutating the data from originalData should not change the existing data on the filesystem")
		originalData["file2.data"] = []byte("mutated this data")
		file2, err := testfs.Open("file2.data")
		require.NoError(t, err)
		assert.NotNil(t, file2)
		t.Log("reading content of file2.data")
		filecontent2, err := io.ReadAll(file2)
		require.NoError(t, err)
		assert.Equal(t, []byte("somedata2"), filecontent2)

	})
}
