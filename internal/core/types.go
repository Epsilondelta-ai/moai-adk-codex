package core

type Manifest struct {
	Files map[string]ManifestEntry `json:"files"`
}

type ManifestEntry struct {
	Managed bool   `json:"managed"`
	Hash    string `json:"hash"`
}

type ScaffoldResult struct {
	Command     string   `json:"command"`
	ProjectRoot string   `json:"projectRoot"`
	Created     []string `json:"created"`
	Updated     []string `json:"updated"`
	Skipped     []string `json:"skipped"`
	Manifest    string   `json:"manifest"`
}

type ProjectStatus struct {
	ProjectRoot     string `json:"projectRoot"`
	ProjectName     string `json:"projectName"`
	RuntimeMode     string `json:"runtimeMode"`
	CurrentSpec     string `json:"currentSpec"`
	LastCommand     string `json:"lastCommand"`
	DevelopmentMode string `json:"developmentMode"`
	ManagedFiles    int    `json:"managedFiles"`
	Initialized     bool   `json:"initialized"`
}

type Check struct {
	Name    string `json:"name"`
	OK      bool   `json:"ok"`
	Details string `json:"details"`
}

type DoctorReport struct {
	ProjectRoot string  `json:"projectRoot"`
	OK          bool    `json:"ok"`
	Checks      []Check `json:"checks"`
}

type ModeResult struct {
	Command     string `json:"command"`
	RuntimeMode string `json:"runtimeMode"`
	ProjectRoot string `json:"projectRoot"`
}

type RuntimePatch struct {
	CurrentRuntimeMode string
	CurrentSpec        string
	LastCommand        string
}

type WorkflowExecution struct {
	Command  string `json:"command"`
	ExitCode int    `json:"exitCode"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
}

type WorkflowResult struct {
	Command      string             `json:"command"`
	ArtifactPath string             `json:"artifactPath"`
	Summary      string             `json:"summary"`
	SpecID       string             `json:"specId,omitempty"`
	Prompt       string             `json:"prompt"`
	Execution    *WorkflowExecution `json:"execution,omitempty"`
}

type SpecResult struct {
	Command   string             `json:"command"`
	SpecID    string             `json:"specId"`
	SpecPath  string             `json:"specPath"`
	Prompt    string             `json:"prompt"`
	Execution *WorkflowExecution `json:"execution,omitempty"`
}

type GitResult struct {
	Command string `json:"command"`
	OK      bool   `json:"ok"`
	Stdout  string `json:"stdout"`
	Stderr  string `json:"stderr"`
}
