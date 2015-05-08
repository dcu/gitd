package gitd

import (
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
)

var Routes = []RouteMatcher{
	RouteMatcher{Matcher: regexp.MustCompile("(.*?)/info/refs$"), Handler: getInfoRefs},
	RouteMatcher{Matcher: regexp.MustCompile("(.*?)/git-upload-pack$"), Handler: uploadPack},
	RouteMatcher{Matcher: regexp.MustCompile("(.*?)/git-receive-pack$"), Handler: receivePack},
}

type RouteFunc func(route *Route, w http.ResponseWriter, r *http.Request)
type RouteMatcher struct {
	Matcher *regexp.Regexp
	Handler RouteFunc
}

type Route struct {
	RepoName     string
	RepoPath     string
	File         string
	MatchedRoute RouteMatcher
}

func (route *Route) Dispatch(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	route.MatchedRoute.Handler(route, w, r)
}

func NewParsedRoute(repoName string, repoPath string, file string, matcher RouteMatcher) *Route {
	return &Route{RepoName: repoName, RepoPath: repoPath, File: file, MatchedRoute: matcher}
}

func MatchRoute(reposRoot string, r *http.Request) *Route {
	path := r.URL.Path[1:]

	for _, routeHandler := range Routes {
		matches := routeHandler.Matcher.FindStringSubmatch(path)
		if matches != nil {
			repoName := matches[1]
			file := strings.Replace(path, repoName+"/", "", 1)
			repoPath := filepath.Join(reposRoot, repoName)

			return NewParsedRoute(repoName, repoPath, file, routeHandler)
		}
	}

	log.Printf("No route found for: %s", path)
	return nil
}
