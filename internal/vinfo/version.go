// Package vinfo exposes build metadata for CLI version reporting.
package vinfo

import (
	"fmt"
	"runtime/debug"
	"strings"
)

const (
	developmentVersion = "dev"
	unknownCommit      = "none"
	unknownBuildDate   = "unknown"
)

// String returns a human-readable build version string.
func String() string {
	metadata := readBuildMetadata()

	return fmt.Sprintf(
		"%s (commit=%s, built=%s)",
		metadata.version,
		metadata.commit,
		metadata.buildDate,
	)
}

func readBuildMetadata() buildMetadata {
	metadata := buildMetadata{
		version:   developmentVersion,
		commit:    unknownCommit,
		buildDate: unknownBuildDate,
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return metadata
	}

	metadata.version = fallbackVersion(info.Main.Version)

	if info.Main.Version == "" || info.Main.Version == "(devel)" {
		metadata.version = revisionVersion(info.Settings)
	}

	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			metadata.commit = setting.Value
		case "vcs.time":
			metadata.buildDate = setting.Value
		}
	}

	return metadata
}

func revisionVersion(settings []debug.BuildSetting) string {
	for _, setting := range settings {
		if setting.Key == "vcs.revision" && setting.Value != "" {
			return setting.Value
		}
	}

	return developmentVersion
}

func fallbackVersion(version string) string {
	if version == "" || version == "(devel)" {
		return developmentVersion
	}

	return strings.TrimSpace(version)
}
