package fileupload

import (
	"bytes"
	"github.com/google/uuid" // external dependency
	"github.com/pkg/errors"  // external dependency
	"gopkg.in/h2non/bimg.v1" // external dependency
	"io"
	"math"
	"mime/multipart"
	"os"
)

/*
	Uses the Bimg library to save a thumbnail of a large image
*/
func saveThumbnail(buffer []byte, path, name string, width, height int) error {
	widthRatio := float64(width) / float64(height)
	thumbHeight := 75
	thumbWidth := int(math.Floor((float64(thumbHeight) * widthRatio) + 0.5)) // round the float
	options := bimg.Options{Quality: 90, Type: bimg.JPEG, Width: thumbWidth, Height: thumbHeight}

	newImage, err := bimg.NewImage(buffer).Process(options) // do image parsing
	if err != nil {
		return errors.Wrap(err, "saveThumbnail()")
	}
	err = bimg.Write(path+string(os.PathSeparator)+name, newImage) // save image to a file
	if err != nil {
		return errors.Wrap(err, "saveThumbnail()")
	}

	return nil
}

/*
	Uses the Bimg library to save a copy of an image
	Returns data about the image (filesize, width, height, name, ...)
*/
func saveImage(buffer []byte, oldName, newName, directory string) (*FileInfo, error) {
	options := bimg.Options{Quality: 90, Type: bimg.JPEG}

	img := bimg.NewImage(buffer)
	imgSize, err := img.Size()
	if err != nil {
		return nil, errors.Wrap(err, "saveImage()")
	}

	newImage, err := img.Process(options) // do image parsing
	if err != nil {
		return nil, errors.Wrap(err, "saveImage()")
	}
	err = bimg.Write(directory+string(os.PathSeparator)+newName, newImage) // save image to a file
	if err != nil {
		return nil, errors.Wrap(err, "saveImage()")
	}

	var size int64
	if fi, err := os.Stat(directory + string(os.PathSeparator) + newName); err != nil {
		return nil, errors.Wrap(err, "saveImage()")
	} else {
		size = fi.Size()
	}

	return &FileInfo{Name: newName, OriginalName: oldName, Size: size, MimeType: "image/jpeg", IsImage: true, Directory: directory, Width: imgSize.Width, Height: imgSize.Height}, nil
}

func UploadImageWithThumbnail(header *multipart.FileHeader, imageDir string, thumbnailDir string) (*FileInfo, error) {
	// check if the directories exist
	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		return nil, ErrDirectoryDoesNotExist
	}
	if _, err := os.Stat(thumbnailDir); os.IsNotExist(err) {
		return nil, ErrDirectoryDoesNotExist
	}

	file, err := header.Open()
	if err != nil {
		return nil, errors.Wrap(err, "UploadImageWithThumbnail()")
	}
	defer file.Close()

	// is it an image?
	mimetype, err := getMimeType(file)
	if err != nil {
		return nil, errors.Wrap(err, "UploadImageWithThumbnail()")
	}
	if isFileImage(mimetype) == false {
		return nil, ErrNotImageType
	}

	// copy file to a buffer
	buffer := &bytes.Buffer{}
	if _, err := io.Copy(buffer, file); err != nil {
		return nil, errors.Wrap(err, "UploadImageWithThumbnail()")
	}

	// UUIDv4 is used to avoid name conflicts (filename already exists errors)
	uuidStr := uuid.New().String() + ".jpg"

	fi, err := saveImage(buffer.Bytes(), header.Filename, uuidStr, imageDir) // re-save the uploaded image
	if err != nil {
		return nil, errors.Wrap(err, "UploadImageWithThumbnail()")
	}
	if err := saveThumbnail(buffer.Bytes(), thumbnailDir, uuidStr, fi.Width, fi.Height); err != nil { // create a thumbnail
		return nil, errors.Wrap(err, "UploadImageWithThumbnail()")
	}

	return fi, nil
}

/*
	Upload and resave fullsize images and create thumbnails. Fullsize images and thumbnails are saved to two separate directories.
	Return data about each file, useful for Javascript upload tools.
*/
func UploadAllImages(files map[string][]*multipart.FileHeader, imageDir, thumbnailDir string) ([]FileInfo, error) {

	// check if the directories exist
	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		return nil, ErrDirectoryDoesNotExist
	}
	if _, err := os.Stat(thumbnailDir); os.IsNotExist(err) {
		return nil, ErrDirectoryDoesNotExist
	}

	var slice []FileInfo

	for _, fieldSlice := range files {
		for _, header := range fieldSlice {
			fi, err := UploadImageWithThumbnail(header, imageDir, thumbnailDir)
			if err != nil {
				return nil, errors.Wrap(err, "UploadAllImages()")
			}
			slice = append(slice, *fi)
		}
	}

	return slice, nil
}
