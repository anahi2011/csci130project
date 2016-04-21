package csci130project

import (
	//"golang.org/x/net/context"
	"net/http"
	//"google.golang.org/cloud/storage"
	//"io"
	"html/template"
)

const gcsBucket = "csci-130group.appspot.com"


//we could use this to pass it to our functions
//that way we don't have to have too many parameters
//NOT SURE IF WE WILL USE THIS THOUGH
//type demo struct{
//	ctx	context.Context
//	w	http.ResponseWriter
//	bucket	*storage.BucketHandle
//	client 	*storage.Client
//}

var tpl* template.Template

func init(){
	tpl, _ = template.ParseGlob("templates/")
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets/"))))
	http.HandleFunc("/", homePage)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080",nil)
}

func homePage(res http.ResponseWriter, req* http.Request){
	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}
	tpl.ExecuteTemplate(res, "home.html", nil)
}

//func putFile(ctx context.Context, name string, rdr io.Reader) error {
//
//	client, err := storage.NewClient(ctx)
//	if err != nil {
//		return err
//	}
//	defer client.Close()
//
//	writer := client.Bucket(gcsBucket).Object(name).NewWriter(ctx)
//
//	io.Copy(writer, rdr)
//	// check for errors on io.Copy in production code!
//	return writer.Close()
//}
//
//func getFile(ctx context.Context, name string) (io.ReadCloser, error) {
//	client, err := storage.NewClient(ctx)
//	if err != nil {
//		return nil, err
//	}
//	defer client.Close()
//
//	return client.Bucket(gcsBucket).Object(name).NewReader(ctx)
//}
//
//func getFileLink(ctx context.Context, name string) (string, error) {
//	client, err := storage.NewClient(ctx)
//	if err != nil {
//		return "", err
//	}
//	defer client.Close()
//
//	attrs, err := client.Bucket(gcsBucket).Object(name).Attrs(ctx)
//	if err != nil {
//		return "", err
//	}
//	return attrs.MediaLink, nil
//}