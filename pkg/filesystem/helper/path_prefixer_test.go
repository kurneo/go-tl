package helper

import (
	"testing"
)

func setupPreFixer() PathPreFixer {
	prefix := "storage"
	separator := "/"
	return NewPreFixer(prefix, separator)
}

func TestPrefixPath(t *testing.T) {
	preFixer := setupPreFixer()
	path := "/image/image.png"
	expect := "storage/image/image.png"
	actual := preFixer.PrefixPath(path)
	if actual == expect {
		t.Logf("PrefixPath(\"%s\") PASS. Expect \"%s\", got \"%s\"", path, expect, actual)
	} else {
		t.Errorf("PrefixPath(\"%s\") FAILED. Expect \"%s\", got \"%s\"", path, expect, actual)
	}
}

func TestStripPrefix(t *testing.T) {
	preFixer := setupPreFixer()
	path := "storage/image/image.png"
	expect := "image/image.png"
	actual := preFixer.StripPrefix(path)
	if actual == expect {
		t.Logf("StripPrefix(\"%s\") PASS. Expect \"%s\", got \"%s\"", path, expect, actual)
	} else {
		t.Errorf("StripPrefix(\"%s\") FAILED. Expect \"%s\", got \"%s\"", path, expect, actual)
	}
}

func TestStripDirectoryPrefix(t *testing.T) {
	preFixer := setupPreFixer()
	path := "storage/image/"
	expect := "image"
	actual := preFixer.StripDirectoryPrefix(path)
	if actual == expect {
		t.Logf("StripDirectoryPrefix(\"%s\") PASS. Expect \"%s\", got \"%s\"", path, expect, actual)
	} else {
		t.Errorf("StripDirectoryPrefix(\"%s\") FAILED. Expect \"%s\", got \"%s\"", path, expect, actual)
	}
}

func TestPrefixDirectoryPath(t *testing.T) {
	preFixer := setupPreFixer()
	path := "image"
	expect := "storage/image/"
	actual := preFixer.PrefixDirectoryPath(path)
	if actual == expect {
		t.Logf("PrefixDirectoryPath(\"%s\") PASS. Expect \"%s\", got \"%s\"", path, expect, actual)
	} else {
		t.Errorf("PrefixDirectoryPath(\"%s\") FAILED. Expect \"%s\", got \"%s\"", path, expect, actual)
	}

	path = ""
	expect = "storage/"
	actual = preFixer.PrefixDirectoryPath(path)
	if actual == expect {
		t.Logf("PrefixDirectoryPath(\"%s\") PASS. Expect \"%s\", got \"%s\"", path, expect, actual)
	} else {
		t.Errorf("PrefixDirectoryPath(\"%s\") FAILED. Expect \"%s\", got \"%s\"", path, expect, actual)
	}
}
