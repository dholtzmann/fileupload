package fileupload

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func Test_saveThumbnail(t *testing.T) {
	tempDir := "testing-filevalidator"
	testName := "test.jpg"
	dir, err := ioutil.TempDir("", tempDir) // make a temp directory
	if err != nil {
		t.Fatalf("ioutil.TempDir: %s", err)
	}
	defer os.RemoveAll(dir) // delete the temp directory

	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(gopherPNG))
	buf, err := ioutil.ReadAll(decoder)
	if err != nil {
		t.Errorf("base64.NewDecoder(): Error reading base64 encoded file! [%s]", err)
	}
	if err := saveThumbnail(buf, dir, testName, 250, 340); err != nil {
		t.Errorf("saveThumbnail(): Returned an error! [%s]", err)
	}

	if _, err := os.Stat(dir + string(os.PathSeparator) + testName); err != nil { // check if the file exists
		t.Errorf("saveThumbnail(): Thumbnail does not exist! [%s]", err)
	}
}

func Test_saveImage(t *testing.T) {
	tempDir := "testing-filevalidator"
	testName := "test.jpg"
	dir, err := ioutil.TempDir("", tempDir) // make a temp directory
	if err != nil {
		t.Fatalf("ioutil.TempDir: %s", err)
	}
	defer os.RemoveAll(dir) // delete the temp directory

	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(gopherPNG))
	buf, err := ioutil.ReadAll(decoder)
	if err != nil {
		t.Errorf("base64.NewDecoder(): Error reading base64 encoded file! [%s]", err)
	}

	fi, err := saveImage(buf, "gopher.png", testName, dir)
	if err != nil {
		t.Errorf("saveImage(): Returned an error! [%s]", err)
	}
	if _, err := os.Stat(fi.Directory + string(os.PathSeparator) + fi.Name); err != nil { // check if the file exists
		t.Errorf("saveImage(): Image does not exist! [%s]", err)
	}
}

func Test_UploadImageWithThumbnail(t *testing.T) {
	tempDir1, tempDir2 := "testing-filevalidator-images", "testing-filevalidator-thumbnails"
	dir1, err := ioutil.TempDir("", tempDir1) // make two temporary directories
	if err != nil {
		t.Fatalf("ioutil.TempDir: %s", err)
	}
	dir2, err := ioutil.TempDir("", tempDir2)
	if err != nil {
		t.Fatalf("ioutil.TempDir: %s", err)
	}
	defer os.RemoveAll(dir1) // delete the temp directories
	defer os.RemoveAll(dir2)

	req := setupRequestMultipartForm(&testFile{gopherPNG, "imageupload", "gopher.png"}, &testFile{blueJPG, "imageupload", "blue.jpg"})

	for _, fieldSlice := range req.MultipartForm.File { // map[string][]*FileHeader
		for _, header := range fieldSlice {
			fi, err := UploadImageWithThumbnail(header, dir1, dir2)
			if err != nil {
				t.Errorf("UploadImageWithThumbnail(): Returned an error! [%s]", err)
			}
			if _, err := os.Stat(dir1 + string(os.PathSeparator) + fi.Name); err != nil { // check if the image exists
				t.Errorf("UploadImageWithThumbnail(): Image does not exist! [%s]", err)
			}
			if _, err := os.Stat(dir2 + string(os.PathSeparator) + fi.Name); err != nil { // check if the thumbnail exists
				t.Errorf("UploadImageWithThumbnail(): Thumbnail does not exist! [%s]", err)
			}
		}
	}
}

func Test_UploadAllImages(t *testing.T) {
	tempDir1, tempDir2 := "testing-filevalidator-images", "testing-filevalidator-thumbnails"
	dir1, err := ioutil.TempDir("", tempDir1) // make two temporary directories
	if err != nil {
		t.Fatalf("ioutil.TempDir: %s", err)
	}
	dir2, err := ioutil.TempDir("", tempDir2)
	if err != nil {
		t.Fatalf("ioutil.TempDir: %s", err)
	}
	defer os.RemoveAll(dir1) // delete the temp directories
	defer os.RemoveAll(dir2)

	req := setupRequestMultipartForm(&testFile{gopherPNG, "imageupload", "gopher.png"}, &testFile{blueJPG, "imageupload", "blue.jpg"})

	fis, err := UploadAllImages(req.MultipartForm.File, dir1, dir2)
	if err != nil {
		t.Errorf("UploadAllImages(): Returned an error! [%s]", err)
	}

	for _, fi := range fis {
		if _, err := os.Stat(dir1 + string(os.PathSeparator) + fi.Name); err != nil { // check if the image exists
			t.Errorf("UploadAllImages(): Image does not exist! [%s]", err)
		}
		if _, err := os.Stat(dir2 + string(os.PathSeparator) + fi.Name); err != nil { // check if the thumbnail exists
			t.Errorf("UploadAllImages(): Thumbnail does not exist! [%s]", err)
		}
	}
}
