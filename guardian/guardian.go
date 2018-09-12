package guardian

import (
	"errors"
	"os/exec"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// New returns a new GladiusGuardian object with the specified spawn timeout
func New() *GladiusGuardian {
	return &GladiusGuardian{mux: &sync.Mutex{}}
}

// GladiusGuardian manages the various gladius processes
type GladiusGuardian struct {
	mux          *sync.Mutex
	spawnTimeout *time.Duration
	networkd     *exec.Cmd
	controld     *exec.Cmd
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
		Running: true,
	}
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
	services["networkd"] = newServiceStatus(gg.networkd)
	services["controld"] = newServiceStatus(gg.controld)

	return services
}

func (gg *GladiusGuardian) StopAll() error {
	gg.mux.Lock()
	defer gg.mux.Unlock()

	return nil
}

func (gg *GladiusGuardian) StartControld(env []string) error {
	gg.mux.Lock()
	defer gg.mux.Unlock()

	if err := gg.checkTimeout(); err != nil {
		return err
	}

	// TODO: Let this location be configurable
	p, err := spawnProcess("gladius-controld", env, gg.spawnTimeout)
	if err != nil {
		return nil
	}
	gg.controld = p
	return nil
}

func (gg *GladiusGuardian) StartNetworkd(env []string) error {
	gg.mux.Lock()
	defer gg.mux.Unlock()

	if err := gg.checkTimeout(); err != nil {
		return err
	}

	// TODO: Let this location be configurable
	p, err := spawnProcess("gladius-networkd", env, gg.spawnTimeout)
	if err != nil {
		return nil
	}
	gg.controld = p
	return nil
}

func (gg *GladiusGuardian) StopControld() error {
	gg.mux.Lock()
	defer gg.mux.Unlock()

	return nil
}

func (gg *GladiusGuardian) StopNetworkd() error {
	gg.mux.Lock()
	defer gg.mux.Unlock()

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
			log.Warn("Couldn't spawn process " + err.Error())
		}
	}(p)

	// Wait for the process to start
	time.Sleep(*timeout)

	return p, nil
}
