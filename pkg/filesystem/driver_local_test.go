package filesystem

import (
	"github.com/kurneo/go-template/pkg/filesystem/helper"
	"log"
	"os"
	"testing"
)

func setupLocalDriver() *LocalDriver {
	separator := "/"
	prefix := "./storage/testing/unit"

	err := os.MkdirAll(prefix, 0777)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(prefix+"/test", 0777)
	if err != nil {
		log.Fatal(err)
	}

	f, _ := os.Create(prefix + "/test.txt")
	_, err = f.WriteString("test file")
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}

	return &LocalDriver{preFixer: helper.NewPreFixer(prefix, separator)}
}

func teardownLocalDriver() {
	prefix := "./storage"
	err := os.RemoveAll(prefix)
	if err != nil {
		log.Fatal(err)
	}
}

func TestExist(t *testing.T) {
	d := setupLocalDriver()
	defer func() { teardownLocalDriver() }()
	filename := "test.txt"
	result, _ := d.FileExists(filename)
	if result {
		t.Logf("FileExists(\"%s\") PASS. Expected %t, got %t\n", filename, true, result)
	} else {
		t.Errorf("FileExists(\"%s\") FAILED. Expected %t, got %t\n", filename, true, result)
	}

	filename = "test2.txt"
	result, _ = d.FileExists(filename)
	if !result {
		t.Logf("FileExists(\"%s\") PASS. Expected %t, got %t\n", filename, false, result)
	} else {
		t.Errorf("FileExists(\"%s\") FAILED. Expected %t, got %t\n", filename, false, result)
	}

	dirName := "test"
	result, _ = d.DirExists(dirName)
	if result {
		t.Logf("DirExists(\"%s\") PASS. Expected %t, got %t\n", dirName, true, result)
	} else {
		t.Errorf("DirExists(\"%s\") FAILED. Expected %t, got %t\n", dirName, true, result)
	}

	dirName = "test2"
	result, _ = d.DirExists(dirName)
	if !result {
		t.Logf("DirExists(\"%s\") PASS. Expected %t, got %t\n", dirName, true, result)
	} else {
		t.Errorf("DirExists(\"%s\") FAILED. Expected %t, got %t\n", dirName, true, result)
	}
}

func TestPut(t *testing.T) {
	d := setupLocalDriver()
	defer func() { teardownLocalDriver() }()
	fileName := "data.txt"
	content := "data test"
	err := d.Put(fileName, []byte(content))
	if err == nil {
		t.Logf("Put(\"%s\") PASS. Expected error nil, got nil\n", content)
	} else {
		t.Errorf("Put(\"%s\") FAILED. Expected error nil, got error \"%s\"\n", content, err.Error())
	}
}

func TestGet(t *testing.T) {
	d := setupLocalDriver()
	defer func() { teardownLocalDriver() }()
	fileName := "test.txt"
	expect := "test file"
	content, err := d.Get(fileName)
	if err == nil {
		t.Logf("Get(\"%s\") PASS. Expected error nil, got nil\n", fileName)
	} else {
		t.Errorf("Get(\"%s\") FAILED. Expected error nil, got error \"%s\"\n", fileName, err.Error())
	}

	if expect == string(content) {
		t.Logf("Get(\"%s\") PASS. Expected \"%s\", got \"%s\"\n", fileName, expect, content)
	} else {
		t.Errorf("Get(\"%s\") FAILED. Expected \"%s\", got \"%s\"\n", fileName, expect, content)
	}
}

func TestMakeDir(t *testing.T) {
	d := setupLocalDriver()
	defer func() { teardownLocalDriver() }()
	dirName := "dir name"
	var perm os.FileMode = 777
	err := d.MakeDir(dirName, perm)

	if err == nil {
		t.Logf("MakeDir(\"%s\", %d) PASS. Expected error nil, got nil\n", dirName, perm)
	} else {
		t.Errorf("MakeDir(\"%s\", %d) FAILED. Expected error nil, got error \"%s\"\n", dirName, perm, err.Error())
	}

	result, _ := d.DirExists(dirName)
	if result {
		t.Logf("DirExists(\"%s\") PASS. Expected %t, got %t\n", dirName, true, result)
	} else {
		t.Errorf("DirExists(\"%s\") FAILED. Expected %t, got %t\n", dirName, true, result)
	}
}

