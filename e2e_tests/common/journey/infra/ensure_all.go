package infra

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

// provMode is how a stack's DOWN services are handled.
type provMode string

const (
	modeDocker provMode = "docker" // bring a down service up via docker compose
	modeLocal  provMode = "local"  // the service runs from source (go run) — EnsureAll only reuses it; a down one is a clear "start it" error, never a compose up
)

// stack is one provisioning group: its logical services, the compose file that
// defines them, and the mode for its down services.
type stack struct {
	name        string
	mode        provMode
	composeFile func() string
	services    []string
}

// Option configures a single EnsureAll run — the suite-wide knobs the runner's
// TestMain declares once (e.g. per-service env). Options never mutate global state.
type Option func(*ensureConfig)

// ensureConfig accumulates the options applied to one EnsureAll run.
type ensureConfig struct {
	// env holds per-service extra KEY=VALUE entries (keyed by logical service
	// name), applied to a COPY of that service's localSpec at spawn.
	env map[string][]string
}

// WithServiceEnv declares extra KEY=VALUE environment for a locally-spawned
// service, e.g. WithServiceEnv("assets", "OTA_SCAN_INTERVAL=5"). Suite-wide by
// design (one shared stack serves every journey) and LOCAL mode only — docker
// services read their env from the compose file.
func WithServiceEnv(service string, kv ...string) Option {
	return func(c *ensureConfig) { c.env[service] = append(c.env[service], kv...) }
}

// EnsureAll brings up (or reuses) BOTH stacks once for the runner's TestMain:
//
//	infra  — nats / broker / minio        (docker by default)
//	mapex  — the app services             (local / `go run` by default)
//
// Each stack's mode is env-overridable (MAPEX_E2E_INFRA_MODE / MAPEX_E2E_MAPEX_MODE =
// docker|local). A service already answering its health probe is ALWAYS reused
// (skipped); only a down service is acted on — brought up via that stack's compose
// file (docker mode) or reported as a clear "start it" error (local mode). It uses
// log.* because TestMain has no *testing.T. Returns a teardown that stops ONLY the
// services this run started (docker mode), in reverse.
func EnsureAll(opts ...Option) func() {
	cfg := &ensureConfig{env: map[string][]string{}}
	for _, opt := range opts {
		opt(cfg)
	}
	stacks := []stack{
		{"infra", envMode("MAPEX_E2E_INFRA_MODE", modeDocker), infraComposeFile, infraServices},
		{"mapex", envMode("MAPEX_E2E_MAPEX_MODE", modeLocal), servicesComposeFile, mapexServices},
	}
	var teardowns []func()
	for _, st := range stacks {
		teardowns = append(teardowns, ensureStack(st, cfg.env))
	}
	return func() {
		for i := len(teardowns) - 1; i >= 0; i-- {
			teardowns[i]()
		}
	}
}

// ensureStack provisions one stack: reuse what is up, then bring DOWN services up —
// docker mode via compose, local mode by spawning `go run`/`npm run` from source. It
// walks the WHOLE list (a failure does not skip the rest), starts every down service
// FIRST and waits for readiness AFTER (so cold `go run` compiles and npm/ts-node boots
// warm up concurrently instead of one-bring-up-per-wait), logs one
// reused/brought-up/failed summary, and returns a teardown that stops only what it
// started — spawned processes then docker-owned services, all in reverse.
func ensureStack(st stack, envOverrides map[string][]string) func() {
	var owned, reused, up, failed, errs, pending []string
	var spawned []runningProc
	for _, svc := range st.services {
		ready, known := probe(svc)
		switch {
		case !known:
			failed, errs = append(failed, svc), append(errs, fmt.Sprintf("unknown service %q", svc))
		case ready:
			reused = append(reused, svc) // already answering — always skip
		case st.mode == modeLocal:
			base := serviceMap[svc].local
			if base == nil {
				failed, errs = append(failed, svc), append(errs, fmt.Sprintf("%q not responding — %s stack is LOCAL and has no local run spec; start it yourself", svc, st.name))
				continue
			}
			// Copy the spec and layer any per-run env override so the shared
			// serviceMap is never mutated (append into a fresh slice, not base.env).
			spec := *base
			if extra := envOverrides[svc]; len(extra) > 0 {
				spec.env = append(append([]string{}, base.env...), extra...)
			}
			p, e := spawner.start(svc, spec)
			if e != nil {
				failed, errs = append(failed, svc), append(errs, fmt.Sprintf("spawn %q: %v", svc, e))
				continue
			}
			spawned, pending = append(spawned, p), append(pending, svc)
		default: // docker: bring it up via this stack's compose file
			if e := compose.Up(st.composeFile(), serviceMap[svc].compose); e != nil {
				failed, errs = append(failed, svc), append(errs, fmt.Sprintf("bring up %q: %v", svc, e))
				continue
			}
			owned, pending = append(owned, svc), append(pending, svc)
		}
	}

	// Wait for everything started this run (spawned OR docker) to answer its probe.
	for _, svc := range pending {
		if e := waitReady(context.Background(), svc); e != nil {
			failed, errs = append(failed, svc), append(errs, e.Error())
			continue
		}
		up = append(up, svc)
	}
	log.Printf("[INFRA] %s stack (%s): reused=%v brought-up=%v failed=%v", st.name, st.mode, reused, up, failed)
	if len(errs) > 0 {
		log.Printf("[INFRA] %s stack issues: %s", st.name, strings.Join(errs, "; "))
	}

	cf := st.composeFile()
	return func() {
		for i := len(spawned) - 1; i >= 0; i-- {
			if e := spawned[i].stop(); e != nil {
				log.Printf("[INFRA] stop spawned failed (non-fatal): %v", e)
			}
		}
		for i := len(owned) - 1; i >= 0; i-- {
			if e := compose.Stop(cf, serviceMap[owned[i]].compose); e != nil {
				log.Printf("[INFRA] stop %q failed (non-fatal): %v", owned[i], e)
			}
		}
	}
}

// envMode reads a stack mode from env, defaulting to def; only "docker"/"local" are
// honored (anything else falls back to def).
func envMode(key string, def provMode) provMode {
	switch provMode(os.Getenv(key)) {
	case modeDocker:
		return modeDocker
	case modeLocal:
		return modeLocal
	default:
		return def
	}
}
