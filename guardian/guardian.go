package guardian

import (
	"errors"
	"os"
	"os/exec"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// New returns a new GladiusGuardian object with the specified spawn timeout
func New(timeout time.Duration) *GladiusGuardian {
	return &GladiusGuardian{mux: &sync.Mutex{}}
}

// GladiusGuardian manages the various gladius processes
type GladiusGuardian struct {
	mux          *sync.Mutex
	spawnTimeout *time.Duration
	networkd     *os.Process
	controld     *os.Process
}

func (gg *GladiusGuardian) SetTimeout(t *time.Duration) {
	gg.mux.Lock()
	defer gg.mux.Unlock()

	gg.spawnTimeout = t
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

func spawnProcess(location string, env []string, timeout *time.Duration) (*os.Process, error) {
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

	return p.Process, nil
}