func TestDelete(t *testing.T) {
	d := setupLocalDriver()
	defer func() { teardownLocalDriver() }()
	filename := "test.txt"

	err := d.Delete(filename)

	if err == nil {
		t.Logf("Delete(\"%s\") PASS. Expected error nil, got nil\n", filename)
	} else {
		t.Errorf("Delete(\"%s\") FAILED. Expected error nil, got error \"%s\"\n", filename, err.Error())
	}

	result, _ := d.FileExists(filename)
	if !result {
		t.Logf("FileExists(\"%s\") PASS. Expected %t, got %t\n", filename, false, result)
	} else {
		t.Errorf("FileExists(\"%s\") FAILED. Expected %t, got %t\n", filename, false, result)
	}
}

func TestRename(t *testing.T) {
	d := setupLocalDriver()
	defer func() { teardownLocalDriver() }()
	filename := "test.txt"
	renameTo := "test2.txt"

	err := d.Rename(filename, renameTo)

	if err == nil {
		t.Logf("Rename(\"%s\", \"%s\") PASS. Expected error nil, got nil\n", filename, renameTo)
	} else {
		t.Errorf("Rename(\"%s\", \"%s\") FAILED. Expected error nil, got error \"%s\"\n", filename, renameTo, err.Error())
	}

	result, _ := d.FileExists(renameTo)
	if result {
		t.Logf("FileExists(\"%s\") PASS. Expected %t, got %t\n", renameTo, true, result)
	} else {
		t.Errorf("FileExists(\"%s\") FAILED. Expected %t, got %t\n", renameTo, true, result)
	}

	result, _ = d.FileExists(filename)
	if !result {
		t.Logf("FileExists(\"%s\") PASS. Expected %t, got %t\n", filename, false, result)
	} else {
		t.Errorf("FileExists(\"%s\") FAILED. Expected %t, got %t\n", filename, false, result)
	}
}

func TestListContents(t *testing.T) {
	d := setupLocalDriver()
	defer func() { teardownLocalDriver() }()
	path := ""
	files, dirs, err := d.ListContents(path)

	if err == nil {
		t.Logf("ListContents(\"%s\") PASS. Expected error nil, got nil\n", path)
	} else {
		t.Errorf("ListContents(\"%s\") FAILED. Expected error nil, got error \"%s\"\n", path, err.Error())
	}

	if len(files) == 1 {
		t.Logf("ListContents(\"%s\") PASS. Expected count file 1, got %d\n", path, len(files))
	} else {
		t.Errorf("ListContents(\"%s\") FAILED. Expected count files 1, got %d\n", path, len(files))
	}

	if len(dirs) == 1 {
		t.Logf("ListContents(\"%s\") PASS. Expected count directories 1, got %d\n", path, len(dirs))
	} else {
		t.Errorf("ListContents(\"%s\") FAILED. Expected count directories 1, got %d\n", path, len(dirs))
	}
}

func TestCopy(t *testing.T) {
	d := setupLocalDriver()
	defer func() { teardownLocalDriver() }()

	fileName := "test.txt"
	copyFile := "test2.txt"

	err := d.Copy(fileName, copyFile)

	if err == nil {
		t.Logf("Copy(\"%s\", \"%s\") PASS. Expected error nil, got nil\n", fileName, copyFile)
	} else {
		t.Errorf("Copy(\"%s\", \"%s\") FAILED. Expected error nil, got error \"%s\"\n", fileName, copyFile, err.Error())
	}

	result, _ := d.FileExists(copyFile)
	if result {
		t.Logf("FileExists(\"%s\") PASS. Expected %t, got %t\n", copyFile, true, result)
	} else {
		t.Errorf("FileExists(\"%s\") FAILED. Expected %t, got %t\n", copyFile, true, result)
	}

	fileName = "test.txt"
	copyFile = "copy/test2.txt"

	err = d.Copy(fileName, copyFile)

	if err == nil {
		t.Logf("Copy(\"%s\", \"%s\") PASS. Expected error nil, got nil\n", fileName, copyFile)
	} else {
		t.Errorf("Copy(\"%s\", \"%s\") FAILED. Expected error nil, got error \"%s\"\n", fileName, copyFile, err.Error())
	}

	result, _ = d.FileExists(copyFile)
	if result {
		t.Logf("FileExists(\"%s\") PASS. Expected %t, got %t\n", copyFile, true, result)
	} else {
		t.Errorf("FileExists(\"%s\") FAILED. Expected %t, got %t\n", copyFile, true, result)
	}
}
