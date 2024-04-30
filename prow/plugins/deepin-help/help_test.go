/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package deepinhelp

import (
	"fmt"
	"testing"

	"k8s.io/test-infra/prow/github"
)

type fakePruner struct{}

func (fp *fakePruner) PruneComments(shouldPrune func(github.IssueComment) bool) {}

func formatLabels(labels ...string) []string {
	r := []string{}
	for _, l := range labels {
		r = append(r, fmt.Sprintf("%s/%s#%d:%s", "org", "repo", 1, l))
	}
	if len(r) == 0 {
		return nil
	}
	return r
}

func TestHelpCommands(t *testing.T) {
	commandsSummary := "deepin pr test commands"
	type testCase struct {
		name        string
		expectedMsg string
	}
	testCases := []testCase{
		{
			name: "Deepin Help message with commands summary",
			expectedMsg: fmt.Sprintf(`
### deepin pr commands help details.

%s`, commandsSummary),
		},
	}

	for _, tc := range testCases {
		ig := prCommands{
			prCommnadsSummary: commandsSummary,
		}

		returnedMsg := ig.helpMsg()
		if returnedMsg != tc.expectedMsg {
			t.Errorf("(%s): Expected message: %sReturned message: %s", tc.name, tc.expectedMsg, returnedMsg)
		}
	}
}
