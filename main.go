package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
"os"
	"github.com/julienschmidt/httprouter"
)

func homepageHanlder(w http.ResponseWriter, r *http.Request, p httprouter.Params) *HttpError {
	fmt.Fprintf(w, "<html>Go Server is working fine</html>")
	return nil
}

var Log *log.Logger

func init(){

var logpath = "log_file.log"

var file, err1 = os.Create(logpath)

   if err1 != nil {
      panic(err1)
   }
      Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
      Log.Println("LogFile : " + logpath)


}

func main() {
	InitDb()

	router := httprouter.New()

	//router.GET("/", ErrorHandler(homepageHanlder))

	router.POST("/signup",ErrorHandler(signUpHandler))
	router.POST("/create_post",ErrorHandler(createPostHandler))
	router.GET("/posts",ErrorHandler(getPostsHandler))

  router.GET("/loaderio-fd91954a77340e382e3b566eb7018c9f.txt", ErrorHandler(LoaderIoFileHandler))

	http.ListenAndServe(":8080", router)

}

func PanicHandler(funcName string) {
	if r := recover(); r != nil {
		//TODO:generate a issue number and log the details here. maybe take w, r from handler and log params from there also for debugging later.
		log.Println("Recovered from panic in "+funcName+"\nPanic: ", r)
	}
}

type AppHandler func(http.ResponseWriter, *http.Request, httprouter.Params) *HttpError

func ErrorHandler(a AppHandler) httprouter.Handle { //respond to user with error in requested encoding
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer PanicHandler(runtime.FuncForPC(reflect.ValueOf(a).Pointer()).Name())
		if e := a(w, r, p); e != nil {

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(e.Code)

			json.NewEncoder(w).Encode(struct {
				*HttpError
				IsError bool `json:"error"`
			}{
				HttpError: e,
				IsError:   true,
			})
		}
	}
}

func PrintSuccessJson(w http.ResponseWriter, j []byte) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,Token")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if bytes.Equal(j, []byte("null")) {
		fmt.Fprintf(w, "%s", "{}")
	} else {
		fmt.Fprintf(w, "%s", j)
	}
}

func LoaderIoFileHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) *HttpError{
http.ServeFile(w, r, "loaderio-fd91954a77340e382e3b566eb7018c9f.txt")
return nil
}
