package backend

import (
	"bufio"
	"github.com/prometheus/common/log"
	"github.com/zeebo/errs"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"sync"
)

type Supervisor struct {
	log     *zap.Logger
	cmd     *exec.Cmd
	mu      sync.Mutex
	lastErr error
	running bool
}

type RunCfg struct {
	Bucket string
	Grant  string
}

func (s *Supervisor) Start() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cmd = exec.Command("registry", "serve", "config.yml")
	cfg, err := s.readConfig()
	if err != nil {
		return errs.New("Storj gateway is not yet configured")
	}
	s.cmd.Env = []string{"REGISTRY_STORAGE_STORJ_BUCKET=" + cfg.Bucket, "REGISTRY_STORAGE_STORJ_ACCESSGRANT=" + cfg.Grant}

	pipeOut, err := s.cmd.StdoutPipe()
	if err != nil {
		return errs.Wrap(err)
	}

	pipeErr, err := s.cmd.StderrPipe()
	if err != nil {
		return errs.Wrap(err)
	}

	err = s.cmd.Start()
	if err != nil {
		return err
	}
	s.running = true

	go func() {

		scanner := bufio.NewScanner(pipeOut)

		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			log.Info(scanner.Text())
		}
	}()

	go func() {

		scanner := bufio.NewScanner(pipeErr)

		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			log.Info(scanner.Text())
		}
	}()
	go func() {
		defer func() {
			s.running = false
		}()
		s.lastErr = s.cmd.Wait()
		log.Info("Registry process has been stopped")
	}()
	return nil
}

func (s *Supervisor) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.running {
		return nil
	}
	s.running = false
	return s.cmd.Process.Kill()
}

func (s *Supervisor) Configure(bucket string, grant string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	out, err := yaml.Marshal(&RunCfg{
		Bucket: bucket,
		Grant:  grant,
	})
	if err != nil {
		return err
	}
	return os.WriteFile("storj.yaml", out, 0600)
}

func (s *Supervisor) readConfig() (cfg RunCfg, err error) {
	in, err := os.ReadFile("storj.yaml")
	if err != nil {
		return
	}
	err = yaml.Unmarshal(in, &cfg)
	return
}

func (s *Supervisor) Status() (string, error) {
	if s.running {
		return "running", nil
	}
	if _, err := s.readConfig(); err != nil {
		return "missing", nil
	}
	return "stopped", nil
}
