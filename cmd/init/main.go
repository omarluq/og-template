// Package main provides the interactive project initializer for og-template.
package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/huh/v2"
	"github.com/samber/oops"
	"github.com/samber/lo"
)

const (
	oldModule = "github.com/username/myapp"
	oldBinary = "myapp"
	oldPrefix = "MYAPP"
)

// allHarnesses lists all supported AI coding assistant harnesses.
var allHarnesses = []string{
	".adal",
	".augment",
	".claude",
	".codebuddy",
	".commandcode",
	".continue",
	".cortex",
	".crush",
	".factory",
	".goose",
	".iflow",
	".junie",
	".kilocode",
	".kiro",
	".kode",
	".mcpjam",
	".mux",
	".neovate",
	".openhands",
	".pi",
	".pochi",
	".qoder",
	".qwen",
	".roo",
	".trae",
	".vibe",
	".windsurf",
	".zencoder",
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	doneStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("76"))
	errStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))
)

// rewriteExts lists file extensions to rewrite.
var rewriteExts = map[string]bool{
	".go": true, ".yml": true, ".yaml": true,
	".mod": true, ".md": true,
}

// skipDirs lists directories to skip.
var skipDirs = map[string]bool{
	".git": true, ".gocache": true, ".gomodcache": true,
	".tmp": true, "bin": true, ".agents": true,
}

type replacement struct {
	label string
	old   string
	new   string
}

func main() {
	if err := run(); err != nil {
		writeOut(errStyle.Render(err.Error()))
		os.Exit(1)
	}
}

func run() error {
	writeOut(titleStyle.Render("og-template init"))
	writeOut("")

	module, binary, prefix, keep, promptErr := promptUser()
	if promptErr != nil {
		return promptErr
	}

	if module == "" {
		return nil
	}

	writeOut("")

	return applyAndFinalize(module, binary, prefix, keep)
}

func applyAndFinalize(module, binary, prefix string, keepHarnesses []string) error {
	root, rootErr := newProjectRoot()
	if rootErr != nil {
		return rootErr
	}

	if pruneErr := pruneHarnesses(keepHarnesses); pruneErr != nil {
		return pruneErr
	}

	files, collectErr := collectFiles(".")
	if collectErr != nil {
		return oops.Wrapf(collectErr, "collect files")
	}

	repls := []replacement{
		{label: "Updating module paths...", old: oldModule, new: module},
		{label: "Updating binary name...", old: oldBinary, new: binary},
	}

	if replErr := applyReplacements(root, files, repls); replErr != nil {
		return replErr
	}

	if renErr := renameCmdDir(binary); renErr != nil {
		return renErr
	}

	fixPathsAfterRename(files, binary)

	writeOut("Updating env prefix...")

	if replErr := applyReplacements(root, files, []replacement{
		{label: "", old: oldPrefix, new: prefix},
	}); replErr != nil {
		return replErr
	}

	return finalize(root, binary)
}

func applyReplacements(root projectRoot, files []string, repls []replacement) error {
	for _, repl := range repls {
		if repl.label != "" {
			writeOut(repl.label)
		}

		for _, fpath := range files {
			if err := root.replaceInFile(fpath, repl.old, repl.new); err != nil {
				return err
			}
		}
	}

	return nil
}

func fixPathsAfterRename(files []string, binary string) {
	oldDir := filepath.Join("cmd", oldBinary)
	newDir := filepath.Join("cmd", binary)

	for idx, fpath := range files {
		files[idx] = strings.Replace(fpath, oldDir, newDir, 1)
	}
}

func promptUser() (module, binary, prefix string, keepHarnesses []string, err error) {
	var confirm bool
	var selectedHarnesses []string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Module path").
				Placeholder("github.com/yourname/yourproject").
				Value(&module),
			huh.NewInput().
				Title("Binary name").
				Placeholder("yourproject").
				Value(&binary),
			huh.NewInput().
				Title("Environment variable prefix").
				Placeholder("YOURPROJECT").
				Value(&prefix),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select AI coding assistant harnesses to keep").
				Description("Unselected directories are removed; .agents/skills/ source is preserved.").
				Options(lo.Map(allHarnesses, func(h string, _ int) huh.Option[string] {
					return huh.NewOption(h, h)
				})...).
				Filterable(true).
				Height(10).
				Value(&selectedHarnesses),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Apply changes?").
				Affirmative("Yes").
				Negative("No").
				Value(&confirm),
		),
	)

	if formErr := form.Run(); formErr != nil {
		return "", "", "", nil, oops.Wrapf(formErr, "run form")
	}

	if !confirm {
		writeOut("Aborted.")

		return "", "", "", nil, nil
	}

	return module, binary, prefix, selectedHarnesses, nil
}

