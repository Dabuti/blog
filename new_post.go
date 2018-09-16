package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func newPost(w http.ResponseWriter, r *http.Request) {
	var Buf bytes.Buffer
	maxUploadSize := int64(4096)

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		fmt.Println(w, "FILE_TOO_BIG", http.StatusBadRequest)
		return
	}

	fileType := r.PostFormValue("type")
	fmt.Println(fileType)
	file, _, err := r.FormFile("data")
	if err != nil {
		fmt.Println(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}

	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Println(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}

	// Print data contents
	file, _, _ = r.FormFile("data")
	io.Copy(&Buf, file)
	contents := Buf.String()
	fmt.Println(contents)

	Buf.Reset()

	newPath := "./posts/current/id.org"
	newFile, err := os.Create(newPath)
	if err != nil {
		fmt.Println(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()
	if _, err := newFile.Write(fileBytes); err != nil {
		fmt.Println(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("SUCCESS"))
}
