package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalculateNames(t *testing.T) {
	t.Run("export only var, type is exported", func(t *testing.T) {
		t.Parallel()

		e := &enumarr{
			TypeName:   "TestType",
			ExportVar:  true,
			ExportFunc: false,
		}

		e.calculateNames()

		require.Equal(t, "", e.FuncName)
		require.Equal(t, "TestTypeArray", e.VarName)
	})

	t.Run("export var and func, type is exported", func(t *testing.T) {
		t.Parallel()

		e := &enumarr{
			TypeName:   "TestType",
			ExportVar:  true,
			ExportFunc: true,
		}

		e.calculateNames()

		require.Equal(t, "TestTypeAll", e.FuncName)
		require.Equal(t, "TestTypeArray", e.VarName)
	})

	t.Run("export only func, type is exported", func(t *testing.T) {
		t.Parallel()

		e := &enumarr{
			TypeName:   "TestType",
			ExportVar:  false,
			ExportFunc: true,
		}

		e.calculateNames()

		require.Equal(t, "TestTypeAll", e.FuncName)
		require.Equal(t, "_TestTypeArray", e.VarName)
	})
}

// PrepareTestFile creates temporary file, fills it with data and returns path to it.
func PrepareTestFile(t *testing.T, data string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "test.go")

	file, err := os.Create(path)
	require.NoError(t, err, "failed to create file for parsing")
	bytes, err := file.Write([]byte(data))
	require.NoError(t, err, "failed to write test data to file")
	require.Equal(t, len(data), bytes, "failed to write test data to file")
	err = file.Close()
	require.NoError(t, err, "failed to close file")

	return path
}

func TestParse(t *testing.T) {
	t.Run("default configuration", func(t *testing.T) {
		t.Parallel()
		path := PrepareTestFile(t, `package defaultpkg

type DefaultEnum int

const (
  Enum1 DefaultEnum = iota
  Enum2
  Enum3
)`)

		e := &enumarr{
			TypeName: "DefaultEnum",
		}

		err := e.parse(path)
		require.NoError(t, err, "parse shouldn't return error")

		require.Equal(t, "DefaultEnum", e.TypeName, "type")
		require.Equal(t, "defaultpkg", e.Parsed.Pkg, "package")
		require.Equal(t, []string{"Enum1", "Enum2", "Enum3"}, e.Parsed.Names, "enum names")
	})

	t.Run("string enu", func(t *testing.T) {
		t.Parallel()
		path := PrepareTestFile(t, `package stringpkg
type StringEnum string

const (
	Enum1        StringEnum = "enum1"
	NotEnum                 = "it's a trap"
	Enum2, Enum3 StringEnum = "enum2", "enum3"
)
`)
		e := &enumarr{
			TypeName: "StringEnum",
		}
		err := e.parse(path)
		require.NoError(t, err, "parse shouldn't return error")

		require.Equal(t, "StringEnum", e.TypeName, "type")
		require.Equal(t, "stringpkg", e.Parsed.Pkg, "package")
		require.Equal(t, []string{"Enum1", "Enum2", "Enum3"}, e.Parsed.Names, "enum names")
	})
}