// pruneHarnesses removes harness directories not included in keep.
// Because each harness's skills/ subdir contains symlinks into .agents/skills/,
// os.RemoveAll only unlinks the symlinks and never follows them to the source.
func pruneHarnesses(keep []string) error {
	keepSet := lo.SliceToMap(keep, func(h string) (string, struct{}) {
		return h, struct{}{}
	})

	writeOut("Pruning unselected harnesses...")

	for _, harness := range allHarnesses {
		if _, ok := keepSet[harness]; ok {
			continue
		}

		if _, statErr := os.Stat(harness); os.IsNotExist(statErr) {
			continue
		} else if statErr != nil {
			return oops.Wrapf(statErr, "stat harness %s", harness)
		}

		if rmErr := os.RemoveAll(harness); rmErr != nil {
			return oops.Wrapf(rmErr, "remove harness %s", harness)
		}
	}

	return nil
}

func renameCmdDir(binary string) error {
	oldCmd := filepath.Join("cmd", oldBinary)
	newCmd := filepath.Join("cmd", binary)

	if oldBinary == binary {
		return nil
	}

	if _, statErr := os.Stat(oldCmd); os.IsNotExist(statErr) {
		return nil
	} else if statErr != nil {
		return oops.Wrapf(statErr, "stat cmd dir")
	}

	writeOut("Renaming cmd/" + oldBinary + " → cmd/" + binary + "...")

	return oops.Wrapf(os.Rename(oldCmd, newCmd), "rename cmd dir")
}

func finalize(root projectRoot, binary string) error {
	writeOut("Cleaning up init scaffolding...")

	if taskErr := removeTaskfileBlock(root, "Taskfile.yml", "init:"); taskErr != nil {
		return taskErr
	}

	if rmErr := os.RemoveAll("cmd/init"); rmErr != nil {
		return oops.Wrapf(rmErr, "remove cmd/init")
	}

	writeOut("Running go mod tidy...")

	cmd := exec.CommandContext(context.Background(), "go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if tidyErr := cmd.Run(); tidyErr != nil {
		return oops.Wrapf(tidyErr, "go mod tidy")
	}

	writeOut("")
	writeOut(doneStyle.Render("Done!"))
	writeOut("")
	writeOut("Next steps:")
	writeOut("  mise exec -- task ci")
	writeOut("  git add -A && git commit -m 'feat: initialize " + binary + "'")

	return nil
}

func removeTaskfileBlock(root projectRoot, fpath, blockName string) error {
	data, readErr := root.readFile(fpath)
	if readErr != nil {
		return readErr
	}

	lines := strings.Split(string(data), "\n")
	result := make([]string, 0, len(lines))
	skip := false

	for _, line := range lines {
		if strings.TrimSpace(line) == blockName && strings.HasPrefix(line, "  ") {
			skip = true

			continue
		}

		if skip {
			if isBlockEnd(line) {
				skip = false
			} else {
				continue
			}
		}

		result = append(result, line)
	}

	return root.writeFile(fpath, []byte(strings.Join(result, "\n")), 0o600)
}

func isBlockEnd(line string) bool {
	if line == "" || (line[0] != ' ' && line[0] != '\t') {
		return true
	}

	trimmed := strings.TrimSpace(line)

	return strings.HasPrefix(line, "  ") &&
		!strings.HasPrefix(line, "    ") &&
		strings.HasSuffix(trimmed, ":")
}

func collectFiles(root string) ([]string, error) {
	var files []string

	walkErr := filepath.Walk(root, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if skipDirs[filepath.Base(fpath)] {
				return filepath.SkipDir
			}

			return nil
		}

		if rewriteExts[filepath.Ext(fpath)] || strings.HasPrefix(filepath.Base(fpath), "Taskfile") {
			files = append(files, fpath)
		}

		return nil
	})

	return files, walkErr
}

func writeOut(msg string) {
	if _, err := os.Stdout.WriteString(msg + "\n"); err != nil {
		panic(err)
	}
}
