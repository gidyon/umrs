package main

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/gidyon/gateway"
	http_middleware "github.com/gidyon/micros/pkg/http"
	"net/http"
	"os"
	"strings"
)

const (
	certPath         = "/home/gideon/go/src/github.com/gidyon/umrs/certs/ledger/cert.pem"
	keyPath          = "/home/gideon/go/src/github.com/gidyon/umrs/certs/ledger/key.pem"
	servicesFilePath = "/home/gideon/go/src/github.com/gidyon/umrs/configs/services.dev.yml"
	pushFilesPath    = "/home/gideon/go/src/github.com/gidyon/umrs/configs/pushfiles.yml"
)

var (
	staticPathPrefix = flag.String("static-path-prefix", "/app/", "URL path prefix for static files")
	pushFilesFile    = flag.String("push-file", pushFilesPath, "File containing list of server push files")
	servicesFile     = flag.String("services-file", servicesFilePath, "File containing list of to services")
	certFile         = flag.String("cert", certPath, "Path to public key")
	keyFile          = flag.String("key", keyPath, "Path to tls key")
	port             = flag.String("port", ":443", "Port to serve files")
	dir              = flag.String("dir", "dist/", "Static files directory")
	insecure         = flag.Bool("insecure", false, "Whether to use insecure http")
	cors             = flag.Bool("cors", false, "Whether to run app with CORS enabled")
	env              = flag.Bool("env", false, "Read parameters from env variables")
)

func main() {
	flag.Parse()

	if *env {
		// Set from environemnt variables
		*staticPathPrefix = setIfEmpty(strings.TrimSpace(os.Getenv("STATIC_URL_PREFIX")), *staticPathPrefix)
		*pushFilesFile = setIfEmpty(os.Getenv("PUSH_FILE"), *pushFilesFile)
		*servicesFile = setIfEmpty(strings.TrimSpace(os.Getenv("SERVICES_FILE")), *servicesFile)
		*certFile = setIfEmpty(os.Getenv("TLS_CERT_FILE"), *certFile)
		*keyFile = setIfEmpty(os.Getenv("TLS_KEY_FILE"), *keyFile)
		*port = setIfEmpty(os.Getenv("SERVICE_PORT"), *port)
		*dir = setIfEmpty(os.Getenv("STATIC_DIR"), *dir)
	}

	g, err := gateway.NewFromFile(*servicesFile)
	handleErr(err)

	// Update documentation handler
	docHandler := apiDocumentationHandler()
	docPath := "/api/umrs/documentation/"
	g.Handle(docPath, http.StripPrefix(docPath, docHandler))

	// Update static files handler
	sprefix := strings.TrimPrefix(strings.TrimSuffix(*staticPathPrefix, "/"), "/")
	if *staticPathPrefix == "/" {
		sprefix = "/"
	} else {
		sprefix = "/" + sprefix + "/"
	}
	// Static file server
	staticHandler := staticFilesHandler(*dir, *pushFilesFile, sprefix)
	g.Handle(sprefix, http.StripPrefix(sprefix, staticHandler))

	// Update endpoints
	updateEndpoints(g)

	if *cors {
		// Enable CORS
		g.AddMiddlewares(http_middleware.SupportCORS)
	}

	gh := g.Handler()

	*port = ":" + strings.TrimPrefix(*port, ":")

	logrus.Infof("started gateway on port: %s", *port)

	if *insecure {
		logrus.Fatalln(http.ListenAndServe(*port, gh))
	} else {
		logrus.Fatalln(http.ListenAndServeTLS(*port, *certFile, *keyFile, gh))
	}
}

func setIfEmpty(val, def string) string {
	if val == "" && def == "" {
		return ""
	}
	if val == "" {
		return def
	}
	return val
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
