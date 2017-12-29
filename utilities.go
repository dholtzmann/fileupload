package fileupload

import (
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/pkg/errors" // external dependency
)

func isFileImage(mimetype string) bool {
	mimetype = strings.ToLower(mimetype) // make it case insensitive
	var types []string = []string{"image/jpeg", "image/png", "image/gif"}
	return inSlice(types, mimetype)
}

func getMimeType(file multipart.File) (string, error) {
	buffer := make([]byte, 512)

	if _, err := file.Seek(0, 0); err != nil { // set position to the start of the file
		return "", errors.Wrap(err, "getMimeType()")
	}

	if n, err := file.Read(buffer); n <= 0 && err != nil { // Copy bytes into a buffer
		return "", errors.Wrap(err, "getMimeType()")
	}

	if _, err := file.Seek(0, 0); err != nil { // reset position to start of the file
		return "", errors.Wrap(err, "getMimeType()")
	}

	// http.DetectContentType() `always returns a valid MIME type: if it cannot determine a more specific one, it returns "application/octet-stream"`
	return strings.TrimSpace(strings.ToLower(http.DetectContentType(buffer))), nil
}

// Is some value in the slice?
func inSlice(slice []string, val string) bool {
	for _, j := range slice {
		if j == val {
			return true
		}
	}
	return false
}

func getFileExtension(name string) string {
	if pos := strings.LastIndex(name, "."); pos != -1 {
		return strings.ToLower(name[pos+1:])
	}
	return ""
}
