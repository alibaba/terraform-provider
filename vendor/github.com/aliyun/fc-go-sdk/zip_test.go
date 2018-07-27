package fc

import (
	"testing"
	"os"
	"os/exec"
	"io/ioutil"
	"path/filepath"
)

func TestZipEmptyDir(t *testing.T) {
	// Cleanup the last data to prevent from confliction.
	os.RemoveAll("/tmp/TestZipEmptyDir/")

	// Create directory with empty sub-directory.
	err := os.MkdirAll("/tmp/TestZipEmptyDir/src/empty", 0777)
	if err != nil {
		panic(err)
	}
	f, err := os.Create("/tmp/TestZipEmptyDir/src/t")
	if err != nil {
		panic(err)
	}
	f.Close()

	// Zip the directory.
	err = os.MkdirAll("/tmp/TestZipEmptyDir/dst", 0777)
	if err != nil {
		panic(err)
	}
	zip, err := os.Create("/tmp/TestZipEmptyDir/dst/empty.zip")
	defer zip.Close()
	err = ZipDir("/tmp/TestZipEmptyDir/src", zip)
	if err != nil {
		panic(err)
	}

	// Use system command to unzip the file.
	cmd := exec.Command("sh", "-c", "unzip /tmp/TestZipEmptyDir/dst/empty.zip")
	cmd.Dir = "/tmp/TestZipEmptyDir/dst/"
	_, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%v", err)
	}

	f, err = os.Open("/tmp/TestZipEmptyDir/dst")
	entries, err := f.Readdir(-1)
	// 3 elems: zip file, empty dir and a file.
	if len(entries) != 3 {
		t.Fatalf("%v", entries)
	}
	f.Close()

	// Check the unzipped empty directory is valid.
	info, err := os.Stat("/tmp/TestZipEmptyDir/dst/empty")
	if err != nil {
		panic(err)
	}
	if !info.IsDir() {
		t.Fatalf("%v", info)
	}
}

func TestZipDirWithSymbolLinks(t *testing.T)  {
	// Cleanup the last data to prevent from confliction.
	os.RemoveAll("./TestZipDirWithSymbolLinks/")

	// Create file and symbol links.
	err := os.MkdirAll("./TestZipDirWithSymbolLinks/src/", 0777)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("./TestZipDirWithSymbolLinks/src/main.py", []byte("this is a file."), 0666)
	if err != nil {
		panic(err)
	}
	// Write link with absolute path.
	absPath, _ := filepath.Abs("./TestZipDirWithSymbolLinks/src/main.py")
	err = os.Symlink(absPath, "./TestZipDirWithSymbolLinks/src/symbol_link")
	if err != nil {
		panic(err)
	}

	// Zip the directory.
	err = os.MkdirAll("./TestZipDirWithSymbolLinks/dst", 0777)
	if err != nil {
		panic(err)
	}
	zip, err := os.Create("./TestZipDirWithSymbolLinks/dst/link.zip")
	defer zip.Close()
	err = ZipDir("./TestZipDirWithSymbolLinks/src", zip)
	if err != nil {
		panic(err)
	}

	// Use system command to unzip the file.
	cmd := exec.Command("sh", "-c", "unzip link.zip")
	cmd.Dir = "./TestZipDirWithSymbolLinks/dst/"
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%v: %s", err, out)
	}

	// Verify the file attributes.
	dir, err := os.Open("./TestZipDirWithSymbolLinks/dst")
	entries, err := dir.Readdir(-1)
	// 3 elems: zip file, file and its links.
	if len(entries) != 3 {
		t.Fatalf("%v", entries)
	}
	dir.Close()
	info, err := os.Lstat("./TestZipDirWithSymbolLinks/dst/symbol_link")
	if err != nil {
		panic(err)
	}
	if info.Mode() & os.ModeSymlink == 0 {
		t.Fatalf("%v", info)
	}

	// Verify the file and its symbol link.
	data, err := ioutil.ReadFile("./TestZipDirWithSymbolLinks/dst/main.py")
	if err != nil {
		panic(err)
	}
	if string(data) != "this is a file." {
		t.Fatalf("%s", string(data))
	}
	linkData, err := ioutil.ReadFile("./TestZipDirWithSymbolLinks/dst/symbol_link")
	if err != nil {
		panic(err)
	}
	if string(linkData) != "this is a file." {
		t.Fatalf("%s", string(data))
	}
}
