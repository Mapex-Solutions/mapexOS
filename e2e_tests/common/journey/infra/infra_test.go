package infra

import (
	"reflect"
	"testing"
)

// fakeCompose records Up/Stop calls (by compose service name) and marks a service
// "upped" so the fake probe reports it ready afterwards — simulating a real bring-up.
type fakeCompose struct {
	ups   []string
	stops []string
	upped map[string]bool
}

func (f *fakeCompose) Up(_, service string) error {
	f.ups = append(f.ups, service)
	if f.upped == nil {
		f.upped = map[string]bool{}
	}
	f.upped[service] = true
	return nil
}

func (f *fakeCompose) Stop(_, service string) error {
	f.stops = append(f.stops, service)
	return nil
}

// fakeSpawner records start/stop calls (by logical service name) and marks a service
// ready once started, so the fake probe reports it up — simulating a local `go run`.
type fakeSpawner struct {
	started []string
	stopped []string
	ready   map[string]bool
}

func (f *fakeSpawner) start(service string, _ localSpec) (runningProc, error) {
	f.started = append(f.started, service)
	if f.ready == nil {
		f.ready = map[string]bool{}
	}
	f.ready[service] = true
	return &fakeProc{svc: service, owner: f}, nil
}

type fakeProc struct {
	svc   string
	owner *fakeSpawner
}

func (p *fakeProc) stop() error {
	p.owner.stopped = append(p.owner.stopped, p.svc)
	return nil
}

// withFakes swaps probe + compose + spawner for the test. preReady lists services
// already up (by logical name); everything else is down until its compose service is
// upped or it is spawned. Uses the real serviceMap so the mapping is not duplicated.
func withFakes(t *testing.T, preReady ...string) (*fakeCompose, *fakeSpawner) {
	t.Helper()
	ready := map[string]bool{}
	for _, s := range preReady {
		ready[s] = true
	}
	fc, fs := &fakeCompose{}, &fakeSpawner{}
	origProbe, origCompose, origSpawner := probe, compose, spawner
	t.Cleanup(func() { probe, compose, spawner = origProbe, origCompose, origSpawner })
	compose, spawner = fc, fs
	probe = func(svc string) (bool, bool) {
		e, ok := serviceMap[svc]
		if !ok {
			return false, false
		}
		return ready[svc] || fc.upped[e.compose] || fs.ready[svc], true
	}
	return fc, fs
}

// dockerStack builds a docker-mode stack over the given logical services for tests.
func dockerStack(services ...string) stack {
	return stack{name: "test", mode: modeDocker, composeFile: func() string { return "test-compose.yml" }, services: services}
}

// localStack builds a local-mode stack over the given logical services for tests.
func localStack(services ...string) stack {
	return stack{name: "mapex", mode: modeLocal, composeFile: func() string { return "x" }, services: services}
}

func TestEnsureStack_AllReady_NoOp(t *testing.T) {
	fc, fs := withFakes(t, "assets", "broker")
	teardown := ensureStack(dockerStack("assets", "broker"))
	if len(fc.ups) != 0 || len(fs.started) != 0 {
		t.Fatalf("Up=%v started=%v, want none (both already up)", fc.ups, fs.started)
	}
	teardown()
	if len(fc.stops) != 0 || len(fs.stopped) != 0 {
		t.Fatalf("Stop=%v stopped=%v, want none (nothing owned)", fc.stops, fs.stopped)
	}
}

func TestEnsureStack_OneDown_DockerBringsUpAndStopsOnlyIt(t *testing.T) {
	fc, _ := withFakes(t, "assets") // assets up, broker down
	teardown := ensureStack(dockerStack("assets", "broker"))
	if !reflect.DeepEqual(fc.ups, []string{"mapex-broker-mqtt"}) {
		t.Fatalf("compose Up = %v, want [mapex-broker-mqtt]", fc.ups)
	}
	teardown()
	if !reflect.DeepEqual(fc.stops, []string{"mapex-broker-mqtt"}) {
		t.Fatalf("compose Stop = %v, want [mapex-broker-mqtt] (assets reused, not stopped)", fc.stops)
	}
}

func TestEnsureStack_LocalMode_SpawnsDownServiceAndStopsIt(t *testing.T) {
	fc, fs := withFakes(t) // nothing up
	teardown := ensureStack(localStack("assets"))
	if len(fc.ups) != 0 {
		t.Fatalf("compose Up = %v, want none in local mode (spawn, not compose)", fc.ups)
	}
	if !reflect.DeepEqual(fs.started, []string{"assets"}) {
		t.Fatalf("spawner started = %v, want [assets]", fs.started)
	}
	teardown()
	if !reflect.DeepEqual(fs.stopped, []string{"assets"}) {
		t.Fatalf("spawner stopped = %v, want [assets]", fs.stopped)
	}
	if len(fc.stops) != 0 {
		t.Fatalf("compose Stop = %v, want none", fc.stops)
	}
}

func TestEnsureStack_LocalMode_ReusedServiceNotSpawned(t *testing.T) {
	_, fs := withFakes(t, "assets") // assets already up
	teardown := ensureStack(localStack("assets"))
	if len(fs.started) != 0 {
		t.Fatalf("spawner started = %v, want none (assets reused)", fs.started)
	}
	teardown()
	if len(fs.stopped) != 0 {
		t.Fatalf("spawner stopped = %v, want none (nothing spawned)", fs.stopped)
	}
}

func TestEnsureStack_LocalMode_NoLocalSpec_NotSpawned(t *testing.T) {
	_, fs := withFakes(t) // nothing up
	// "nats" has no local run spec (docker-only infra service) — a down one in local
	// mode is flagged, never spawned.
	teardown := ensureStack(localStack("nats"))
	if len(fs.started) != 0 {
		t.Fatalf("spawner started = %v, want none (nats has no local spec)", fs.started)
	}
	teardown()
}

func TestEnvMode(t *testing.T) {
	if got := envMode("NOPE_UNSET", modeLocal); got != modeLocal {
		t.Errorf("unset = %q, want default local", got)
	}
	t.Setenv("MAPEX_TEST_MODE", "docker")
	if got := envMode("MAPEX_TEST_MODE", modeLocal); got != modeDocker {
		t.Errorf("docker env = %q, want docker", got)
	}
	t.Setenv("MAPEX_TEST_MODE", "garbage")
	if got := envMode("MAPEX_TEST_MODE", modeDocker); got != modeDocker {
		t.Errorf("garbage env = %q, want default docker", got)
	}
}
