package main

import (
	"flag"
	"fmt"
	"github.com/dcu/gitd/gitd"
	"log"
	"net/http"
	"path/filepath"
)

var (
	listenAddressFlag = flag.String("web.listen-address", ":4000", "Address on which the git server will be served.")
	reposRoot         = flag.String("repos.root", "/var/repos", "Location of the repositories.")
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Processing %s\n", r.URL.Path[0:])

	parsedRoute := gitd.MatchRoute(*reposRoot, r)
	if parsedRoute != nil {
		parsedRoute.Dispatch(w, r)
	} else {
		fmt.Fprintf(w, "nothing to see here\n")
	}
}

func init() {
	flag.Parse()
	*reposRoot, _ = filepath.Abs(*reposRoot)
}

func main() {
	log.Printf("Starting server on %s, repos=%s", *listenAddressFlag, *reposRoot)

	http.HandleFunc("/", handler)
	http.ListenAndServe(*listenAddressFlag, nil)
}
