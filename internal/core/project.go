package core

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	primaryControlDir = ".coai"
	legacyControlDir  = ".moai"
)

type runtimeState struct {
	CurrentRuntimeMode string `json:"currentRuntimeMode"`
	CurrentSpec        string `json:"currentSpec"`
	LastCommand        string `json:"lastCommand"`
	UpdatedAt          string `json:"updatedAt,omitempty"`
}

func EnsureProjectScaffold(projectRoot string, forceUpdate bool) (ScaffoldResult, error) {
	root, err := filepath.Abs(projectRoot)
	if err != nil {
		return ScaffoldResult{}, err
	}
	if err := EnsureDir(root); err != nil {
		return ScaffoldResult{}, err
	}

	manifestPath := filepath.Join(primaryControlRoot(root), "manifest.json")
	manifest, err := loadOrMigrateManifest(root)
	if err != nil {
		return ScaffoldResult{}, err
	}

	created := []string{}
	updated := []string{}
	skipped := []string{}

	for _, template := range BuildTemplates(filepath.Base(root), time.Now().UTC().Format(time.RFC3339)) {
		destination := filepath.Join(root, template.Path)
		hash := SHA256(template.Content)
		entry, hasEntry := manifest.Files[template.Path]

		if err := EnsureDir(filepath.Dir(destination)); err != nil {
			return ScaffoldResult{}, err
		}

		currentContent, err := os.ReadFile(destination)
		if err != nil {
			if os.IsNotExist(err) {
				if err := os.WriteFile(destination, []byte(template.Content), 0o644); err != nil {
					return ScaffoldResult{}, err
				}
				manifest.Files[template.Path] = ManifestEntry{Managed: true, Hash: hash}
				created = append(created, template.Path)
				continue
			}
			return ScaffoldResult{}, err
		}

		currentHash := SHA256(string(currentContent))
		if forceUpdate && hasEntry && entry.Hash == currentHash {
			if err := os.WriteFile(destination, []byte(template.Content), 0o644); err != nil {
				return ScaffoldResult{}, err
			}
			manifest.Files[template.Path] = ManifestEntry{Managed: true, Hash: hash}
			updated = append(updated, template.Path)
			continue
		}

		if !hasEntry {
			manifest.Files[template.Path] = ManifestEntry{Managed: false, Hash: currentHash}
			skipped = append(skipped, template.Path)
			continue
		}

		if entry.Hash != currentHash {
			skipped = append(skipped, template.Path)
			continue
		}

		manifest.Files[template.Path] = ManifestEntry{Managed: true, Hash: hash}
	}

	if err := WriteJSON(manifestPath, manifest); err != nil {
		return ScaffoldResult{}, err
	}

	command := "init"
	if forceUpdate {
		command = "update"
	}

	return ScaffoldResult{
		Command:     command,
		ProjectRoot: root,
		Created:     created,
		Updated:     updated,
		Skipped:     skipped,
		Manifest:    filepath.ToSlash(filepath.Join(primaryControlDir, "manifest.json")),
	}, nil
}

func ReadProjectStatus(projectRoot string) (ProjectStatus, error) {
	root := FindProjectRoot(projectRoot)
	controlRoot, source := controlRootForRead(root)

	state := runtimeState{CurrentRuntimeMode: "cg"}
	_ = ReadJSON(filepath.Join(controlRoot, "state", "runtime.json"), &state)
	projectConfig := ReadSimpleYAML(filepath.Join(controlRoot, "config", "sections", "project.yaml"))
	qualityConfig := ReadSimpleYAML(filepath.Join(controlRoot, "config", "sections", "quality.yaml"))
	manifest := Manifest{Files: map[string]ManifestEntry{}}
	_ = ReadJSON(filepath.Join(controlRoot, "manifest.json"), &manifest)

	projectName := filepath.Base(root)
	if projectSection, ok := projectConfig["project"].(map[string]any); ok {
		if name, ok := projectSection["name"].(string); ok && name != "" {
			projectName = name
		}
	}

	developmentMode := "tdd"
	if constitution, ok := qualityConfig["constitution"].(map[string]any); ok {
		if mode, ok := constitution["development_mode"].(string); ok && mode != "" {
			developmentMode = mode
		}
	}

	return ProjectStatus{
		ProjectRoot:     root,
		ProjectName:     projectName,
		RuntimeMode:     defaultString(state.CurrentRuntimeMode, "cg"),
		CurrentSpec:     state.CurrentSpec,
		LastCommand:     state.LastCommand,
		DevelopmentMode: developmentMode,
		ManagedFiles:    len(manifest.Files),
		Initialized:     source != "",
	}, nil
}

