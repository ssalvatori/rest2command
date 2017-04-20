package main

import (
	"encoding/json"
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"
	log "github.com/Sirupsen/logrus"

	"io/ioutil"
	"regexp"
	"fmt"

)

var (
	Port = ":8891"
	ConfigurationFile = "./configuration.json"
	Version = "0.0.0"
	API_Version = ""
	CredentialsFile = "./credentials.json"
)

func main() {

	setUpLog()
	setUp()

	// subscribe to SIGINT signals
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	mux := buildHandlers()

	srv := &http.Server{
		Addr:	Port,
		Handler: mux,
	}

	go func() {
		// service connections
		err := srv.ListenAndServe()

		if  err != nil {
			log.Info("listen: %s\n", err.Error())
		}
	}()

	<-stopChan // wait for SIGINT
	log.Info("Shutting down server...")

	// shut down gracefully, but wait no longer than 5 seconds before halting
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)

	log.Info("Server gracefully stopped")

}

type Configuration struct {
	Url string `json:"url"`
	Command   string `json:"command"`
}

func getConfigurations() []Configuration {
	raw, err := ioutil.ReadFile(ConfigurationFile)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	var c []Configuration
	json.Unmarshal(raw, &c)
	return c
}

func buildCommands() map[string]string {
	commands := make(map[string] string)

	for _, configuration := range getConfigurations() {
		log.Debug(API_Version+configuration.Url, " -> ", configuration.Command)
		commands[API_Version+configuration.Url] = configuration.Command
	}

	return commands
}

func buildHandlers() *http.ServeMux {
	mux := http.NewServeMux()
	log.Info("Setting handlers")

	commands := buildCommands()

	for key := range commands {
		mux.Handle(key, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request, ) {
			log.Debug("Command: ", commands[r.RequestURI])
			log.Debug("Request origin: ", r.RemoteAddr)
			io.WriteString(w, "Finished!")
		}))
	}

	return mux
}

func setUp() {
	if os.Getenv("PORT") != "" {
		Port = ":"+os.Getenv("PORT")
	}
	if os.Getenv("FILE_CONFIGURATION") != "" {
		ConfigurationFile = os.Getenv("FILE_CONFIGURATION")
	}
	if os.Getenv("FILE_CREDENTIALS") != "" {
		CredentialsFile = os.Getenv("FILE_CREDENTIALS")
	}

	API_Version = fmt.Sprintf("/v%s",getAPIVersion(Version))
	log.Info("Context: ",API_Version)
}

func getAPIVersion(version string) string {
	versionRegExp := regexp.MustCompile(`(\d+).(\d+).(\d+)`)
	if versionRegExp.MatchString(version) {
		versions := versionRegExp.FindStringSubmatch(version)
		return versions[1]
	}
	return "0"
}

func setUpLog() {

	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		log.SetLevel(log.DebugLevel)
		break
	case "info":
		log.SetLevel(log.InfoLevel)
		break
	case "warn":
		log.SetLevel(log.WarnLevel)
		break
	case "error":
		log.SetLevel(log.ErrorLevel)
		break
	case "fatal":
		log.SetLevel(log.FatalLevel)
		break
	case "panic":
		log.SetLevel(log.PanicLevel)
		break
	default:
		log.SetLevel(log.InfoLevel)
		break
	}

}
