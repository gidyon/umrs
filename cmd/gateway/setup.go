package main

import (
	"github.com/gidyon/file-handlers/static"
	"github.com/pkg/errors"
	"net/http"
)

func staticFilesHandler(dir, pushFile, URLPathPrefix string) http.Handler {
	// parse list of push files
	pushFiles, err := readPushFiles(pushFile)
	handleErr(err)

	notFound := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "dist/index.html")
	})

	pushMap := map[string][]string{
		"dist/index.html": pushFiles,
	}

	sh, err := static.NewHandler(&static.ServerOptions{
		RootDir:         dir,
		Index:           "index.html",
		NotFoundHandler: notFound,
		URLPathPrefix:   URLPathPrefix,
		PushContent:     pushMap,
		FallBackIndex:   true,
	})
	handleErr(err)

	return sh
}

func apiDocumentationHandler() http.Handler {
	handler, err := static.NewHandler(&static.ServerOptions{
		RootDir: "./",
		Index:   "static/dist/index.html",
		AllowedDirs: []string{
			"api/swagger",
			"static/dist",
		},
		NotFoundHandler: http.NotFoundHandler(),
	})
	handleErr(errors.Wrap(err, "failed to setup API documentation"))

	return handler
}
