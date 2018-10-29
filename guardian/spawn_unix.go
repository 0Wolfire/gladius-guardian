// +build linux darwin

package guardian

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func (gg *GladiusGuardian) spawnProcess(name, location string, env []string, timeout *time.Duration) (*exec.Cmd, error) {
	p := exec.Command(location)
	p.Env = env

	// Create standard err and out pipes
	stdOut, err := p.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("Error creating StdoutPipe for command: %s", err)
	}
	stdErr, err := p.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("Error creating StderrPipe for command: %s", err)
	}

	// Read both of those in
	scanner := bufio.NewScanner(stdOut)
	stdErrScanner := bufio.NewScanner(stdErr)
	go func() {
		defer stdOut.Close()
		for scanner.Scan() {
			gg.AppendToLog(name, scanner.Text())
		}
	}()
	go func() {
		defer stdErr.Close()
		for stdErrScanner.Scan() {
			gg.AppendToLog(name, stdErrScanner.Text())
		}
	}()

	// Start the command
	err = p.Start()
	if err != nil {
		log.WithFields(log.Fields{
			"exec_location":    location,
			"environment_vars": strings.Join(env, ", "),
			"err":              err,
		}).Warn("Couldn't spawn process")
		return nil, fmt.Errorf("Error starting process: %s", err)
	}

	go func() {
		err := p.Wait()
		gg.services[name] = nil // Set out service to nil when it dies
		if err != nil {
			// Only log errors if we didn't kill it
			if err.Error() != "signal: killed" {
				log.WithFields(log.Fields{
					"exec_location":    location,
					"environment_vars": strings.Join(env, ", "),
					"err":              err,
				}).Error("Service errored out")
				gg.AppendToLog(name, "Exiting... "+err.Error())
			}
		}
	}()

	// Wait for the process to start
	time.Sleep(*timeout)
	if p.ProcessState != nil { // ProcessState is only non-nil if p.Wait() concludes
		if p.ProcessState.Exited() {
			return nil, fmt.Errorf("process %s already exited, check the logs for errors", name)
		}
	}
	return p, nil
}

func killProcess(gg *GladiusGuardian, name string) error {
	service := gg.services[name]
	err := service.Process.Kill()
	if err != nil {
		return errors.New("could not kill unix process")
	}
	return nil
}
