package infra

import (
	"path/filepath"
	"testing"
)

// TestComposeFiles_CWDIndependent proves both stack compose files resolve to the same
// absolute path regardless of the process working directory, and that the two stacks
// resolve to DISTINCT files.
func TestComposeFiles_CWDIndependent(t *testing.T) {
	infraBefore, svcBefore := infraComposeFile(), servicesComposeFile()

	t.Chdir(t.TempDir())
	if got := infraComposeFile(); got != infraBefore {
		t.Errorf("infraComposeFile is CWD-dependent: %q before chdir, %q after", infraBefore, got)
	}
	if got := servicesComposeFile(); got != svcBefore {
		t.Errorf("servicesComposeFile is CWD-dependent: %q before chdir, %q after", svcBefore, got)
	}
	if !filepath.IsAbs(infraBefore) || !filepath.IsAbs(svcBefore) {
		t.Errorf("compose files must be absolute: infra=%q services=%q", infraBefore, svcBefore)
	}
	if infraBefore == svcBefore {
		t.Errorf("infra and services compose files must differ, both = %q", infraBefore)
	}
}

// TestComposeFiles_EnvOverride proves each stack's env override wins verbatim.
func TestComposeFiles_EnvOverride(t *testing.T) {
	t.Setenv("MAPEX_COMPOSE_FILE", "/custom/infra.yml")
	t.Setenv("MAPEX_SERVICES_COMPOSE_FILE", "/custom/services.yml")

	if got := infraComposeFile(); got != "/custom/infra.yml" {
		t.Errorf("infraComposeFile = %q, want /custom/infra.yml", got)
	}
	if got := servicesComposeFile(); got != "/custom/services.yml" {
		t.Errorf("servicesComposeFile = %q, want /custom/services.yml", got)
	}
}
