package workflows

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Epsilondelta-ai/coai/internal/core"
)

func CreateSpec(projectRoot, description string) (core.SpecResult, error) {
	root := core.FindProjectRoot(projectRoot)
	coaiRoot := filepath.Join(root, ".coai")
	specID := "SPEC-" + strings.ToUpper(core.Slugify(description)) + "-001"
	filePath := filepath.Join(coaiRoot, "specs", specID+".md")

	if err := core.EnsureDir(filepath.Dir(filePath)); err != nil {
		return core.SpecResult{}, err
	}

	content := fmt.Sprintf(`# %s

## Summary

%s

## Acceptance Criteria

- Implementation is compatible with Codex/OMX.
- Changes are verifiable via tests or command outputs.
- Documentation is updated when behavior changes.
`, specID, description)

	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		return core.SpecResult{}, err
	}
	if err := core.UpdateRuntime(root, core.RuntimePatch{CurrentSpec: specID, LastCommand: "plan"}); err != nil {
		return core.SpecResult{}, err
	}

	return core.SpecResult{
		Command:  "plan",
		SpecID:   specID,
		SpecPath: filepath.ToSlash(strings.TrimPrefix(filePath, root+string(filepath.Separator))),
		Prompt:   fmt.Sprintf("Create or refine implementation for %s: %s", specID, description),
	}, nil
}

func CreateWorkflowArtifact(projectRoot, command, summary, specID string) (core.WorkflowResult, error) {
	root := core.FindProjectRoot(projectRoot)
	coaiRoot := filepath.Join(root, ".coai")
	reportsDir := filepath.Join(coaiRoot, "project", "reports")
	if err := core.EnsureDir(reportsDir); err != nil {
		return core.WorkflowResult{}, err
	}
	if err := core.EnsureDir(filepath.Join(coaiRoot, "project", "codemaps")); err != nil {
		return core.WorkflowResult{}, err
	}

	effectiveSummary := strings.TrimSpace(summary)
	if effectiveSummary == "" {
		effectiveSummary = defaultSummary(command)
	}

	if command == "project" {
		if err := writeProjectDocs(root, effectiveSummary); err != nil {
			return core.WorkflowResult{}, err
		}
	}
	if command == "codemaps" {
		content := fmt.Sprintf("# Codemaps\n\nUpdated at %s\n\nSummary: %s\n", time.Now().UTC().Format(time.RFC3339), effectiveSummary)
		if err := os.WriteFile(filepath.Join(coaiRoot, "project", "codemaps", "overview.md"), []byte(content), 0o644); err != nil {
			return core.WorkflowResult{}, err
		}
	}

	fileName := fmt.Sprintf("%s-%s.md", command, core.TimestampUTC())
	artifactPath := filepath.Join(reportsDir, fileName)
	prompt := buildPrompt(command, effectiveSummary, specID)
	content := fmt.Sprintf(`# %s

- Timestamp: %s
- Summary: %s
- Spec: %s
- Mode: compatibility

## Next Codex Prompt

%s
`, command, time.Now().UTC().Format(time.RFC3339), effectiveSummary, defaultString(specID, "n/a"), prompt)

	if err := os.WriteFile(artifactPath, []byte(content), 0o644); err != nil {
		return core.WorkflowResult{}, err
	}
	if err := core.UpdateRuntime(root, core.RuntimePatch{CurrentSpec: specID, LastCommand: command}); err != nil {
		return core.WorkflowResult{}, err
	}

	return core.WorkflowResult{
		Command:      command,
		ArtifactPath: filepath.ToSlash(strings.TrimPrefix(artifactPath, root+string(filepath.Separator))),
		Summary:      effectiveSummary,
		SpecID:       specID,
		Prompt:       prompt,
	}, nil
}

