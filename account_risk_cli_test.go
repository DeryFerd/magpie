package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yetone/magpie/internal/provider"
)

// guiRisk reads the window's Add sheet risk list out of app.js: the agent of
// each entry marked `risk: true`, with its riskNote text when it carries one
// (an entry without one takes the dialog's default note).
func guiRisk(t *testing.T) map[string]string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("internal", "gui", "assets", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	all := string(b)
	out := map[string]string{}
	for i := 0; ; {
		j := strings.Index(all[i:], "risk: true")
		if j < 0 {
			break
		}
		j += i
		before := all[:j]
		a := strings.LastIndex(before, `agent: "`)
		if a < 0 {
			t.Fatalf("app.js: `risk: true` at byte %d has no agent before it; this scanner needs updating", j)
		}
		id := before[a+len(`agent: "`):]
		id = id[:strings.IndexByte(id, '"')]
		note := ""
		rest := all[j:]
		if n := strings.Index(rest, `riskNote: "`); n >= 0 {
			if na := strings.Index(rest, `agent: "`); na < 0 || n < na {
				s := rest[n+len(`riskNote: "`):]
				note = s[:strings.IndexByte(s, '"')]
			}
		}
		out[id] = note
		i = j + len("risk: true")
	}
	if len(out) == 0 {
		t.Fatal("app.js has no `risk: true` entry; this scanner needs updating with the page")
	}
	return out
}

// The CLI signs in the same subscriptions the window does, so it must ask
// about the risk of every one the window warns about: a warning on one
// surface only reaches half the users (accounts add, plugin login).
func TestAccountRiskCoversEveryGUIRisk(t *testing.T) {
	window := guiRisk(t)
	for id := range window {
		if provider.RiskOf(id) == "" {
			t.Errorf("the window asks before signing in to a %s account (app.js `risk: true`), the CLI does not: provider.RiskOf(%q) is empty", id, id)
		}
	}
	for id := range provider.RiskNotes() {
		if _, ok := window[id]; !ok {
			t.Errorf("the CLI warns about %q before sign-in, the window does not: app.js has no `risk: true` for it", id)
		}
	}
}

// The note the CLI prints is the one the window shows: one fact, one text.
func TestAccountRiskNotesAreTheGUITexts(t *testing.T) {
	notes := provider.RiskNotes()
	for id, note := range guiRisk(t) {
		if note == "" {
			continue // the entry takes the dialog's default note
		}
		if notes[id] != note {
			t.Errorf("%s's risk note has drifted apart\n window: %s\n cli:    %s", id, note, notes[id])
		}
	}
}

// A risky subscription stops the CLI until the user says y, as the window's
// Sign in anyway does.
func TestConfirmRiskStopsWithoutAYes(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = old }()
	go func() {
		_, _ = w.WriteString("n\n")
		_ = w.Close()
	}()
	if err := confirmRisk("claude"); err == nil {
		t.Fatal(`confirmRisk("claude") with "n" on stdin: expected the sign-in to be canceled`)
	}
}

// Anything the window doesn't warn about is signed in without a question.
func TestConfirmRiskAsksNothingForTheRest(t *testing.T) {
	if err := confirmRisk("codex"); err != nil {
		t.Fatalf(`confirmRisk("codex"): %v; expected nothing to ask`, err)
	}
}
