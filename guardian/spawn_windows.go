package guardian

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gladiusio/gladius-guardian/win"
	log "github.com/sirupsen/logrus"
)

func (gg *GladiusGuardian) spawnProcess(name, location string, env []string, timeout *time.Duration) (*exec.Cmd, error) {
	p := exec.Command("cmd.exe", "/C", "start", location)
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
		return nil, fmt.Errorf("\nError starting process: %s", err)
	}

	// Timeout test, can we find the process before the timeout?
	time.Sleep(*timeout)
	process, err := GetProcess(location + ".exe")
	if err != nil {
		return nil, fmt.Errorf("could not finding process %s or failed to start before timeout, check the logs for errors", name)
	}

	// when process exits, call this
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
	processID, err := win.GetProcessWindows(name)
	if err != nil {
		return nil, fmt.Errorf("error getting process id")
	}
	process, err := os.FindProcess(processID)
	if err != nil {
		return nil, fmt.Errorf("error loading process")
	}

	return process, nil
}
