package updater_test

import (
	"testing"

	"github.com/gladiusio/gladius-guardian/updater"
)

func TestGetOfficialVersions(t *testing.T) {
	version, err := updater.GetOfficialVersions()
	if err != nil {
		t.Error(err)
	}
	if version["gladius-guardian"] != "0.8.0" {
		t.Error("Expected version 0.8.0")
	}
}

func TestGetVersion(t *testing.T) {
	version, err := updater.GetVersion("guardian")
	if err != nil {
		t.Error(err)
	}
	if version != "0.8.0" {
		t.Error("Expected version 0.8.0")
	}
}

func TestCompareVersion(t *testing.T) {
	officialVersion, err := updater.GetOfficialVersions()
	if err != nil {
		t.Error(err)
	}
	version, err := updater.CompareVersion("0.8.0", officialVersion["gladius-guardian"])
	if err != nil {
		t.Error(err)
	}

	switch version {
	case -1:
		t.Error("You are on an older version")
	case 1:
		t.Error("You are on a newer version")
	}
}
