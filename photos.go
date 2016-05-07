package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func uploadPhoto(src multipart.File, hdr *multipart.FileHeader, c *http.Cookie) *http.Cookie {
	defer src.Close()
	// get filename with extension and store it in fName
	fName := getSha(src) + filepath.Ext(hdr.Filename)
	//Get root folder
	wd, _ := os.Getwd()
	//Get model for its username to store in a folder of the same name
	m := Model(c)
	//Join the path together
	path := filepath.Join(wd, "assets", "imgs", m.Name, fName)
	//Create the path
	dst, _ := os.Create(path)
	defer dst.Close()
	//I have no idea what this does but I need it
	src.Seek(0,0)
	//Copy the file in the path
	io.Copy(dst, src)
	return addPhoto("/imgs/" + m.Name + "/" + fName, c)
}

//Stores the file path inside the Model
func addPhoto(fName string, c *http.Cookie) *http.Cookie {
	//Get Model for m.Pictures
	m := Model(c)
	//Store the file path in string slice, m.Picures
	m.Pictures = append(m.Pictures, fName)
	//Get id from old Model and update the cookie with updated model
	xs := strings.Split(c.Value, "|")
	id := xs[0]
	cookie := currentVisitor(m, id)
	return cookie
}

//encryption stuff
func getSha(src multipart.File) string{
	h := sha1.New()
	io.Copy(h, src)
	return fmt.Sprintf("%x", h.Sum(nil))
}
