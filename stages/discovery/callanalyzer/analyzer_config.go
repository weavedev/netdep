package callanalyzer

// DiscoveryAction indicates what to do when encountering
// a certain call. Used in interestingCalls
type DiscoveryAction int64

const (
	Output DiscoveryAction = iota
	Substitute
)

// defaultMaxTraversalDepth is the default max traversal depth for the analyser
const defaultMaxTraversalDepth = 16

type InterestingCall struct {
	action          DiscoveryAction
	interestingArgs []int
}

// AnalyserConfig holds the properties to adjust the analyser's behaviour for different use cases
// i.e. for client/server side code detection
type AnalyserConfig struct {
	// interestingCallsClient is a map from target to action that is to be taken when encountering the target.
	// Used internally to distinguish whether a call is to be:
	// outputted as a party in a dependency (0) or substituted with a constant (1)
	interestingCallsClient map[string]InterestingCall
	interestingCallsServer map[string]InterestingCall
	// ignoreList is a set of function names to not recurse into
	ignoreList        map[string]bool
	maxTraversalDepth int
}

// DefaultConfigForFindingHTTPClientCalls returns the default config
// for locating client calls
func DefaultConfigForFindingHTTPClientCalls() AnalyserConfig {
	return AnalyserConfig{
		interestingCallsClient: map[string]InterestingCall{
			"(*net/http.Client).Do":          {action: Output, interestingArgs: []int{0}},
			"(*net/http.Client).Get":         {action: Output, interestingArgs: []int{0, 1}},
			"(*net/http.Client).Post":        {action: Output, interestingArgs: []int{0, 1}},
			"(*net/http.Client).Head":        {action: Output, interestingArgs: []int{0, 1}},
			"net/http.Get":                   {action: Output, interestingArgs: []int{0, 1}},
			"net/http.Post":                  {action: Output, interestingArgs: []int{0, 1}},
			"net/http.NewRequestWithContext": {action: Output, interestingArgs: []int{2}},
			"os.Getenv":                      {action: Substitute, interestingArgs: []int{0}}, // TODO: implement env var substitution
			// "net/http.NewRequest":  ...
			// "(*net/http.Client).PostForm": ...
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
		maxTraversalDepth: defaultMaxTraversalDepth,
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
