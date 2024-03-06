package datasource

import (
	"strings"
	"testing"
)

func setupPathHelper() PathHelper {
	return PathHelper{
		separator: "/",
	}
}

func TestStripSlash(t *testing.T) {
	helper := setupPathHelper()
	assertMessage := func(t *testing.T, path, expect, got string) {
		t.Helper()
		if expect == got {
			t.Logf("StripSlash(\"%s\", \".\\/\") PASS. Expected \"%s\", get \"%s\"", path, expect, got)
		} else {
			t.Errorf("StripSlash(\"%s\", \".\\/\") FAILED. Expected \"%s\", get \"%s\"", path, expect, got)
		}
	}
	t.Run("path have prefix slash and dot", func(t *testing.T) {
		path := "./images/image.png"
		expect := "images/image.png"
		got := helper.StripSlash(path, "./")
		assertMessage(t, path, expect, got)

	})

	t.Run("path have both prefix and postix slash", func(t *testing.T) {
		path := "./images/"
		expect := "images"
		got := helper.StripSlash(path, "./")
		assertMessage(t, path, expect, got)
	})
}

func TestDirPath(t *testing.T) {
	helper := setupPathHelper()
	assertMessage := func(t *testing.T, path, expect, got string) {
		t.Helper()
		if expect == got {
			t.Logf("DirPath(\"%s\") PASS. Expected \"%s\", got \"%s\"", path, expect, got)
		} else {
			t.Errorf("DirPath(\"%s\") FAILED. Expected \"%s\", got \"%s\"", path, expect, got)
		}
	}
	t.Run("path with directory and file name", func(t *testing.T) {
		path := "images/image.png"
		expect := "images"
		got := helper.DirPath(path)
		assertMessage(t, path, expect, got)
	})

	t.Run("path with only file name", func(t *testing.T) {
		path := "image.png"
		expect := ""
		got := helper.DirPath(path)
		assertMessage(t, path, expect, got)
	})
}

func TestConcat(t *testing.T) {
	helper := setupPathHelper()
	paths := []string{"a", "b", "c"}
	expect := "a/b/c"
	got := helper.Concat(paths...)

	if expect == got {
		t.Logf("Concat([%s]) PASS. Expected \"%s\", got \"%s\"", strings.Join(paths, ","), expect, got)
	} else {
		t.Errorf("Concat([%s]) FAILED. Expected \"%s\", got \"%s\"", strings.Join(paths, ","), expect, got)
	}
}
