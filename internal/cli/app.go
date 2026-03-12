package cli

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/Epsilondelta-ai/moai-adk-codex/internal/core"
	"github.com/Epsilondelta-ai/moai-adk-codex/internal/workflows"
)

const HelpText = `MoAI Codex Compatibility CLI

Usage:
  moai-codex <command> [args]

Core commands:
  init [path]        Initialize a Codex-compatible MoAI scaffold
  update [path]      Reapply managed scaffold files
  status             Show scaffold and mode status
  doctor             Check runtime and scaffold health

Workflow commands:
  project [summary]
  plan <description>
  run <spec-id|description>
  sync [spec-id]
  review
  coverage
  clean
  fix
  loop
  codemaps

Mode commands:
  cc | cg | glm

Git worktree:
  worktree list
  worktree status
  worktree new <name>
  worktree remove <path>

Flags:
  --json       Render machine-readable output where supported
  --execute    For workflow commands, also invoke codex exec
`

type getwdFunc func() (string, error)

type parsedArgs struct {
	Command string
	Args    []string
	Flags   flags
}

type flags struct {
	JSON    bool
	Execute bool
}

func Run(argv []string, stdout io.Writer, _ io.Writer, getwd getwdFunc) error {
	parsed := parseArgs(argv)
	cwd, err := getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	switch parsed.Command {
	case "", "help", "--help", "-h":
		_, _ = io.WriteString(stdout, HelpText+"\n")
		return nil
	case "version", "--version", "-V":
		_, _ = io.WriteString(stdout, "moai-codex 0.2.0\n")
		return nil
	case "init":
		target, err := resolveTargetPath(cwd, firstOrDefault(parsed.Args, "."))
		if err != nil {
			return err
		}
		result, err := core.EnsureProjectScaffold(target, false)
		if err != nil {
			return err
		}
		return writeOutput(stdout, parsed.Flags.JSON, result)
	case "update":
		target, err := resolveTargetPath(cwd, firstOrDefault(parsed.Args, "."))
		if err != nil {
			return err
		}
		result, err := core.EnsureProjectScaffold(target, true)
		if err != nil {
			return err
		}
		return writeOutput(stdout, parsed.Flags.JSON, result)
	case "status":
		result, err := core.ReadProjectStatus(cwd)
		if err != nil {
			return err
		}
		return writeOutput(stdout, parsed.Flags.JSON, result)
	case "doctor":
		result := core.RunDoctor(cwd)
		return writeOutput(stdout, parsed.Flags.JSON, result)
	case "cc", "cg", "glm":
		result, err := core.SetRuntimeMode(cwd, parsed.Command)
		if err != nil {
			return err
		}
		return writeOutput(stdout, parsed.Flags.JSON, result)
	case "project", "sync", "review", "coverage", "clean", "fix", "loop", "codemaps":
		if _, err := core.EnsureProjectScaffold(cwd, false); err != nil {
			return err
		}
		result, err := workflows.CreateWorkflowArtifact(cwd, parsed.Command, strings.TrimSpace(strings.Join(parsed.Args, " ")), "")
		if err != nil {
			return err
		}
		if parsed.Flags.Execute {
			execResult := workflows.MaybeExecuteWithCodex(cwd, parsed.Command, result.Prompt)
			result.Execution = &execResult
		}
		return writeOutput(stdout, parsed.Flags.JSON, result)
	case "plan":
		if len(parsed.Args) == 0 {
			return fmt.Errorf("plan requires a description")
		}
		if _, err := core.EnsureProjectScaffold(cwd, false); err != nil {
			return err
		}
		result, err := workflows.CreateSpec(cwd, strings.Join(parsed.Args, " "))
		if err != nil {
			return err
		}
		if parsed.Flags.Execute {
			execResult := workflows.MaybeExecuteWithCodex(cwd, "plan", result.Prompt)
			result.Execution = &execResult
		}
		return writeOutput(stdout, parsed.Flags.JSON, result)
	case "run":
		if len(parsed.Args) == 0 {
			return fmt.Errorf("run requires a SPEC id or description")
		}
		if _, err := core.EnsureProjectScaffold(cwd, false); err != nil {
			return err
		}
		raw := strings.TrimSpace(strings.Join(parsed.Args, " "))
		specID := raw
		if !strings.HasPrefix(specID, "SPEC-") {
			specID = "SPEC-" + strings.ToUpper(core.Slugify(raw)) + "-001"
		}
		result, err := workflows.CreateWorkflowArtifact(cwd, "run", raw, specID)
		if err != nil {
			return err
		}
		if parsed.Flags.Execute {
			execResult := workflows.MaybeExecuteWithCodex(cwd, "run", result.Prompt)
			result.Execution = &execResult
		}
		return writeOutput(stdout, parsed.Flags.JSON, result)
	case "worktree":
		result, err := workflows.RunWorktreeCommand(cwd, parsed.Args)
		if err != nil {
			return err
		}
		return writeOutput(stdout, parsed.Flags.JSON, result)
	default:
		return fmt.Errorf("unknown command: %s", parsed.Command)
	}
}

func parseArgs(argv []string) parsedArgs {
	out := parsedArgs{Flags: flags{}, Args: []string{}}
	for _, value := range argv {
		switch value {
		case "--json":
			out.Flags.JSON = true
		case "--execute":
			out.Flags.Execute = true
		default:
			if out.Command == "" {
				out.Command = value
			} else {
				out.Args = append(out.Args, value)
			}
		}
	}
	return out
}

func writeOutput(stdout io.Writer, asJSON bool, payload any) error {
	var rendered string
	if asJSON {
		rendered = core.RenderJSON(payload)
	} else {
		rendered = core.RenderText(payload)
	}
	_, err := io.WriteString(stdout, rendered+"\n")
	return err
}

func firstOrDefault(values []string, fallback string) string {
	if len(values) == 0 {
		return fallback
	}
	return values[0]
}

func resolveTargetPath(cwd, target string) (string, error) {
	if filepath.IsAbs(target) {
		return filepath.Abs(target)
	}
	return filepath.Abs(filepath.Join(cwd, target))
}
