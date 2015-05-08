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
	reposRootFlag     = flag.String("repos.root", "/var/repos", "Location of the repositories.")
	authUserFlag      = flag.String("auth.user", "", "Username for basic auth.")
	authPassFlag      = flag.String("auth.password", "", "Password for basic auth.")
)

func hasAuth() bool {
	return len(*authUserFlag) > 0 && len(*authPassFlag) > 0
}

func validateAuth(w http.ResponseWriter, r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	if ok && username == *authUserFlag && password == *authPassFlag {
		return true
	}

	return false
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Processing %s\n", r.URL.Path[0:])

	if hasAuth() && !validateAuth(w, r) {
		w.Header().Set("WWW-Authenticate", `Basic realm="gitd"`)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
	}

	parsedRoute := gitd.MatchRoute(*reposRootFlag, r)
	if parsedRoute != nil {
		parsedRoute.Dispatch(w, r)
	} else {
		fmt.Fprintf(w, "nothing to see here\n")
	}
}

func init() {
	flag.Parse()
	*reposRootFlag, _ = filepath.Abs(*reposRootFlag)
}

func main() {
	log.Printf("Starting server on %s, repos=%s", *listenAddressFlag, *reposRootFlag)

	http.HandleFunc("/", handler)
	http.ListenAndServe(*listenAddressFlag, nil)
}
