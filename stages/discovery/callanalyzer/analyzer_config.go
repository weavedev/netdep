package callanalyzer

// DiscoveryAction indicates what to do when encountering
// a certain call. Used in interestingCalls
type DiscoveryAction int64

const (
	Output DiscoveryAction = iota
	Substitute
)

const depth = 16

type AnalyzerConfig struct {
	// interestingCalls is a map from target to action that is to be taken when encountering the target.
	// Used internally to distinguish whether a call is to be:
	// outputted as a party in a dependency (0) or substituted with a constant (1)
	interestingCalls map[string]DiscoveryAction
	// ignoreList is a set of function names to not recurse into
	ignoreList  map[string]bool
	maxRecDepth int
}

func DefaultConfig() AnalyzerConfig {
	return AnalyzerConfig{
		interestingCalls: map[string]DiscoveryAction{
			"(*net/http.Client).Do":   Output,
			"os.Getenv":               Substitute,
			"(*net/http.Client).Get":  Output,
			"(*net/http.Client).Post": Output,
			"(*net/http.Client).Head": Output,
			"NewRequestWithContext":   Output,
			// "net/http.NewRequest":   Output,
			// "(*net/http.Client).PostForm":      Output,
		},
		maxRecDepth: depth,
		ignoreList: map[string]bool{
			"fmt":                  true,
			"reflect":              true,
			"net/url":              true,
			"strings":              true,
			"bytes":                true,
			"io":                   true,
			"errors":               true,
			"runtime":              true,
			"math/bits":            true,
			"internal/reflectlite": true,
		},
	}
}
