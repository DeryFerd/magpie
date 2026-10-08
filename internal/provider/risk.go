package provider

// The warnings for the subscriptions whose vendors may treat magpie's use
// as a breach. The texts are the ones the window's Add sheet shows (app.js
// riskNote); Claude's and Antigravity's are named in google.go beside the
// rest of their provider code.
const (
	CommandCodeRisk = "A Go plan account is used through Command Code's private interface, which Command Code may treat as a breach of its terms and ban the account for. Pro, Max and the other plans use its Provider API. Use a Go account you can afford to lose."
	QoderRisk       = "Qoder has no public API for this; magpie signs requests as its desktop client would, which Qoder may treat as third-party use and act on. Use an account you can afford to lose."
	ZedRisk         = "Zed serves these models to its own editor; magpie signs requests as the editor would, which Zed may treat as third-party use and act on. Use an account you can afford to lose."
	FactoryRisk     = "Factory serves these models to its own Droid CLI; magpie signs requests as Droid would, which Factory may treat as third-party use and act on. Use an account you can afford to lose."
	MimoRisk        = "Xiaomi serves these models to its own MiMo app; magpie signs requests as the app would, which Xiaomi may treat as third-party use and act on. Use an account you can afford to lose."
)

// riskNotes is the source of truth for the CLI: subscription id → what
// magpie says before signing in to an account of it.
var riskNotes = map[string]string{
	"claude":           ClaudeRisk,
	"antigravity":      AntigravityRisk,
	"commandcode-plan": CommandCodeRisk,
	"qoder":            QoderRisk,
	"qoder-cn":         QoderRisk,
	"zed":              ZedRisk,
	"factory":          FactoryRisk,
	"mimo-app":         MimoRisk,
}

// RiskNotes is, per subscription id, what magpie says before signing in to
// an account of it: the vendor may treat magpie's use as a breach and act on
// the account.
//
// The window's Add sheet keeps the same list (app.js, `risk: true` with its
// `riskNote`), and TestAccountRiskCoversEveryGUIRisk in the root package
// fails when the two drift apart, so a new risky subscription reaches the
// CLI as well as the window.
func RiskNotes() map[string]string {
	out := make(map[string]string, len(riskNotes))
	for id, note := range riskNotes {
		out[id] = note
	}
	return out
}

// RiskOf is the warning for one subscription's accounts, "" for the rest.
func RiskOf(id string) string { return riskNotes[id] }