func SetRuntimeMode(projectRoot, mode string) (ModeResult, error) {
	root := FindProjectRoot(projectRoot)
	if err := EnsureDir(filepath.Join(primaryControlRoot(root), "state")); err != nil {
		return ModeResult{}, err
	}
	if err := EnsureDir(filepath.Join(primaryControlRoot(root), "config", "sections")); err != nil {
		return ModeResult{}, err
	}

	state := runtimeState{}
	_ = ReadJSON(filepath.Join(primaryControlRoot(root), "state", "runtime.json"), &state)
	state.CurrentRuntimeMode = mode
	state.LastCommand = mode
	if err := WriteJSON(filepath.Join(primaryControlRoot(root), "state", "runtime.json"), state); err != nil {
		return ModeResult{}, err
	}

	if err := os.WriteFile(filepath.Join(primaryControlRoot(root), "config", "sections", "llm.yaml"), []byte("llm:\n  current_runtime_mode: "+mode+"\n  provider: codex\n"), 0o644); err != nil {
		return ModeResult{}, err
	}

	return ModeResult{
		Command:     mode,
		RuntimeMode: mode,
		ProjectRoot: root,
	}, nil
}

func RunDoctor(projectRoot string) DoctorReport {
	root, _ := filepath.Abs(projectRoot)
	scaffoldOK := dirExists(primaryControlRoot(root)) || dirExists(legacyControlRoot(root))
	scaffoldDetails := "missing"
	switch {
	case dirExists(primaryControlRoot(root)):
		scaffoldDetails = "present (.coai)"
	case dirExists(legacyControlRoot(root)):
		scaffoldDetails = "present (.moai legacy)"
	}
	checks := []Check{
		checkBinary("go", "Go runtime"),
		checkBinary("git", "Git"),
		checkBinary("codex", "Codex CLI"),
		{
			Name:    ".coai scaffold",
			OK:      scaffoldOK,
			Details: scaffoldDetails,
		},
	}

	ok := true
	for _, check := range checks {
		if !check.OK {
			ok = false
		}
	}
	return DoctorReport{ProjectRoot: root, OK: ok, Checks: checks}
}

func FindProjectRoot(startDir string) string {
	initial, _ := filepath.Abs(startDir)
	current := initial
	homeDir, _ := os.UserHomeDir()
	homeDir, _ = filepath.Abs(homeDir)

	for {
		if current != homeDir && (dirExists(primaryControlRoot(current)) || dirExists(legacyControlRoot(current))) {
			return current
		}
		parent := filepath.Dir(current)
		if parent == current {
			return initial
		}
		current = parent
	}
}

func UpdateRuntime(projectRoot string, patch RuntimePatch) error {
	root := FindProjectRoot(projectRoot)
	filePath := filepath.Join(primaryControlRoot(root), "state", "runtime.json")
	if err := EnsureDir(filepath.Dir(filePath)); err != nil {
		return err
	}
	current := runtimeState{}
	_ = ReadJSON(filePath, &current)

	if patch.CurrentRuntimeMode != "" {
		current.CurrentRuntimeMode = patch.CurrentRuntimeMode
	}
	current.CurrentSpec = patch.CurrentSpec
	if patch.LastCommand != "" {
		current.LastCommand = patch.LastCommand
	}
	current.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return WriteJSON(filePath, current)
}

func checkBinary(binary, name string) Check {
	if err := exec.Command("which", binary).Run(); err != nil {
		if binary == "go" {
			homeDir, homeErr := os.UserHomeDir()
			if homeErr == nil {
				localGo := filepath.Join(homeDir, ".local", "go", "bin", "go")
				if _, statErr := os.Stat(localGo); statErr == nil {
					return Check{Name: name, OK: true, Details: "available at ~/.local/go/bin/go"}
				}
			}
		}
		return Check{Name: name, OK: false, Details: "missing"}
	}
	return Check{Name: name, OK: true, Details: "available"}
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func ternary[T any](condition bool, left, right T) T {
	if condition {
		return left
	}
	return right
}

func primaryControlRoot(root string) string {
	return filepath.Join(root, primaryControlDir)
}

func legacyControlRoot(root string) string {
	return filepath.Join(root, legacyControlDir)
}

func controlRootForRead(root string) (string, string) {
	switch {
	case dirExists(primaryControlRoot(root)):
		return primaryControlRoot(root), primaryControlDir
	case dirExists(legacyControlRoot(root)):
		return legacyControlRoot(root), legacyControlDir
	default:
		return primaryControlRoot(root), ""
	}
}

func loadOrMigrateManifest(root string) (Manifest, error) {
	manifest := Manifest{Files: map[string]ManifestEntry{}}
	if err := ReadJSON(filepath.Join(primaryControlRoot(root), "manifest.json"), &manifest); err == nil {
		if manifest.Files == nil {
			manifest.Files = map[string]ManifestEntry{}
		}
		return manifest, nil
	}

	legacyManifest := Manifest{Files: map[string]ManifestEntry{}}
	if err := ReadJSON(filepath.Join(legacyControlRoot(root), "manifest.json"), &legacyManifest); err == nil {
		migrated := Manifest{Files: map[string]ManifestEntry{}}
		for key, entry := range legacyManifest.Files {
			migrated.Files[strings.Replace(key, legacyControlDir+"/", primaryControlDir+"/", 1)] = entry
		}
		return migrated, nil
	}

	return manifest, nil
}
