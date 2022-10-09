package server

import (
	"encoding/json"
	"fmt"

	app "github.com/kormiltsev/url-testtask/internal/internal/app"
	st "github.com/kormiltsev/url-testtask/internal/storage"

	//"github.com/gorilla/mux"
	"log"
	"net/http"
)

type responseToClient struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    map[string]string `json:"data"`
}

// operates requests
func defaultFunc(w http.ResponseWriter, r *http.Request) {
	var ok = "success"
	fmt.Println("client connect success ", r.RemoteAddr)
	fmt.Println(r.Method, r.RequestURI)
	data := make(map[string]string)
	//check header for content type
	ct, k := r.Header["Content-Type"]
	if k {
		// check for json
		if ct[0] == "application/json" {
			buf := make([]byte, 2048)
			n, _ := r.Body.Read(buf)
			json.Unmarshal(buf[:n], &data)
		}
	} else {
		// simple form
		r.ParseForm()
		for w, v := range r.Form {
			data[w] = v[0]
		}
	}
	var NewLink = st.Request{}

	// form/json request
	if data["url"] != "" || data["short_url"] != "" {
		NewLink.Url = data["url"]
		NewLink.Surl = data["short_url"]
		NewLink, ok = MetodSwitcher(NewLink, r)

		// sending json
		data["url"] = NewLink.Url
		data["short_url"] = NewLink.Surl
		m := responseToClient{200, ok, data}
		mjson, e := json.Marshal(m)
		if e != nil {
			fmt.Println(e)
		}
		fmt.Fprintf(w, "%v\n", string(mjson))

	}
	// in case "get"&"post" in param
	if data["get"] != "" {
		NewLink.Url = ""
		NewLink.Surl = data["get"]
		NewLink, ok = app.Get(NewLink)
		fmt.Fprintf(w, NewLink.Url)
	}
	if data["post"] != "" {
		NewLink.Url = data["post"]
		NewLink.Surl = ""
		NewLink, ok = app.Post(NewLink)
		fmt.Fprintf(w, NewLink.Surl)
	}
}

// switcher
func MetodSwitcher(req st.Request, r *http.Request) (st.Request, string) {
	if r.Method == "GET" {
		return app.Get(req)
	}
	if r.Method == "POST" {
		return app.Post(req)
	}
	return req, "unknown Method"
}

// start router and listener
func StartServe() {
	//default
	http.HandleFunc("/", defaultFunc)

	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
