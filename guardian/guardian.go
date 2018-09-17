package guardian

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	multierror "github.com/hashicorp/go-multierror"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// New returns a new GladiusGuardian object with the specified spawn timeout
func New() *GladiusGuardian {
	return &GladiusGuardian{
		mux:                &sync.Mutex{},
		registeredServices: make(map[string]*serviceSettings),
		services:           make(map[string]*exec.Cmd),
	}
}

// GladiusGuardian manages the various gladius processes
type GladiusGuardian struct {
	mux                *sync.Mutex
	spawnTimeout       *time.Duration
	registeredServices map[string]*serviceSettings
	services           map[string]*exec.Cmd
}

type serviceSettings struct {
	env      []string
	execName string
}

type serviceStatus struct {
	Running  bool     `json:"running"`
	PID      int      `json:"pid"`
	Env      []string `json:"environment_vars"`
	Location string   `json:"executable_location"`
}

func newServiceStatus(p *exec.Cmd) *serviceStatus {
	if p != nil {
		return &serviceStatus{
			Running:  true,
			PID:      p.Process.Pid,
			Env:      p.Env,
			Location: p.Path,
		}
	}
	return &serviceStatus{
		Running: false,
	}
}

func (gg *GladiusGuardian) RegisterService(name, execLocation string, env []string) {
	log.WithFields(log.Fields{
		"service_name":     name,
		"exec_location":    execLocation,
		"environment_vars": strings.Join(env, ", "),
	}).Debug("Registered new service")
	gg.registeredServices[name] = &serviceSettings{env: env, execName: execLocation}
	gg.services[name] = nil // So it's still returned when we list services
}

func (gg *GladiusGuardian) SetTimeout(t *time.Duration) {
	gg.mux.Lock()
	defer gg.mux.Unlock()

	gg.spawnTimeout = t
}

func (gg *GladiusGuardian) GetServicesStatus() map[string]*serviceStatus {
	gg.mux.Lock()
	defer gg.mux.Unlock()

	services := make(map[string]*serviceStatus)
	for serviceName, service := range gg.services {
		services[serviceName] = newServiceStatus(service)
	}

	return services
}

func (gg *GladiusGuardian) StopAll() error {
	gg.mux.Lock()
	defer gg.mux.Unlock()

	var result *multierror.Error

	for sName, s := range gg.services {
		if s != nil {
			err := s.Process.Kill()
			result = multierror.Append(result, fmt.Errorf("error stopping service %s: %s", sName, err))
		}
		result = multierror.Append(result, fmt.Errorf("service not running: %s", sName))
	}
	err := result.ErrorOrNil()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Warn("Error stoping one or more service")
	}
	return result.ErrorOrNil()
}

func (gg *GladiusGuardian) StartService(name string, env []string) error {
	gg.mux.Lock()
	defer gg.mux.Unlock()

	serviceSettings, ok := gg.registeredServices[name]
	if !ok {
		return errors.New("attempted to start unregistered service")
	}

	if len(env) == 0 {
		env = viper.GetStringSlice("DefaultEnvironment")
	}

	if err := gg.checkTimeout(); err != nil {
		return err
	}

	p, err := spawnProcess(serviceSettings.execName, serviceSettings.env, gg.spawnTimeout)
	if err != nil {
		return err
	}
	gg.services[name] = p
	log.WithFields(log.Fields{
		"service_name":     name,
		"exec_location":    serviceSettings.execName,
		"environment_vars": strings.Join(env, ", "),
	}).Debug("Started service")
	return nil
}

func (gg *GladiusGuardian) StopService(name string) error {
	gg.mux.Lock()
	defer gg.mux.Unlock()

	serviceSettings, ok := gg.registeredServices[name]
	if !ok {
		return errors.New("attempted to stop unregistered service")
	}

	service := gg.services[name]
	if service == nil {
		return errors.New("service is not running so can not stop")
	}

	err := service.Process.Kill()
	if err != nil {
		log.WithFields(log.Fields{
			"service_name":     name,
			"exec_location":    serviceSettings.execName,
			"environment_vars": strings.Join(serviceSettings.env, ", "),
			"err":              err,
		}).Warn("Couldn't kill service")
		return errors.New("couldn't kill service, error was: " + err.Error())
	}

	return nil
}

func (gg *GladiusGuardian) checkTimeout() error {
	if gg.spawnTimeout == nil {
		return errors.New("spawn timeout not set, please set it before a process is spawned")
	}
	return nil
}

func spawnProcess(location string, env []string, timeout *time.Duration) (*exec.Cmd, error) {
	p := exec.Command(location)
	p.Env = env

	go func(proc *exec.Cmd) {
		// TODO: Configure logging through API/defualts
		_, err := proc.CombinedOutput()
		if err != nil {
			log.WithFields(log.Fields{
				"exec_location":    location,
				"environment_vars": strings.Join(env, ", "),
				"err":              err,
			}).Warn("Couldn't spawn process")
		}
	}(p)

	// Wait for the process to start
	time.Sleep(*timeout)

	return p, nil
}
