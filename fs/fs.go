package fs

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/Jagiello95/common/random"

	"github.com/gabriel-vasile/mimetype"
)

func NewFileSystemUtil() FileSystemUtil {
	return FileSystemUtil{}
}

type FileSystemUtil struct{}

// CreateDirIfNotExist creates a directory, and all necessary parent directories, if it does not exist.
func (f *FileSystemUtil) CreateDirIfNotExist(path string) error {
	const mode = 0755
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, mode)
		if err != nil {
			return err
		}
	}
	return nil
}

// DownloadStaticFile downloads a file, and tries to force the browser to avoid displaying it in
// the browser window by setting content-disposition. It also allows specification of the display name.
func (f *FileSystemUtil) DownloadStaticFile(w http.ResponseWriter, r *http.Request, p, file, displayName string) {
	fp := path.Join(p, file)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", displayName))

	http.ServeFile(w, r, fp)
}

// UploadOneFile uploads one file to a specified directory, and gives it a random name.
// It returns the newly named file, the original file name, and potentially an error.
func (f *FileSystemUtil) UploadOneFile(r *http.Request, uploadDir string) (string, string, error) {
	u := random.NewRandomUtil()
	// parse the form so we have access to the file
	err := r.ParseMultipartForm(1024 * 1024 * 1024)
	if err != nil {
		return "", "", err
	}

	var filename, fileNameDisplay string
	for _, fHeaders := range r.MultipartForm.File {
		for _, hdr := range fHeaders {
			infile, err := hdr.Open()
			if err != nil {
				return "", "", err
			}
			defer infile.Close()

			ext, err := mimetype.DetectReader(infile)
			if err != nil {
				fmt.Println(err)
				return "", "", err
			}

			_, err = infile.Seek(0, 0)
			if err != nil {
				fmt.Println(err)
				return "", "", err
			}

			filename = u.RandomString(25) + ext.Extension()
			fileNameDisplay = hdr.Filename

			var outfile *os.File
			defer outfile.Close()

			if outfile, err = os.Create(uploadDir + filename); nil != err {
				fmt.Println(err)
			} else {
				_, err := io.Copy(outfile, infile)
				if err != nil {
					fmt.Println(err)
					return "", "", err
				}
			}
		}

	}
	return filename, fileNameDisplay, nil
}
