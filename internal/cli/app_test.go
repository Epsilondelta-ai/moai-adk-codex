package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestInitCreatesScaffoldAndManifest(t *testing.T) {
	dir := t.TempDir()
	stdout := &bytes.Buffer{}

	if err := Run([]string{"init", ".", "--json"}, stdout, &bytes.Buffer{}, func() (string, error) { return dir, nil }); err != nil {
		t.Fatalf("Run(init): %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}

	if payload["command"] != "init" {
		t.Fatalf("command = %v, want init", payload["command"])
	}
	if _, err := os.Stat(filepath.Join(dir, ".moai", "manifest.json")); err != nil {
		t.Fatalf("manifest missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "AGENTS.md")); err != nil {
		t.Fatalf("AGENTS.md missing: %v", err)
	}
}

func TestStatusReflectsInitializedScaffold(t *testing.T) {
	dir := t.TempDir()
	if err := Run([]string{"init", "."}, &bytes.Buffer{}, &bytes.Buffer{}, func() (string, error) { return dir, nil }); err != nil {
		t.Fatalf("Run(init): %v", err)
	}

	stdout := &bytes.Buffer{}
	if err := Run([]string{"status", "--json"}, stdout, &bytes.Buffer{}, func() (string, error) { return dir, nil }); err != nil {
		t.Fatalf("Run(status): %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if payload["initialized"] != true {
		t.Fatalf("initialized = %v, want true", payload["initialized"])
	}
	if payload["runtimeMode"] != "cg" {
		t.Fatalf("runtimeMode = %v, want cg", payload["runtimeMode"])
	}
}

func TestInitRespectsAbsoluteTargetPath(t *testing.T) {
	dir := t.TempDir()
	stdout := &bytes.Buffer{}

	if err := Run([]string{"init", dir, "--json"}, stdout, &bytes.Buffer{}, func() (string, error) {
		return t.TempDir(), nil
	}); err != nil {
		t.Fatalf("Run(init absolute): %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if payload["projectRoot"] != dir {
		t.Fatalf("projectRoot = %v, want %s", payload["projectRoot"], dir)
	}
}

func TestUpdatePreservesUserModifiedFiles(t *testing.T) {
	dir := t.TempDir()
	if err := Run([]string{"init", "."}, &bytes.Buffer{}, &bytes.Buffer{}, func() (string, error) { return dir, nil }); err != nil {
		t.Fatalf("Run(init): %v", err)
	}

	agentsPath := filepath.Join(dir, "AGENTS.md")
	if err := os.WriteFile(agentsPath, []byte("# custom\n"), 0o644); err != nil {
		t.Fatalf("os.WriteFile: %v", err)
	}

	stdout := &bytes.Buffer{}
	if err := Run([]string{"update", ".", "--json"}, stdout, &bytes.Buffer{}, func() (string, error) { return dir, nil }); err != nil {
		t.Fatalf("Run(update): %v", err)
	}

	var payload struct {
		Skipped []string `json:"skipped"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	found := false
	for _, value := range payload.Skipped {
		if value == "AGENTS.md" {
			found = true
		}
	}
	if !found {
		t.Fatalf("AGENTS.md was not skipped")
	}
}

func TestPlanCreatesSpecArtifact(t *testing.T) {
	dir := t.TempDir()
	if err := Run([]string{"init", "."}, &bytes.Buffer{}, &bytes.Buffer{}, func() (string, error) { return dir, nil }); err != nil {
		t.Fatalf("Run(init): %v", err)
	}

	stdout := &bytes.Buffer{}
	if err := Run([]string{"plan", "Add", "compatibility", "routing", "--json"}, stdout, &bytes.Buffer{}, func() (string, error) { return dir, nil }); err != nil {
		t.Fatalf("Run(plan): %v", err)
	}

	var payload struct {
		SpecID   string `json:"specId"`
		SpecPath string `json:"specPath"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if payload.SpecID == "" {
		t.Fatalf("specId was empty")
	}
	if _, err := os.Stat(filepath.Join(dir, payload.SpecPath)); err != nil {
		t.Fatalf("spec file missing: %v", err)
	}
}

func TestRunCreatesWorkflowArtifact(t *testing.T) {
	dir := t.TempDir()
	if err := Run([]string{"init", "."}, &bytes.Buffer{}, &bytes.Buffer{}, func() (string, error) { return dir, nil }); err != nil {
		t.Fatalf("Run(init): %v", err)
	}

	stdout := &bytes.Buffer{}
	if err := Run([]string{"run", "SPEC-TEST-001", "--json"}, stdout, &bytes.Buffer{}, func() (string, error) { return dir, nil }); err != nil {
		t.Fatalf("Run(run): %v", err)
	}

	var payload struct {
		SpecID       string `json:"specId"`
		ArtifactPath string `json:"artifactPath"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if payload.SpecID != "SPEC-TEST-001" {
		t.Fatalf("specId = %s, want SPEC-TEST-001", payload.SpecID)
	}
	if _, err := os.Stat(filepath.Join(dir, payload.ArtifactPath)); err != nil {
		t.Fatalf("artifact missing: %v", err)
	}
}

func TestModeSwitchUpdatesRuntimeState(t *testing.T) {
	dir := t.TempDir()
	if err := Run([]string{"init", "."}, &bytes.Buffer{}, &bytes.Buffer{}, func() (string, error) { return dir, nil }); err != nil {
		t.Fatalf("Run(init): %v", err)
	}
	if err := Run([]string{"glm"}, &bytes.Buffer{}, &bytes.Buffer{}, func() (string, error) { return dir, nil }); err != nil {
		t.Fatalf("Run(glm): %v", err)
	}

	stdout := &bytes.Buffer{}
	if err := Run([]string{"status", "--json"}, stdout, &bytes.Buffer{}, func() (string, error) { return dir, nil }); err != nil {
		t.Fatalf("Run(status): %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if payload["runtimeMode"] != "glm" {
		t.Fatalf("runtimeMode = %v, want glm", payload["runtimeMode"])
	}
}

func TestDoctorReportsScaffoldAvailability(t *testing.T) {
	dir := t.TempDir()
	if err := Run([]string{"init", "."}, &bytes.Buffer{}, &bytes.Buffer{}, func() (string, error) { return dir, nil }); err != nil {
		t.Fatalf("Run(init): %v", err)
	}

	stdout := &bytes.Buffer{}
	if err := Run([]string{"doctor", "--json"}, stdout, &bytes.Buffer{}, func() (string, error) { return dir, nil }); err != nil {
		t.Fatalf("Run(doctor): %v", err)
	}

	var payload struct {
		Checks []struct {
			Name string `json:"name"`
			OK   bool   `json:"ok"`
		} `json:"checks"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}

	found := false
	for _, check := range payload.Checks {
		if check.Name == ".moai scaffold" && check.OK {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected .moai scaffold check to pass")
	}
}

func TestStatusDoesNotTreatHomeMoaiAsProjectRoot(t *testing.T) {
	dir := t.TempDir()
	stdout := &bytes.Buffer{}
	if err := Run([]string{"status", "--json"}, stdout, &bytes.Buffer{}, func() (string, error) { return dir, nil }); err != nil {
		t.Fatalf("Run(status): %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if payload["projectRoot"] != dir {
		t.Fatalf("projectRoot = %v, want %s", payload["projectRoot"], dir)
	}
	if payload["initialized"] != false {
		t.Fatalf("initialized = %v, want false", payload["initialized"])
	}
}

func TestProjectAndCodemapsRefreshDocs(t *testing.T) {
	dir := t.TempDir()
	if err := Run([]string{"init", "."}, &bytes.Buffer{}, &bytes.Buffer{}, func() (string, error) { return dir, nil }); err != nil {
		t.Fatalf("Run(init): %v", err)
	}
	if err := Run([]string{"project", "Compatibility", "docs", "for", "the", "repo"}, &bytes.Buffer{}, &bytes.Buffer{}, func() (string, error) { return dir, nil }); err != nil {
		t.Fatalf("Run(project): %v", err)
	}
	if err := Run([]string{"codemaps"}, &bytes.Buffer{}, &bytes.Buffer{}, func() (string, error) { return dir, nil }); err != nil {
		t.Fatalf("Run(codemaps): %v", err)
	}

	product, err := os.ReadFile(filepath.Join(dir, ".moai", "project", "product.md"))
	if err != nil {
		t.Fatalf("os.ReadFile(product): %v", err)
	}
	if !bytes.Contains(product, []byte("Compatibility docs")) {
		t.Fatalf("product.md did not contain updated content")
	}
	if _, err := os.Stat(filepath.Join(dir, ".moai", "project", "codemaps", "overview.md")); err != nil {
		t.Fatalf("overview.md missing: %v", err)
	}
}
