// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "net/http"
    "embed"
)

//go:embed templates
var templates embed.FS

var addr = flag.String("addr", ":8100", "http service address")

func serveRegister(w http.ResponseWriter, r *http.Request) {
    http.ServeFileFS(w, r, templates, "templates/register.html")
}

func servePlay(w http.ResponseWriter, r *http.Request) {
    http.ServeFileFS(w, r, templates, "templates/play.html")
}

func main() {
    flag.Parse()
    hub := newHub()
    go hub.run()
    hub.newQuestion()

    mux := http.NewServeMux()
    mux.HandleFunc("GET /quizlounge", serveRegister)
    mux.HandleFunc("GET /quizlounge/play", servePlay)
    mux.HandleFunc("/quizlounge/ws", func(w http.ResponseWriter, r *http.Request) {
        serveWs(hub, w, r)
    })

    srv := &http.Server{
        Addr: *addr,
        Handler: mux,
    }
    fmt.Printf("serving quizlounge at %v\n", *addr)
    err := srv.ListenAndServe()
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
    os.Exit(1)
}
