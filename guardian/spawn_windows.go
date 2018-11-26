package guardian

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gladiusio/gladius-guardian/win"
	log "github.com/sirupsen/logrus"
)

// spawnProcess - spawn a windows process
func (gg *GladiusGuardian) spawnProcess(name, location string, env []string, timeout *time.Duration) (*exec.Cmd, error) {
	log.Info("Starting service")
	p := exec.Command("cmd.exe", "/C", location)
	p.Env = append(os.Environ(), env...)

	// Create standard err and out pipes
	stdOut, err := p.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("Error creating StdoutPipe for command: %s", err)
	}
	stdErr, err := p.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("Error creating StderrPipe for command: %s", err)
	}

	// Pipe stdout to the logs
	scanner := bufio.NewScanner(stdOut)
	stdErrScanner := bufio.NewScanner(stdErr)
	go func() {
		defer stdOut.Close()
		for scanner.Scan() {
			gg.AppendToLog(name, scanner.Text())
		}
		err = scanner.Err()
		if err != nil {
			gg.AppendToLog(name, "STDOUT ERR: "+err.Error())
		}
	}()
	go func() {
		defer stdErr.Close()
		for stdErrScanner.Scan() {
			gg.AppendToLog(name, stdErrScanner.Text())
		}
		err = stdErrScanner.Err()
		if err != nil {
			gg.AppendToLog(name, "STDERR ERR: "+err.Error())
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
		return nil, fmt.Errorf("\nError starting process: %s", err)
	}

	// Timeout test, can we find the process before the timeout?
	time.Sleep(*timeout)
	process, err := GetProcess(location + ".exe")
	if err != nil {
		return nil, fmt.Errorf("could not finding process %s or failed to start before timeout, check the logs for errors", name)
	}

	// this waits for the process to end
	go func() {
		_, err := process.Wait()
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

	return p, nil
}

// GetProcess - Returns process obj (Windows only)
func GetProcess(name string) (*os.Process, error) {
	// find process id by name
	processID, err := win.GetProcessWindows(name)
	if err != nil {
		return nil, fmt.Errorf("error getting process id")
	}
	// load the process up to return it
	process, err := os.FindProcess(processID)
	if err != nil {
		return nil, fmt.Errorf("error loading process")
	}

	return process, nil
}

// killProcess - kill a windows process
func killProcess(gg *GladiusGuardian, name string) error {
	// get the process by name
	process, err := GetProcess("gladius-" + name + ".exe")
	if err != nil {
		return errors.New("could not find windows process")
	}
	log.Info("Stopping service")

	// kill the process
	err = process.Kill()
	if err != nil {
		return errors.New("could not kill windows process")
	}

	return nil
}
