package fileUpload

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"phyEcom.com/auth"
)

func remove(fileNames *[]string) {
	for _, name := range *fileNames {
		os.Remove(name)
	}
}

func getFileName(path string) string {
	arr := strings.Split(path, "\\")
	return arr[len(arr)-1]
}

type UploadData struct {
	FormValue bool
	Req       *http.Request
	Res       http.ResponseWriter
	MaxSize   int64
	FileName  string
	Directory string
	FileMatch []string
	UsePath   string
}

func (u *UploadData) generator(done <-chan any, form *multipart.Form) <-chan *multipart.FileHeader {
	formChan := make(chan *multipart.FileHeader)
	go func() {
		defer close(formChan)
		files, ok := form.File[u.FileName]
		if !ok {
			fmt.Println("File doesn't have the given field name")
			return
		}
		for _, file := range files {
			select {
			case <-done:
				return
			case formChan <- file:
			}
		}

	}()
	return formChan
}

func (u *UploadData) typeCheck(done <-chan any, file <-chan *multipart.FileHeader, names *[]string) <-chan *multipart.File {
	fileChan := make(chan *multipart.File)
	go func() {
		defer close(fileChan)
		for f := range file {
			fi, err := f.Open()
			if err != nil {
				remove(names)
				fmt.Printf("Err: error accessing file %v", err)
				return
			}
			select {
			case <-done:
				return
			case fileChan <- &fi:
			}
		}

	}()
	return fileChan

}
func (u *UploadData) typeFit(buf *bufio.Reader) bool {
	matched := false
	sniff, _ := buf.Peek(512)

	fileType := http.DetectContentType(sniff)

	for _, match := range u.FileMatch {
		if fileType == match {
			matched = true
		}
	}
	return matched
}

func (u *UploadData) addFile(done <-chan any, file <-chan *multipart.File, names *[]string) <-chan string {
	filePath := make(chan string)
	go func() {
		defer close(filePath)
		for f := range file {

			buf := bufio.NewReader(*f)
			if !u.typeFit(buf) {
				remove(names)
				fmt.Println("Err: File type doesn't match")
				return
			}
			name := fmt.Sprintf(`%s-*.png`, uuid.NewString())
			tempfile, err1 := os.CreateTemp(u.Directory, name)

			if err1 != nil {
				remove(names)
				fmt.Println("error: error creating temporary file in the said directory")
				return
			}
			defer tempfile.Close()
			lmt := io.MultiReader(buf, io.LimitReader(*f, u.MaxSize-511))
			written, copyErr := io.Copy(tempfile, lmt)

			*names = append(*names, auth.NormalizePath(fmt.Sprintf(`%s/%s`, u.UsePath, getFileName(tempfile.Name()))))
			if copyErr != nil {
				remove(names)
				fmt.Println("Err: Could not Copy")
				return
			}

			if written > u.MaxSize {
				remove(names)
				fmt.Println("error: Larger then expected")
				return
			}

			select {
			case <-done:
				return
			case filePath <- tempfile.Name():
			}

		}

	}()
	return filePath

}

func (u *UploadData) NewUpload() (map[string][]string, []string, error) {
	names := []string{}
	done := make(chan any)
	defer close(done)
	u.Req.Body = http.MaxBytesReader(u.Res, u.Req.Body, u.MaxSize)
	rd, err := u.Req.MultipartReader()
	if err != nil {
		return nil, nil, errors.New("err: error in multipart")
	}
	form, err := rd.ReadForm(u.MaxSize)
	if err != nil {
		log.Println(err)
	}
	formFile := u.generator(done, form)
	file := u.typeCheck(done, formFile, &names)
	pathSavedTo := u.addFile(done, file, &names)
	for path := range pathSavedTo {
		fmt.Println(path)
	}
	if u.FormValue {

		return form.Value, names, nil
	}
	return nil, names, nil
}
