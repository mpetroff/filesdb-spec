/*
 * Serves files from SQLite file database
 * Copyright (c) 2015 Matthew Petroff
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET and HEAD methods
	if r.Method != "GET" && r.Method != "HEAD" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	// Separate request into SQLite database filename and requested file
	if len(r.URL.Path) < 1 {
		http.NotFound(w, r)
		return
	}
	split := strings.SplitN(r.URL.Path[1:], "/", 2)
	if len(split) < 2 {
		http.NotFound(w, r)
		return
	}
	dbFilename := "./" + split[0] + ".filesdb"
	requestedFile := split[1]

	// Check if the database file exists
	if _, err := os.Stat(dbFilename); os.IsNotExist(err) {
		log.Print(dbFilename, " not found")
		http.NotFound(w, r)
		return
	}

	// Open the database
	db, err := sql.Open("sqlite3", dbFilename)
	if err != nil {
		log.Print(err)
		http.NotFound(w, r)
		return
	}
	defer db.Close()

	// Query for the requested file
	var data []byte
	if err := db.QueryRow("select data from files where filename=?", requestedFile).Scan(&data); err != nil {
		if err := db.QueryRow("select data from files where filename=?", requestedFile+"index.htm").Scan(&data); err != nil {
			if err := db.QueryRow("select data from files where filename=?", requestedFile+"index.html").Scan(&data); err != nil {
				log.Print("`", requestedFile, "` not found in db")
				http.NotFound(w, r)
				return
			}
		}
	}

	// Serve the requested file
	w.Header().Set("Content-Type", http.DetectContentType(data))
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Cache-Control", "public")
	if r.Method == "GET" {
		w.Write(data)
	}
}

func main() {
	log.Print("Server started")

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
