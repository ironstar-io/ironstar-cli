package cache

import (
	"strings"
	"testing"
)

func TestPublicCacheCommandHelp(t *testing.T) {
	if CacheCmd.Hidden {
		t.Fatal("cache command must be visible")
	}

	for _, value := range []string{
		"one HTTPS URL",
		"whole-environment invalidations are always",
		"--yes (-y)",
		"--environment dev",
		"--environment stage",
	} {
		if !strings.Contains(CacheCmd.Long+CacheCmd.Example, value) {
			t.Errorf("cache help does not contain %q", value)
		}
	}

	for _, forbidden := range []string{"--environment prod", "--environment production"} {
		if strings.Contains(CacheCmd.Example, forbidden) {
			t.Errorf("cache examples must not contain %q", forbidden)
		}
	}
}

func TestInvalidateCommandHelp(t *testing.T) {
	if InvalidateCmd.Hidden {
		t.Fatal("invalidate command must be visible")
	}

	help := InvalidateCmd.Long + InvalidateCmd.Example
	for _, value := range []string{
		"Supply one --url",
		"soft by default",
		"--type hard",
		"FAILED status",
		"not found in cache",
	} {
		if !strings.Contains(help, value) {
			t.Errorf("invalidate help does not contain %q", value)
		}
	}
}

func TestLegacyInvalidationCommandsAreHidden(t *testing.T) {
	for name, command := range map[string]bool{
		"invalidation": InvalidationCmd.Hidden,
		"create":       CreateCmd.Hidden,
		"list":         LegacyListInvalidationsCmd.Hidden,
		"show":         LegacyShowInvalidationCmd.Hidden,
	} {
		if !command {
			t.Errorf("legacy %s command must be hidden", name)
		}
	}
}
