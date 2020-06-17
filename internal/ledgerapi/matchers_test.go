package ledger

import (
	"fmt"
	"github.com/onsi/gomega/types"
)

type strSliceMatcher struct {
	slice []string
}

func (strSlice *strSliceMatcher) Match(actual interface{}) (bool, error) {
	strV, ok := actual.(string)
	if !ok {
		return false, fmt.Errorf("failed to convert to string")
	}
	for _, str := range strSlice.slice {
		if strV == str {
			return true, nil
		}
	}
	return false, fmt.Errorf("expected %s to be a member of %v but it's not", actual, strSlice.slice)
}

func (strSlice *strSliceMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("expected %v to be a member of the string slice: %v", actual, strSlice.slice)
}

func (strSlice *strSliceMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("%s must be a member of the slice %v", actual, strSlice.slice)
}

func BeMemberOfStringSlice(strSlice []string) types.GomegaMatcher {
	return &strSliceMatcher{strSlice}
}
