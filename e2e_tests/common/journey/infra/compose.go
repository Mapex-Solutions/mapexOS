package infra

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// infraComposeFile resolves the INFRA docker-compose file (mongo/redis/clickhouse/
// nats/broker/minio). Override with MAPEX_COMPOSE_FILE; otherwise anchored at the
// module root (walk to go.mod), so it is CWD-independent.
func infraComposeFile() string { return composeFileFor("MAPEX_COMPOSE_FILE", "infra") }

// servicesComposeFile resolves the mapex SERVICES docker-compose file (mapex-iam,
// http-gateway, assets, router, events, triggers, workflow, lns, js-executor).
// Override with MAPEX_SERVICES_COMPOSE_FILE.
func servicesComposeFile() string { return composeFileFor("MAPEX_SERVICES_COMPOSE_FILE", "services") }

// composeFileFor returns env when set, else the sibling mapexOSDeploy/<subdir>/
// docker-compose.yml anchored at the module root (CWD-independent).
func composeFileFor(env, subdir string) string {
	if f := os.Getenv(env); f != "" {
		return f
	}
	root, err := repoRoot()
	if err != nil {
		return filepath.Join("..", "..", "mapexOSDeploy", subdir, "docker-compose.yml")
	}
	return filepath.Join(root, "..", "..", "mapexOSDeploy", subdir, "docker-compose.yml")
}

// repoRoot returns the e2e_tests module root (the directory holding go.mod) by walking
// up from this source file's compile-time path via runtime.Caller — so the result
// never depends on the current working directory.
func repoRoot() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("infra: cannot resolve caller source path")
	}
	dir := filepath.Dir(file)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("infra: go.mod not found walking up from %s", filepath.Dir(file))
		}
		dir = parent
	}
}

// Up brings a compose service up in the background using the given compose file.
func (dockerCompose) Up(composeFile, service string) error {
	return runCompose(composeFile, "up", "-d", service)
}

// Stop stops a compose service using the given compose file.
func (dockerCompose) Stop(composeFile, service string) error {
	return runCompose(composeFile, "stop", service)
}

// runCompose shells out to `docker compose -f <file> <args...>`, folding the combined
// output into the error on a non-zero exit.
func runCompose(composeFile string, args ...string) error {
	full := append([]string{"compose", "-f", composeFile}, args...)
	cmd := exec.Command("docker", full...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("docker %v: %w: %s", full, err, string(out))
	}
	return nil
}

// compose is the runner infra uses; a package var so tests inject a fake and the unit
// tests never touch docker.
var compose composeRunner = dockerCompose{}
