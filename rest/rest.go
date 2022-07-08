package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Jagiello95/common/slice"
	"github.com/gabriel-vasile/mimetype"
)

func NewRestUtil() *RestUtil {
	return &RestUtil{}
}

type RestUtil struct {
	MaxFileSize int
}

// JSONResponse is the type used for sending JSON around
type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// ReadJSON tries to read the body of a request and converts it into JSON
func (u *RestUtil) ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // one megabyte
	if u.MaxFileSize > 0 {
		maxBytes = u.MaxFileSize
	}

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single json value")
	}

	return nil
}

// WriteJSON takes a response status code and arbitrary data and writes a json response to the client
func (u *RestUtil) WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// ErrorJSON takes an error, and optionally a response status code, and generates and sends
// a json error response
func (u *RestUtil) ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JSONResponse
	payload.Error = true
	payload.Message = err.Error()

	return u.WriteJSON(w, statusCode, payload)
}

// PushJSONToRemote posts arbitrary json to some url, and returns error,
// if any, as well as the response status code
func (u *RestUtil) PushJSONToRemote(client *http.Client, uri string, data any) (int, error) {
	// create json we'll send
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return 0, err
	}

	// build the request and set header
	request, err := http.NewRequest("POST", uri, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}
	request.Header.Set("Content-Type", "application/json")

	// call the uri
	response, err := client.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	return response.StatusCode, nil
}

func (u *RestUtil) GetFileToUpload(r *http.Request, fieldName string) (string, error) {
	var maxUploadSize int64
	sliceUtil := slice.NewSliceUtil()
	if max, err := strconv.Atoi(os.Getenv("MAX_UPLOAD_SIZE")); err != nil {
		maxUploadSize = 10 << 20
	} else {
		maxUploadSize = int64(max)
	}

	_ = r.ParseMultipartForm(maxUploadSize)
	fmt.Println("def")

	file, header, err := r.FormFile(fieldName)
	fmt.Println("ghj")

	if err != nil {
		fmt.Println("err1")
		return "", err
	}

	defer file.Close()
	// look at the first 5000 bytes and see the mime type
	mimeType, err := mimetype.DetectReader(file)
	fmt.Println("jkl")

	if err != nil {
		return "", err
	}

	fmt.Println("mimeType", mimeType)

	//go back to start of file
	_, err = file.Seek(0, 0)

	if err != nil {
		return "", err
	}
	exploded := strings.Split(os.Getenv("ALLOWED_FILETYPES"), ",")
	var mimeTypes []string

	mimeTypes = append(mimeTypes, exploded...)

	if !sliceUtil.InSlice(mimeTypes, mimeType.String()) {
		return "", errors.New("invalid file type uploaded")
	}

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)

	dst, err := os.Create(fmt.Sprintf("./tmp/%s", header.Filename))
	if err != nil {
		fmt.Println("err2")
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("./tmp/%s", header.Filename), nil
}
