package services

import (
	"os"
	"path/filepath"

	"github.com/ironstar-io/ironstar-cli/internal/system/fs"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func InitializeIronstarProject() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	confPath := filepath.Join(wd, ".ironstar", "config.yml")

	exists := fs.CheckExists(confPath)
	if exists {
		return errors.New("A config file has already been initialized in this directory")
	}

	// Package excludes now live in .ironstar/.ironignore (gitignore syntax), not
	// in config.yml's package.exclude.
	projConf := types.ProjectConfig{
		Version: "1",
	}

	newMarhsalled, err := yaml.Marshal(projConf)
	if err != nil {
		return err
	}

	err = SafeTouchConfigYAML(confPath)
	if err != nil {
		return err
	}

	fs.Replace(confPath, newMarhsalled, 0400)

	ignorePath := filepath.Join(wd, ".ironstar", ".ironignore")
	if !fs.CheckExists(ignorePath) {
		if err := fs.TouchByteArray(ignorePath, []byte(defaultIronignore), 0644); err != nil {
			return err
		}
	}

	return nil
}

// ironignoreHeader documents the file format and precedence. It heads every
// generated .ironignore.
const ironignoreHeader = `# .ironignore — files excluded from the package uploaded to Ironstar.
#
# Uses .gitignore syntax (https://git-scm.com/docs/gitignore):
#   - a pattern with no slash matches at ANY depth (e.g. node_modules, *.sql)
#   - a pattern containing a slash is anchored to this project root
#   - a leading "!" re-includes; a trailing "/" matches directories only
#
# This file takes precedence over any "package.exclude" in config.yml.
# Tune it for your stack — these are sensible defaults, not a complete list.
# Note: vendor/ (Composer) and built theme assets are intentionally NOT excluded.
`

// defaultIronignoreRules is the grouped starter exclude set. It targets bloat
// and never-at-runtime content while deliberately keeping things a deploy needs
// (e.g. vendor/, built theme assets).
const defaultIronignoreRules = `# Version control
.git
.svn
.hg

# Build dependencies (rebuilt on deploy; the source tree is never run)
node_modules
bower_components

# Local development environments
.ddev
.lando.yml
.lando.local.yml
docker-compose.override.yml

# Editor / OS cruft
.DS_Store
Thumbs.db
*.swp
.idea/
.vscode/

# Logs, caches, temp
*.log
.cache
tmp/

# Database dumps & backups
*.sql
*.sql.gz
*.dump

# Secrets & local environment
.env
.env.*
settings.local.php
settings.ddev.php

# CI & project/agent tooling
.github/
.gitlab-ci.yml
.opencode
.claude

# Drupal: user-uploaded files & generated/local settings
/private
sites/*/files
web/sites/*/files
docroot/sites/*/files
`

// defaultIronignore is the full file written by `iron init`.
const defaultIronignore = ironignoreHeader + "\n" + defaultIronignoreRules
