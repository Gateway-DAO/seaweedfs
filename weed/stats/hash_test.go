package stats

import (
	"regexp"
	"testing"

	"gotest.tools/assert"
)

func Test_regexp(t *testing.T) {
	paths := []string{
		"/data/gateway-private_1.idx",
		"/data/gateway-private_1.dat",
		"/data/gateway-private_1.dat/edge_case_path",
	}
	expected_matches := 1
	curr_matches := 0

	r, err := regexp.CompilePOSIX(`\.dat$`)
	if err != nil {
		t.Error(err)
	}

	for _, path := range paths {
		if r.MatchString(path) {
			curr_matches++
		}
	}

	assert.Equal(t, expected_matches, curr_matches, "expected regex matches")
}
