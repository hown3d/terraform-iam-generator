package aws

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hown3d/terraform-iam-generator/internal/metrics"
)

func GenerateIamPolicy(msgs []metrics.CsmMessage) ([]byte, error) {
	policy := policyDocument{
		Version: "2012-10-17",
	}
	actions := make(map[string]struct{})
	for _, msg := range msgs {
		action := fmt.Sprintf("%v:%v", strings.ToLower(msg.Service), msg.API)
		actions[action] = struct{}{}
	}

	actionsFromMap := func() []string {
		keys := make([]string, len(actions))
		i := 0
		for k := range actions {
			keys[i] = k
			i++
		}
		return keys
	}

	policy.Statement = []statementEntry{{
		Effect:   "Allow",
		Action:   actionsFromMap(),
		Resource: "*",
	}}
	data, err := json.MarshalIndent(policy, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("marshaling policy: %w", err)
	}
	return data, nil
}

type policyDocument struct {
	Version   string
	Statement []statementEntry
}

type statementEntry struct {
	Effect   string
	Action   []string
	Resource string
}
