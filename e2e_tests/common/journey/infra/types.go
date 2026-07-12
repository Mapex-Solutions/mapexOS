// Package infra provisions the live-stack services the suite needs: EnsureAll()
// brings up (or reuses) the whole known stack once for the runner's TestMain and
// returns a teardown that stops only the services it started. It is the e2e
// Environment layer — idempotent, and a no-op when the stack is already up.
package infra

// serviceEntry maps a logical service name to its docker-compose service name, a
// readiness probe, and (optionally) how to run it from source in local mode.
// serviceMap (services.go) holds one entry per supported service.
type serviceEntry struct {
	compose string
	ready   func() bool
	local   *localSpec // how to `go run`/`npm run` it in local mode; nil = not locally spawnable
}

// localSpec is the command that starts a service from source in local mode: the
// executable, its args, and a resolver for the working directory it runs in.
type localSpec struct {
	cmd  string
	args []string
	dir  func() string
}

// composeRunner brings a compose service up or stops it in a given compose file. It
// is a seam: the default impl shells out to docker compose; tests inject a fake so the
// unit tests need no docker.
type composeRunner interface {
	Up(composeFile, service string) error
	Stop(composeFile, service string) error
}

// dockerCompose is the default composeRunner — it shells out to `docker compose`.
type dockerCompose struct{}

// localSpawner starts a service from source and returns a handle that stops it. A seam:
// the default (osSpawner) runs a real process group; tests inject a fake.
type localSpawner interface {
	start(service string, spec localSpec) (runningProc, error)
}

// runningProc is a started local service; stop() terminates it (and its process group).
type runningProc interface {
	stop() error
}
