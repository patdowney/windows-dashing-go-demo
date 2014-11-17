package main

import (
	"fmt"
	"os"
	"path/filepath"

	"bitbucket.org/kardianos/osext"
	"bitbucket.org/kardianos/service"
	"github.com/patdowney/dashing-go"

	_ "./jobs"
)

var log service.Logger

func main() {
	var name = "WindowsDashingDemo"
	var displayName = "Windows Dashing Demo"
	var desc = "This is a Windows Dashing-go Demo"

	var s, err = service.NewService(name, displayName, desc)
	log = s

	if err != nil {
		fmt.Printf("%s unable to start: %s", displayName, err)
		return
	}

	if len(os.Args) > 1 {
		var err error
		verb := os.Args[1]
		switch verb {
		case "install":
			err = s.Install()
			if err != nil {
				fmt.Printf("Failed to install: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" installed.\n", displayName)
		case "remove":
			err = s.Remove()
			if err != nil {
				fmt.Printf("Failed to remove: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" removed.\n", displayName)
		case "run":
			startDashingServer()
		case "start":
			err = s.Start()
			if err != nil {
				fmt.Printf("Failed to start: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" started.\n", displayName)
		case "stop":
			err = s.Stop()
			if err != nil {
				fmt.Printf("Failed to stop: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" stopped.\n", displayName)
		}
		return
	}
	err = s.Run(
		func() error {
			go startDashingServer()
			return nil
		},
		func() error {
			stopDashingServer()
			return nil
		})
	if err != nil {
		s.Error(err.Error())
	}
}

var exit = make(chan struct{})

func startDashingServer() {
	go func() {
		workingDir, err := osext.ExecutableFolder()
		if err != nil {
			log.Error("failed to fetch executable folder")
			return
		}

		os.Chdir(workingDir)
		dashing.StartWithStaticDirectory(filepath.Join(workingDir, "public"))
	}()
	for {
		select {
		case <-exit:
			return
		}
	}
}

func stopDashingServer() {
	exit <- struct{}{}
}
