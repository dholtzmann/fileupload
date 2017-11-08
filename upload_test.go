package fileupload

import (
	"io/ioutil"
	"os"
	"testing"
)

func Test_Json(t *testing.T) {
	testStr := `{"name":"thing.jpg","size":1000,"path":"images/","mimeType":"image/jpeg"}`
	var fi FileInfo = FileInfo{Name: "thing.jpg", Size: 1000, Directory: "images/", MimeType: "image/jpeg"}
	js, err := fi.Json()
	if err != nil {
		t.Errorf("FileInfo.Json(): Returned an error! [%s]", err)
	}
	if js != testStr {
		t.Errorf("FileInfo.Json(): JSON string does not match! [%s]", js)
	}
}

func Test_SliceJSON(t *testing.T) {
	testStr := `[{"name":"thing.jpg","size":1000,"path":"images/","mimeType":"image/jpeg"},{"name":"blah.tar.gz","size":100000,"path":"uploads/","mimeType":"application/x-gzip"}]`
	var fis []FileInfo = []FileInfo{
		FileInfo{Name: "thing.jpg", Size: 1000, Directory: "images/", MimeType: "image/jpeg"},
		FileInfo{Name: "blah.tar.gz", Size: 100000, Directory: "uploads/", MimeType: "application/x-gzip"},
	}

	js, err := SliceJSON(fis)
	if err != nil {
		t.Errorf("SliceJSON(): Returned an error! [%s]", err)
	}
	if js != testStr {
		t.Errorf("SliceJSON(): JSON string does not match! [%s]", js)
	}
}

func TestUploadFile(t *testing.T) {
	tempDir := "testing-filevalidator"
	dir, err := ioutil.TempDir("", tempDir) // make a temp directory
	if err != nil {
		t.Fatalf("ioutil.TempDir: %s", err)
	}
	defer os.RemoveAll(dir) // delete the temp directory

	req := setupRequestMultipartForm(&testFile{gopherPNG, "imageupload", "gopher.png"}, &testFile{blueJPG, "imageupload", "blue.jpg"}, &testFile{carTARGZ, "imageupload", "car.tar.gz"})

	for _, fieldSlice := range req.MultipartForm.File { // map[string][]*FileHeader
		for _, header := range fieldSlice {
			var mfi *FileInfo
			if mfi, err = UploadFile(header, dir, true); err != nil {
				t.Errorf("UploadFile(): Returned an error! [%s]", err)
			}

			if _, err := os.Stat(dir + string(os.PathSeparator) + mfi.Name); err != nil { // check if the file exists
				t.Errorf("UploadFile(): File failed to upload! [%s]", err)
			}
		}
	}
}

func TestUploadAllFiles(t *testing.T) {
	tempDir := "testing-filevalidator"
	dir, err := ioutil.TempDir("", tempDir) // make a temp directory
	if err != nil {
		t.Fatalf("ioutil.TempDir: %s", err)
	}
	defer os.RemoveAll(dir) // delete the temp directory

	req := setupRequestMultipartForm(&testFile{gopherPNG, "imageupload", "gopher.png"}, &testFile{blueJPG, "imageupload", "blue.jpg"}, &testFile{carTARGZ, "imageupload", "car.tar.gz"})

	var fis []FileInfo
	if fis, err = UploadAllFiles(req.MultipartForm.File, dir, true); err != nil {
		t.Errorf("UploadAllFiles(): Returned an error! [%s]", err)
	}

	for _, fi := range fis {
		if _, err := os.Stat(fi.Directory + string(os.PathSeparator) + fi.Name); err != nil { // check if the file exists
			t.Errorf("UploadAllFiles(): File failed to upload! [%s]", err)
		}
	}
}

func TestUploadFileByCategory(t *testing.T) {
	tempDir1, tempDir2 := "testing-filevalidator-images", "testing-filevalidator-other-uploads"
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

	list := CatSlice(&Category{[]string{"image/png", "image/gif", "image/jpeg"}, dir1}, &Category{[]string{"application/x-gzip"}, dir2}, &Category{[]string{"*"}, "/tmp"}) // the 3rd entry should not be used

	req := setupRequestMultipartForm(&testFile{gopherPNG, "imageupload", "gopher.png"}, &testFile{blueJPG, "imageupload", "blue.jpg"}, &testFile{carTARGZ, "imageupload", "car.tar.gz"})

	for _, fieldSlice := range req.MultipartForm.File { // map[string][]*FileHeader
		for _, header := range fieldSlice {
			var mfi *FileInfo
			if mfi, err = UploadFileByCategory(header, list, true); err != nil {
				t.Errorf("UploadFileByCategory(): Returned an error! [%s]", err)
			}

			if mfi == nil { // in case of an error accessing mfi.* would cause a panic
				return
			}

			if _, err := os.Stat(mfi.Directory + string(os.PathSeparator) + mfi.Name); err != nil { // check if the file exists
				t.Errorf("UploadFileByCategory(): File failed to upload! [%s]", err)
			}

		}
	}
}

func TestUploadAllFilesByCategory(t *testing.T) {
	tempDir1, tempDir2 := "testing-filevalidator-images", "testing-filevalidator-other-uploads"
	dir1, err := ioutil.TempDir("", tempDir1) // make a two temporary directories
	if err != nil {
		t.Fatalf("ioutil.TempDir: %s", err)
	}
	dir2, err := ioutil.TempDir("", tempDir2)
	if err != nil {
		t.Fatalf("ioutil.TempDir: %s", err)
	}
	defer os.RemoveAll(dir1) // delete the temp directories
	defer os.RemoveAll(dir2)

	list := CatSlice(&Category{[]string{"image/png", "image/gif", "image/jpeg"}, dir1}, &Category{[]string{"application/x-gzip"}, dir2}, &Category{[]string{"*"}, "/tmp"}) // the 3rd entry should not be used
	req := setupRequestMultipartForm(&testFile{gopherPNG, "imageupload", "gopher.png"}, &testFile{blueJPG, "imageupload", "blue.jpg"}, &testFile{carTARGZ, "imageupload", "car.tar.gz"})

	var fis []FileInfo

	// test for an error
	if fis, err = UploadAllFilesByCategory(req.MultipartForm.File, CatSlice(), true); err == nil {
		t.Errorf("UploadAllFilesByCategory(): Should Return an error!")
	}

	if fis, err = UploadAllFilesByCategory(req.MultipartForm.File, list, true); err != nil {
		t.Errorf("UploadAllFilesByCategory(): Returned an error! [%s]", err)
	}

	for _, fi := range fis {
		if _, err := os.Stat(fi.Directory + string(os.PathSeparator) + fi.Name); err != nil { // check if the file exists
			t.Errorf("UploadAllFilesByCategory(): File failed to upload! [%s]", err)
		}
	}
}
