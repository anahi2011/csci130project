package myapp

import (
	"crypto/sha1"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"google.golang.org/cloud/storage"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine"
	"golang.org/x/net/context"
)

// I added a parameter.. The request
func uploadPhoto(src multipart.File, hdr *multipart.FileHeader, c *http.Cookie, req *http.Request) *http.Cookie {


	//we gotta close our file but we wait till we
	//are done working with it
	defer src.Close()

	//grabbing context for error checking
	ctx := appengine.NewContext(req)


	//grabbing the extension of the file uploaded
	ext := hdr.Filename[strings.LastIndex(hdr.Filename, ".")+1:]

	//if we wanna log the extension uncomment the next line
	//log.Infof(ctx, "FILE EXTENSION: %s", ext)

	//checking the extension type. we only allow .txt , .png , .jpeg , .jpg
	switch ext {
	case "jpg", "jpeg", "png", "txt":
		//log.Infof(ctx, "GOOD FILE EXTENSION: %s", ext)
	default:
		log.Errorf(ctx, "We do not allow files of type %s. We only allow jpg, jpeg, png, txt extensions.", ext)
		return c
	}

	//grabbing our object data from our cookie cause matt sucks at making things modular
	m := Model(c)

	var fName string

	//making our filenames easier to work with and easy to query
	if(ext == "jpeg" || ext == "jpg" || ext == "png"){
		fName = m.Name + "/image/"
	}else{
		fName = m.Name + "/text/"
	}

	//just setting up a basic file structure
	//bucket/ userName/encryptedFilename.ext
	fName = fName + getSha(src) + "." + ext

	//not sure why i would need this but todds code has it when he uses
	//the multipart file so im gunna include it cause i used the multipart
	//file in the previous lines
	//something about the offset for the next read or write operations
	//so we set it to zero so we read or write to it starting from the start
	src.Seek(0, 0)


	//creating a new client from our context
	client, err := storage.NewClient(ctx)
	if err != nil{
		log.Errorf(ctx, "Error in main client err")
		return c
	}
	defer client.Close()

	//making a writer for our specific bucket and file
	writer := client.Bucket(gcsBucket).Object(fName).NewWriter(ctx)


	//making the file public
	writer.ACL = []storage.ACLRule{
		{storage.AllUsers, storage.RoleReader},
	}

	//setting the type of the file png/jpg/txt
	if(ext == "jpeg" || ext == "jpg") {
		writer.ContentType = "image/jpeg"
	}else if(ext == "png"){
		writer.ContentType = "image/png"
	}else{
		writer.ContentType = "text/plain"
	}

	//writing the file to the gcs bucket
	io.Copy(writer, src)

	err = writer.Close()
	if(err != nil){
		log.Errorf(ctx, "error uploadPhoto writer close", err)
		return c
	}


	return addPhoto(fName, ext, c)
}

//Stores the file path inside the Model
func addPhoto(fName string, ext string, c *http.Cookie) *http.Cookie {
	//Get Model for m.Pictures
	m := Model(c)
	//If the file is an image with jpg or png extension, put it in m.Pictures
	if ext == "jpg" || ext == "png"{
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
