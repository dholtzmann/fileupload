package fileupload

import (
	"encoding/json"
	"github.com/google/uuid" // external dependency
	"github.com/pkg/errors"  // external dependency
	"io"
	"mime/multipart"
	"os"
)

var ErrDirectoryDoesNotExist = errors.New("The provided directory does not exist!")
var ErrNoMatchingMimeType = errors.New("No matching mimetype found!")
var ErrNotImageType = errors.New("This file is not an image! Bad mimetype!")

// this is provided for convenience, javascript uploads can parse this JSON
// fields must be exported for encoding/json Marshal() to work!
type FileInfo struct {
	Name         string `json:"name"`
	OriginalName string `json:"originalName,omitempty"`
	Size         int64  `json:"size"`
	IsImage      bool   `json:"-"`
	Directory    string `json:"path"`
	MimeType     string `json:"mimeType"`
	Width        int    `json:"width,omitempty"`
	Height       int    `json:"height,omitempty"`

	Url           string `json:"url,omitempty"`
	ThumbnailUrl  string `json:"thumbnailUrl,omitempty"`
	DeleteUrl     string `json:"deleteUrl,omitempty"`
	DeleteNoJSUrl string `json:"-"`
	DeleteMethod  string `json:"deleteMethod,omitempty"`
	Error         string `json:"error,omitempty"`
}

type Category struct {
	mimeTypes []string
	directory string
}

func CatSlice(arg ...*Category) []Category {
	var cat []Category
	for _, c := range arg {
		cat = append(cat, *c)
	}
	return cat
}

func (fi *FileInfo) Json() (string, error) {
	b, err := json.Marshal(fi)
	if err != nil {
		return "", errors.Wrap(err, "FileInfo.Json()")
	}

	return string(b), nil
}

func SliceJSON(fis []FileInfo) (string, error) {
	b, err := json.Marshal(fis)
	if err != nil {
		return "", errors.Wrap(err, "SliceJSON()")
	}

	return string(b), nil
}

// copy an uploaded file to a directory
func UploadFile(header *multipart.FileHeader, directory string, includeOldExtension bool) (*FileInfo, error) {
	if _, err := os.Stat(directory); os.IsNotExist(err) { // does the directory exist?
		return nil, ErrDirectoryDoesNotExist
	}

	var newName, mimetype string
	newName = uuid.New().String() //UUIDv4
	if includeOldExtension {
		newName += "." + getFileExtension(header.Filename)
	}

	file, err := header.Open()
	if err != nil {
		return nil, errors.Wrap(err, "UploadFile()")
	}
	defer file.Close()

	if mimetype, err = getMimeType(file); err != nil {
		return nil, errors.Wrap(err, "UploadFile()")
	}

	f, err := os.OpenFile(directory+string(os.PathSeparator)+newName, os.O_WRONLY|os.O_CREATE, 0644) // create a file
	if err != nil {
		return nil, errors.Wrapf(err, "UploadFile() Filename[%s]", directory+newName)
	}
	defer f.Close() // close created file

	size, err := io.Copy(f, file) // copy the uploaded file to the created file
	if err != nil {
		return nil, errors.Wrap(err, "UploadFile()")
	}

	return &FileInfo{Name: newName, OriginalName: header.Filename, Size: size, IsImage: isFileImage(mimetype), Directory: directory, MimeType: mimetype}, nil
}

// copy the uploaded file to a created file
func copyUploadedFile(directory, newName, mimetype, oldName string, file multipart.File) (*FileInfo, error) {
	f, err := os.OpenFile(directory+string(os.PathSeparator)+newName, os.O_WRONLY|os.O_CREATE, 0644) // create a file
	if err != nil {
		return nil, errors.Wrapf(err, "copyUploadedFile() Filename[%s]", directory+newName)
	}
	defer f.Close() // close created file

	size, err := io.Copy(f, file) // copy the uploaded file to the created file
	if err != nil {
		return nil, errors.Wrap(err, "copyUploadedFile()")
	}

	return &FileInfo{Name: newName, OriginalName: oldName, Size: size, IsImage: isFileImage(mimetype), Directory: directory, MimeType: mimetype}, nil
}

/*
	Copy an uploaded file to a directory, sorted by Category, see the struct declaration

	Use []string{"*"} in Category for matching leftover types/backup

	Example:
	uploadFileByCategory(header, []Category{ Category{ []string{"*"}, "uploads-directory"} })
*/
func UploadFileByCategory(header *multipart.FileHeader, list []Category, includeOldExtension bool) (*FileInfo, error) {
	var newName, mimetype, leftoverTypesDirectory string

	newName = uuid.New().String() //UUIDv4
	if includeOldExtension {
		newName += "." + getFileExtension(header.Filename)
	}

	file, err := header.Open()
	if err != nil {
		return nil, errors.Wrap(err, "UploadFileByCategory()")
	}
	defer file.Close()

	if mimetype, err = getMimeType(file); err != nil {
		return nil, errors.Wrap(err, "UploadFileByCategory()")
	}

	for _, l := range list { // loop through the Category list
		if _, err := os.Stat(l.directory); os.IsNotExist(err) { // all supplied directories must exist
			return nil, ErrDirectoryDoesNotExist
		}

		if len(l.mimeTypes) == 1 && l.mimeTypes[0] == "*" { // look for the left-over/backup type
			leftoverTypesDirectory = l.directory
		}

		if inSlice(l.mimeTypes, mimetype) { // found a mimetype in a Category
			return copyUploadedFile(l.directory, newName, mimetype, header.Filename, file)
		}
	}

	if len(leftoverTypesDirectory) > 0 { // no mimetype matched, use the backup directory
		return copyUploadedFile(leftoverTypesDirectory, newName, mimetype, header.Filename, file)
	}

	return nil, ErrNoMatchingMimeType
}

func UploadAllFiles(files map[string][]*multipart.FileHeader, directory string, includeOldExtension bool) ([]FileInfo, error) {
	var slice []FileInfo

	for _, fieldSlice := range files {
		for _, header := range fieldSlice {
			fi, err := UploadFile(header, directory, includeOldExtension)
			if err != nil {
				return nil, errors.Wrap(err, "UploadAllFiles()")
			}
			slice = append(slice, *fi)
		}
	}

	return slice, nil
}

func UploadAllFilesByCategory(files map[string][]*multipart.FileHeader, list []Category, includeOldExtension bool) ([]FileInfo, error) {
	var slice []FileInfo

	for _, fieldSlice := range files {
		for _, header := range fieldSlice {
			fi, err := UploadFileByCategory(header, list, includeOldExtension)
			if err != nil {
				return nil, errors.Wrap(err, "UploadAllFilesByCategory()")
			}
			slice = append(slice, *fi)
		}
	}

	return slice, nil
}
