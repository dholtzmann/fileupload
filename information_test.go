package fileupload

import (
	"io/ioutil"
	"os"
	"testing"
)

func Test_GetDirectoryContentsData(t *testing.T) {
	list := []FileInfo{
		FileInfo{OriginalName: "gopher.png", Size: 17668, Width: 250, Height: 340, MimeType: "image/png", IsImage: true},
		FileInfo{OriginalName: "blue.jpg", Size: 10011, Width: 300, Height: 163, MimeType: "image/jpeg", IsImage: true},
		FileInfo{OriginalName: "car.tar.gz", Size: 10147, MimeType: "application/x-gzip", IsImage: false},
	}
	tempDir := "testing-filevalidator"
	dir, err := ioutil.TempDir("", tempDir) // make a temp directory
	if err != nil {
		t.Fatalf("ioutil.TempDir: %s", err)
	}
	defer os.RemoveAll(dir) // delete the temp directory

	req := setupRequestMultipartForm(&testFile{gopherPNG, "imageupload", "gopher.png"}, &testFile{blueJPG, "imageupload", "blue.jpg"}, &testFile{carTARGZ, "imageupload", "car.tar.gz"})

	// upload files to a temporary directory which will be checked
	var fis []FileInfo
	if fis, err = UploadAllFiles(req.MultipartForm.File, dir, true); err != nil {
		t.Errorf("UploadAllFiles(): Returned an error! [%s]", err)
	}

	// verify that the files exist
	for _, fi := range fis {
		if _, err := os.Stat(fi.Directory + string(os.PathSeparator) + fi.Name); err != nil { // check if the file exists
			t.Errorf("UploadAllFiles(): File failed to upload! [%s]", err)
		}
	}

	fis, err = GetDirectoryContentsData(dir, true) // test with includeMimeType flag
	if err != nil {
		t.Errorf("GetDirectoryContentsData(): Returned an error! [%s] Temporary directory: [%s]", err, dir)
	}

	var foundResultsSize, foundMimeTypes, foundIsImage int
	for _, fi := range fis {
		for _, l := range list { // the returned results [GetDirectoryContentsData()] will have an unknown order due to UUID strings alphabetical order, so must loop over entire list
			if l.Size == fi.Size {
				foundResultsSize += 1
			}
			if l.MimeType == fi.MimeType {
				foundMimeTypes += 1

				if l.IsImage == fi.IsImage { // statement here avoids the counter being increased several times in the loop
					foundIsImage += 1
				}
			}
		}
	}

	if foundResultsSize != len(list) {
		t.Errorf("GetDirectoryContentsData(): Does not match prepared testing list! Size variable: [%d]", foundResultsSize)
	}
	if foundMimeTypes != len(list) {
		t.Errorf("GetDirectoryContentsData(): Does not match prepared testing list! Mimetype variable: [%d]", foundMimeTypes)
	}
	if foundIsImage != len(list) {
		t.Errorf("GetDirectoryContentsData(): Does not match prepared testing list! isImage variable: [%d]", foundIsImage)
	}

	fis, err = GetDirectoryContentsData(dir, false) // test without includeMimeType flag
	if err != nil {
		t.Errorf("GetDirectoryContentsData(): Returned an error! [%s] Temporary directory: [%s]", err, dir)
	}

	foundResultsSize = 0
	for _, fi := range fis {
		for _, l := range list { // the returned results [GetDirectoryContentsData()] will have an unknown order due to UUID strings alphabetical order, so must loop over entire list
			if l.Size == fi.Size {
				foundResultsSize += 1
			}
		}
	}

	if foundResultsSize != len(list) {
		t.Errorf("GetDirectoryContentsData(): Does not match prepared testing list! Size variable: [%d]", foundResultsSize)
	}
}
