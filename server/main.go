package main

import (
"flag"
"log"
"net/http"
)

func main() {
addr := flag.String("addr", ":8080", "listen address")
dir := flag.String("dir", "public", "directory to serve")
flag.Parse()

files := http.FileServer(http.Dir(*dir))
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Cache-Control", "no-store") // handy while developing
files.ServeHTTP(w, r)
})

log.Printf("gallery running at http://localhost%s", *addr)
log.Fatal(http.ListenAndServe(*addr, nil))
}