func MaybeExecuteWithCodex(projectRoot, command, prompt string) core.WorkflowExecution {
	root := core.FindProjectRoot(projectRoot)
	cmd := exec.Command("codex", "exec", "--dangerously-bypass-approvals-and-sandbox", "--cd", root, prompt)
	cmd.Dir = root
	output, err := cmd.CombinedOutput()
	if err != nil {
		exitCode := 1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
		return core.WorkflowExecution{
			Command:  command,
			ExitCode: exitCode,
			Stdout:   "",
			Stderr:   strings.TrimSpace(string(output)),
		}
	}
	return core.WorkflowExecution{
		Command:  command,
		ExitCode: 0,
		Stdout:   strings.TrimSpace(string(output)),
		Stderr:   "",
	}
}

func RunWorktreeCommand(projectRoot string, args []string) (core.GitResult, error) {
	if len(args) == 0 {
		return core.GitResult{}, fmt.Errorf("supported worktree commands: list, status, new, remove")
	}
	cwd, err := filepath.Abs(projectRoot)
	if err != nil {
		return core.GitResult{}, err
	}

	subcommand := args[0]
	switch subcommand {
	case "list":
		return execGit(cwd, []string{"worktree", "list", "--porcelain"}, "worktree list"), nil
	case "status":
		return execGit(cwd, []string{"worktree", "list"}, "worktree status"), nil
	case "new":
		if len(args) < 2 {
			return core.GitResult{}, fmt.Errorf("worktree new requires a name")
		}
		name := args[1]
		return execGit(cwd, []string{"worktree", "add", "../" + name, "-b", name}, "worktree new"), nil
	case "remove":
		if len(args) < 2 {
			return core.GitResult{}, fmt.Errorf("worktree remove requires a path")
		}
		return execGit(cwd, []string{"worktree", "remove", args[1]}, "worktree remove"), nil
	default:
		return core.GitResult{}, fmt.Errorf("supported worktree commands: list, status, new, remove")
	}
}

func execGit(cwd string, args []string, command string) core.GitResult {
	cmd := exec.Command("git", args...)
	cmd.Dir = cwd
	output, err := cmd.CombinedOutput()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			return core.GitResult{
				Command: command,
				OK:      false,
				Stdout:  "",
				Stderr:  strings.TrimSpace(string(exitErr.Stderr) + string(output)),
			}
		}
		return core.GitResult{Command: command, OK: false, Stdout: "", Stderr: strings.TrimSpace(string(output))}
	}
	return core.GitResult{Command: command, OK: true, Stdout: strings.TrimSpace(string(output)), Stderr: ""}
}

func writeProjectDocs(root, summary string) error {
	coaiRoot := filepath.Join(root, ".coai")
	if err := os.WriteFile(filepath.Join(coaiRoot, "project", "product.md"), []byte("# Product\n\n"+summary+"\n"), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(coaiRoot, "project", "structure.md"), []byte("# Structure\n\n- Generated by coai project workflow.\n"), 0o644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(coaiRoot, "project", "tech.md"), []byte("# Tech\n\n- Runtime: Codex CLI\n- Compatibility layer: coai\n"), 0o644)
}

func buildPrompt(command, summary, specID string) string {
	header := strings.ToUpper(command)
	if specID != "" {
		header += " " + specID
	}
	return fmt.Sprintf("%s\n\nSummary: %s\n\nUse Codex/OMX-compatible behavior and preserve .coai artifacts.", header, summary)
}

func defaultSummary(command string) string {
	summaries := map[string]string{
		"project":  "Refresh project-level compatibility docs.",
		"run":      "Advance the active SPEC using Codex-native execution.",
		"sync":     "Synchronize compatibility docs and workflow artifacts.",
		"review":   "Review current changes through the compatibility lens.",
		"coverage": "Assess test coverage and gaps.",
		"clean":    "Identify dead code and cleanup candidates.",
		"fix":      "Apply a single-pass repair loop.",
		"loop":     "Repeat repair and verification until convergence.",
		"codemaps": "Refresh architecture maps.",
	}
	if summary, ok := summaries[command]; ok {
		return summary
	}
	return command + " compatibility workflow"
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
