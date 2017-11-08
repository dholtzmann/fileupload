package fileupload

import (
	"github.com/pkg/errors" // external dependency
	"os"
)

/*
	Loop through a directory and return a slice of structs with information about all the files in said directory
	includeMimeType - flag, whether or not to include each file's mimetype, the operation takes more processing of the files
*/
func GetDirectoryContentsData(directory string, includeMimeType bool) ([]FileInfo, error) {
	var fis []FileInfo

	if _, err := os.Stat(directory); os.IsNotExist(err) { // check if the directory exists
		return nil, ErrDirectoryDoesNotExist
	}

	var file *os.File
	var err error

	if file, err = os.Open(directory); err != nil { // open the directory
		return nil, errors.Wrap(err, "GetDirectoryContentsData()")
	}

	if includeMimeType { // will have to open each file to examine first bytes
		fileSlice, _ := file.Readdir(0)
		for _, fi := range fileSlice {
			if fi.IsDir() == false {
				file, err := os.Open(directory + string(os.PathSeparator) + fi.Name())
				if err != nil { // check for error
					return nil, errors.Wrap(err, "GetDirectoryContentsData()")
				}
				defer file.Close()
				mimetype, err := getMimeType(file)
				if err != nil { // check for error
					return nil, errors.Wrap(err, "GetDirectoryContentsData()")
				}

				fis = append(fis, FileInfo{Name: fi.Name(), Size: fi.Size(), Directory: directory, MimeType: mimetype, IsImage: isFileImage(mimetype)})
			}
		}

		return fis, nil
	}

	// includeMimeType flag is false
	fileSlice, _ := file.Readdir(0)
	for _, fi := range fileSlice {
		if fi.IsDir() == false {
			fis = append(fis, FileInfo{Name: fi.Name(), Size: fi.Size(), Directory: directory})
		}
	}

	return fis, nil
}
