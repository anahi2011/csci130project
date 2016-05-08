package myapp

import (
	"crypto/sha1"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"google.golang.org/cloud/storage"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine"
	"golang.org/x/net/context"
)

// I added a parameter.. The request
func uploadPhoto(src multipart.File, hdr *multipart.FileHeader, c *http.Cookie, req *http.Request) *http.Cookie {
	defer src.Close()
	//Limit kinds of files you can upload
	if(hdr.Filename != ".jpg" ||
	   hdr.Filename != ".png" ||
	   hdr.Filename != ".txt"){
		return c
	}

	m := Model(c)

	var fName string

	//making our filenames easier to work with and easy to query
	if(hdr.Filename == ".jpg" || hdr.Filename == ".png"){
		fName = m.Name + "/photo/"
	}else{
		fName = m.Name + "/text/"
	}

	// get filename with extension and store it in fName.
	//just setting up a basic file structure
	//bucket/ userName/encryptedFilename
	fName += getSha(src) + filepath.Ext(hdr.Filename)


	//grabbing context for error checking
	ctx := appengine.NewContext(req)

	//creating a new client from our context
	client, err := storage.NewClient(ctx)
	if err != nil{
		log.Errorf(ctx, "Error in main client err")
	}
	defer client.Close()

	//grabbing a client fronm our specific bucket
	bucket := client.Bucket(gcsBucket)

	//making a new gcs writer
	writer := bucket.Object(fName).NewWriter(ctx)

	//making the file public
	writer.ACL = []storage.ACLRule{
		{storage.AllUsers, storage.RoleReader},
	}

	//setting the type of the file png/jpg/txt
	writer.ContentType = hdr.Filename

	//writing the file to the gcs bucket
	//NOT SURE IF I'M ALLOWED TO CONVERT OUR FILE TO []byte
	io.Copy(writer, src)

	if(err != nil){
		log.Errorf(ctx, "uploadPhoto: unable to write data to bucket")
		return c
	}

	err = writer.Close()
	if(err != nil){
		log.Errorf(ctx, "uploadPhoto closing writer")
		return c
	}

	return addPhoto(fName, filepath.Ext(hdr.Filename), c)
}

//Stores the file path inside the Model
func addPhoto(fName string, ext string, c *http.Cookie) *http.Cookie {
	//Get Model for m.Pictures
	m := Model(c)
	//If the file is an image with jpg or png extension, put it in m.Pictures
	if ext == ".jpg" || ext == ".png"{
		m.Pictures = append(m.Pictures, fName)
	}
	//Store the file path in string slice, m.Files
	m.Files = append(m.Files, fName)
	//Get id from old Model and update the cookie with updated model
	xs := strings.Split(c.Value, "|")
	id := xs[0]
	cookie := currentVisitor(m, id)
	return cookie
}




//this is some code

func uploadFile(req *http.Request, mpf multipart.File, hdr *multipart.FileHeader) (string, error) {

	ext, err := fileFilter(req, hdr)
	if err != nil {
		return "", err
	}
	name := getSha(mpf) + `.` + ext
	mpf.Seek(0, 0)

	ctx := appengine.NewContext(req)
	return name, putFile(ctx, name, mpf)
}

func fileFilter(req *http.Request, hdr *multipart.FileHeader) (string, error) {

	ext := hdr.Filename[strings.LastIndex(hdr.Filename, ".")+1:]
	ctx := appengine.NewContext(req)
	log.Infof(ctx, "FILE EXTENSION: %s", ext)

	switch ext {
	case "jpg", "jpeg", "txt", "md":
		return ext, nil
	}
	return ext, fmt.Errorf("We do not allow files of type %s. We only allow jpg, jpeg, txt, md extensions.", ext)
}

func getSha(src multipart.File) string {
	h := sha1.New()
	io.Copy(h, src)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func putFile(ctx context.Context, name string, rdr io.Reader) error {

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	writer := client.Bucket(gcsBucket).Object(name).NewWriter(ctx)

	io.Copy(writer, rdr)
	// check for errors on io.Copy in production code!
	return writer.Close()
}
