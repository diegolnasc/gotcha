package cmd

import "testing"

func TestVersion(t *testing.T) {
	cmd := &GotchaVersion{Version: "test-1.0.0"}

	if err := cmd.Init().Execute(); err != nil {
		t.Errorf("unexpected return %s", err.Error())
		t.FailNow()
	}
}
