package infra

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// spawner is the localSpawner infra uses to start services from source in local mode.
// A package var so tests inject a fake and the unit tests never launch real processes.
var spawner localSpawner = osSpawner{}

// Local run source roots, resolved CWD-independently from the e2e module root and
// env-overridable. Layout: mapexOS/{e2e_tests,workspace_go,workspace_js} are siblings;
// mapexLNS is a sibling repo of mapexOS (one level further up).
var (
	goSrcDir  = resolveFromRoot("MAPEX_WORKSPACE_GO", "..", "workspace_go")
	jsSrcDir  = resolveFromRoot("MAPEX_WORKSPACE_JS", "..", "workspace_js")
	lnsSrcDir = resolveFromRoot("MAPEX_LNS_DIR", "..", "..", "mapexLNS")
)

// goRun builds the local spec for a workspace_go service: `go run ./services/<dir>/src/main.go`.
func goRun(serviceDir string) *localSpec {
	return &localSpec{cmd: "go", args: []string{"run", "./services/" + serviceDir + "/src/main.go"}, dir: goSrcDir}
}

// npmDev builds the local spec for a workspace_js service: `npm run dev -w services/<name>`.
func npmDev(workspace string) *localSpec {
	return &localSpec{cmd: "npm", args: []string{"run", "dev", "-w", "services/" + workspace}, dir: jsSrcDir}
}

// lnsRun builds the local spec for the mapexLNS sibling repo: `go run ./src/main.go`.
func lnsRun() *localSpec {
	return &localSpec{cmd: "go", args: []string{"run", "./src/main.go"}, dir: lnsSrcDir}
}

// resolveFromRoot returns a dir resolver: env override when set, else the e2e module
// root (walk to go.mod) joined with rel — so the path is the same from any CWD.
func resolveFromRoot(env string, rel ...string) func() string {
	return func() string {
		if d := os.Getenv(env); d != "" {
			return d
		}
		root, err := repoRoot()
		if err != nil {
			return filepath.Join(rel...)
		}
		return filepath.Join(append([]string{root}, rel...)...)
	}
}

// osSpawner starts a service as a child process in its OWN process group. The group is
// what makes stop() reliable: `go run` and `npm` each fork the actual server as a child
// process, so signalling only the parent would orphan the listener on its port; killing
// the whole group reaches the real server too.
type osSpawner struct{}

// start launches spec's command in spec.dir, tees stdout+stderr to a per-service log
// file (so a service that fails to come up is diagnosable), and returns a handle.
func (osSpawner) start(svc string, spec localSpec) (runningProc, error) {
	logPath := filepath.Join(localLogDir(), svc+".log")
	lf, err := os.Create(logPath)
	if err != nil {
		return nil, fmt.Errorf("create log %s: %w", logPath, err)
	}
	cmd := exec.Command(spec.cmd, spec.args...)
	cmd.Dir = spec.dir()
	// spec.env is layered AFTER os.Environ() so per-service overrides win (exec
	// keeps the last value for a duplicated key). Nil env leaves the inherited
	// environment untouched.
	cmd.Env = append(os.Environ(), spec.env...)
	cmd.Stdout, cmd.Stderr = lf, lf
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		_ = lf.Close()
		return nil, fmt.Errorf("start %s (%s %v in %s): %w", svc, spec.cmd, spec.args, cmd.Dir, err)
	}
	log.Printf("[INFRA] spawned %s pid=%d log=%s", svc, cmd.Process.Pid, logPath)
	return &osProc{svc: svc, cmd: cmd, log: lf}, nil
}

// osProc is a running local service plus its captured log file.
type osProc struct {
	svc string
	cmd *exec.Cmd
	log *os.File
}

// stop signals the process GROUP (negative pid) so the child binary go run/npm forked
// is terminated too, then reaps it and closes the log. Because Setpgid put the child at
// the head of a new group, its pid equals the group id.
func (p *osProc) stop() error {
	defer p.log.Close()
	pgid := p.cmd.Process.Pid
	if err := syscall.Kill(-pgid, syscall.SIGTERM); err != nil {
		return fmt.Errorf("signal %s group (pgid=%d): %w", p.svc, pgid, err)
	}
	_ = p.cmd.Wait()
	return nil
}

// localLogDir returns (creating) the directory spawned services' output is written to.
// Override with MAPEX_E2E_LOG_DIR; defaults under the OS temp dir.
func localLogDir() string {
	dir := os.Getenv("MAPEX_E2E_LOG_DIR")
	if dir == "" {
		dir = filepath.Join(os.TempDir(), "mapex-e2e-logs")
	}
	_ = os.MkdirAll(dir, 0o755)
	return dir
}
