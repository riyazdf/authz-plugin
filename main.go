package main

import (
	"fmt"
	"log/syslog"
	"os"
	"regexp"

	"github.com/Sirupsen/logrus"
	logrus_syslog "github.com/Sirupsen/logrus/hooks/syslog"
	"github.com/docker/go-plugins-helpers/authorization"
)

type authorizer struct{}

func main() {
	fmt.Println("starting authz plugin")
	h := authorization.NewHandler(NewAuthZPlugin())
	if err := h.ServeUnix("", "/run/docker/plugins/authz-plugin.sock"); err != nil {
		fmt.Printf("Error starting plugin: %v\n", err)
		os.Exit(1)
	}
}

// NewAuthZPlugin returns a new authorization plugin, as specified by the AuthZ interface
func NewAuthZPlugin() authorization.Plugin {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	hook, err := logrus_syslog.NewSyslogHook("", "", syslog.LOG_DEBUG, "")
	if err != nil {
		// for now this is best effort due to CI integration
		logrus.Errorf("Error retrieving syslog hook: %v", err)
	} else {
		logrus.Debug("Adding syslog hook for logging...")
		logrus.AddHook(hook)
	}
	return authorizer{}
}

// AuthZReq vets incoming requests from the docker daemon, and rejects any invalid requests
func (a authorizer) AuthZReq(req authorization.Request) authorization.Response {
	var err error
	logrus.Infof("received AuthZ request, method: '%s', url: '%s'", req.RequestMethod, req.RequestURI)
	// If the request matches a `docker volume` command, we'll block it as a POC
	isVol, err := regexp.MatchString("/volumes", req.RequestURI)
	// fail closed on any regex errors
	if err != nil {
		return authorization.Response{
			Allow: false,
			Err:   fmt.Sprintf("unauthorized request due to error parsing API url: '%s': %v", req.RequestURI, err),
		}
	}
	if isVol {
		return authorization.Response{
			Allow: false,
			Err:   fmt.Sprintf("unauthorized volume request: '%s'", req.RequestURI),
		}
	}

	logrus.Infof("accepting request, url: '%s'", req.RequestURI)
	return authorization.Response{
		Allow: true,
		Msg:   "request authorized",
	}
}

// AuthZRes will allow all requests back from the daemon, as all incoming requests have already been handled by AuthZReq
func (a authorizer) AuthZRes(req authorization.Request) authorization.Response {
	return authorization.Response{
		Allow: true,
		Msg:   "request authorized",
	}
}
