package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
)

//go:embed resources/index.html
var helloHtml string

//go:embed files
var resources embed.FS

func main() {
	var mux = http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, r.Method)
		fmt.Fprintln(w, r.RequestURI)
		fmt.Fprintln(w, r.URL)
		fmt.Fprintln(w, r.Header)
		fmt.Fprintln(w, r.Host)
		fmt.Fprintln(w, r.Proto)
		fmt.Fprintln(w, r.RemoteAddr)
		fmt.Fprintln(w, r.TLS)
		fmt.Fprintln(w, r.TransferEncoding)
		fmt.Fprintln(w, r.URL.Scheme)
		fmt.Fprintln(w, r.URL.Path)
		fmt.Fprintln(w, r.URL.Fragment)
		// w.Write([]byte("Hello World!"))
	})

	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, "Hello from the handler")
		// w.Write([]byte("Hello World!"))
	})

	fileServer := http.FileServer(http.Dir("./files"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	fileEmbedDir, _ := fs.Sub(resources, "files")
	fileEmbedServer := http.FileServer(http.FS(fileEmbedDir))
	mux.Handle("/static_embed/", http.StripPrefix("/static_embed/", fileEmbedServer))

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, helloHtml)
	})

	mux.HandleFunc("/hello2", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./resources/index.html")
	})

	mux.HandleFunc("/template", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./templates/simple.gohtml")
		if err != nil {
			fmt.Fprint(w, err)
		}
		t.ExecuteTemplate(w, "simple.gohtml", "hello world")
	})

	mux.HandleFunc("/coba", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./templates/coba.gohtml")
		if err != nil {
			fmt.Fprint(w, err)
		}
		t.ExecuteTemplate(w, "coba.gohtml", map[string]interface{}{
			"Title": "Hello Coba",
			"Name":  "John Doe",
			"Address": map[string]interface{}{
				"Street": "123 Main St",
			},
		})
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server listening on port", server.Addr)

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
