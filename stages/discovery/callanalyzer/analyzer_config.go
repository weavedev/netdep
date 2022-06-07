package callanalyzer

// DiscoveryAction indicates what to do when encountering
// a certain call. Used in interestingCalls
type DiscoveryAction int64

const (
	Output DiscoveryAction = iota
	Substitute
)

// defaultMaxTraversalDepth is the default max traversal depth for the analyser
const defaultMaxTraversalDepth = 64

// defaultMaxTraceDepth defines the maximum number of calls reported
// when verbosely printing unresolved call stack traces
const defaultMaxTraceDepth = 8

// InterestingCall holds information about a call that is to be outputted,
// substituted or otherwise inspected (by the analyzer).
type InterestingCall struct {
	action          DiscoveryAction
	interestingArgs []int
}

// AnalyserConfig holds the properties to adjust the analyser's behaviour for different use cases
type AnalyserConfig struct {
	// interestingCallsClient is a map from target to action that is to be taken when encountering the target.
	// Used internally to distinguish whether a call is to be:
	// outputted as a party in a dependency (0) or substituted with a constant (1)
	interestingCallsClient map[string]InterestingCall
	interestingCallsServer map[string]InterestingCall

	// substitutionCalls are the calls that are to be substituted with environment variable values
	substitutionCalls map[string]InterestingCall

	// environment: map[service name]map[variable name]value
	environment map[string]map[string]string

	// ignoreList is a set of function names to not recurse into
	ignoreList        map[string]bool
	verbose           bool
	maxTraversalDepth int
	maxTraceDepth     int
}

// SetVerbose is a setter for verbose
func (a *AnalyserConfig) SetVerbose(v bool) {
	a.verbose = v
}

// DefaultConfigForFindingHTTPCalls returns the default config
// for locating calls
func DefaultConfigForFindingHTTPCalls(environment map[string]map[string]string) AnalyserConfig {
	return AnalyserConfig{
		interestingCallsClient: map[string]InterestingCall{
			"(*net/http.Client).Do":          {action: Output, interestingArgs: []int{0}},
			"(*net/http.Client).Get":         {action: Output, interestingArgs: []int{1}},
			"(*net/http.Client).Post":        {action: Output, interestingArgs: []int{1}},
			"(*net/http.Client).Head":        {action: Output, interestingArgs: []int{1}},
			"net/http.Get":                   {action: Output, interestingArgs: []int{0, 1}},
			"net/http.Post":                  {action: Output, interestingArgs: []int{0, 1}},
			"net/http.NewRequestWithContext": {action: Output, interestingArgs: []int{2}},
			// "net/http.NewRequest":  ...
		},

		interestingCallsServer: map[string]InterestingCall{
			"net/http.Handle":                                 {action: Output, interestingArgs: []int{0}},
			"net/http.HandleFunc":                             {action: Output, interestingArgs: []int{0}},
			"net/http.ListenAndServe":                         {action: Output, interestingArgs: []int{0}},
			"(*github.com/gin-gonic/gin.RouterGroup).GET":     {action: Output, interestingArgs: []int{1}},
			"(*github.com/gin-gonic/gin.RouterGroup).PUT":     {action: Output, interestingArgs: []int{1}},
			"(*github.com/gin-gonic/gin.RouterGroup).POST":    {action: Output, interestingArgs: []int{1}},
			"(*github.com/gin-gonic/gin.RouterGroup).DELETE":  {action: Output, interestingArgs: []int{1}},
			"(*github.com/gin-gonic/gin.RouterGroup).PATCH":   {action: Output, interestingArgs: []int{1}},
			"(*github.com/gin-gonic/gin.RouterGroup).HEAD":    {action: Output, interestingArgs: []int{1}},
			"(*github.com/gin-gonic/gin.RouterGroup).OPTIONS": {action: Output, interestingArgs: []int{1}},
			"(*github.com/gin-gonic/gin.Engine).Run":          {action: Output, interestingArgs: []int{1}},
		},

		environment: environment,

		substitutionCalls: map[string]InterestingCall{
			"os.Getenv": {action: Substitute, interestingArgs: []int{0}}, // TODO: implement env var substitution
		},

		maxTraversalDepth: defaultMaxTraversalDepth,
		maxTraceDepth:     defaultMaxTraceDepth,
		verbose:           false,

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
			"sync":                 true,
			"syscall":              true,
			"unicode":              true,
			"time":                 true,
			// TODO: expand ignoreList
		},
	}
}
