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

	"os/exec"
)

var (
	Port = ":9999"
	ConfigurationFile = "/etc/rest2command/configuration.json"
	Version = "1.0.0"
	BuildTime = time.Now().String()
	GitHash   = "undefined"
	API_Version = ""
	CredentialsFile = "./credentials.json"
)

type Configuration struct {
	Url     string `json:"url"`
	Command string `json:"command"`
	Args    string `json:"args"`
}
type Command struct {
	Command string
	Args    string
}

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

func getConfigurations(configuration string) []Configuration {
	log.Debug("Loading file ",configuration)
	raw, err := ioutil.ReadFile(configuration)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	var c []Configuration
	json.Unmarshal(raw, &c)
	return c
}

func buildCommands(configurations []Configuration) map[string]Command {
	commands := make(map[string]Command)

	for _, configuration := range configurations {
		log.Debug(API_Version+configuration.Url, " -> ", configuration.Command)
		commands[API_Version+configuration.Url] = Command{
			Command: configuration.Command,
			Args:    configuration.Args,
		}
	}

	return commands
}

func buildHandlers() *http.ServeMux {
	mux := http.NewServeMux()
	log.Info("Setting handlers")

	commands := buildCommands(getConfigurations(ConfigurationFile))

	for key := range commands {
		//TODO wrap handlerfunc to check credentials and log information
		// https://medium.com/@matryer/the-http-handlerfunc-wrapper-technique-in-golang-c60bf76e6124
		mux.Handle(key, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request, ) {
			out := runCommand(commands[r.RequestURI])
			io.WriteString(w, out)
		}))
	}

	return mux
}

func runCommand(cmd Command) string {
	log.Debug("Cmd: ", cmd.Command, " Args ", cmd.Args)

	out, err := exec.Command(cmd.Command, cmd.Args).Output()
	if err != nil {
		log.Error(err)
		out = []byte("error")
	}
	log.Debug("Out: ", fmt.Sprintf("%s",out))

	return fmt.Sprintf("%s",out)
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
