package reconcile

import (
	"fmt"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// Action represents what should be done to reconcile a key.
type Action string

const (
	ActionAdd    Action = "ADD"
	ActionRemove Action = "REMOVE"
	ActionUpdate Action = "UPDATE"
	ActionKeep   Action = "KEEP"
)

// Step describes a single reconciliation step.
type Step struct {
	Key    string
	Action Action
	Value  string
	OldVal string
}

// Plan returns an ordered list of reconciliation steps needed to make
// source match target, based on the diff entries provided.
func Plan(entries []diff.Entry) []Step {
	steps := make([]Step, 0, len(entries))
	for _, e := range entries {
		switch e.Status {
		case diff.Added:
			steps = append(steps, Step{Key: e.Key, Action: ActionAdd, Value: e.TargetVal})
		case diff.Removed:
			steps = append(steps, Step{Key: e.Key, Action: ActionRemove, OldVal: e.SourceVal})
		case diff.Modified:
			steps = append(steps, Step{Key: e.Key, Action: ActionUpdate, Value: e.TargetVal, OldVal: e.SourceVal})
		case diff.Unchanged:
			steps = append(steps, Step{Key: e.Key, Action: ActionKeep, Value: e.SourceVal})
		}
	}
	sort.Slice(steps, func(i, j int) bool { return steps[i].Key < steps[j].Key })
	return steps
}

// Apply takes a source EnvFile and a list of steps and returns the reconciled
// map of key→value that represents the updated environment.
func Apply(source parser.EnvFile, steps []Step) parser.EnvFile {
	result := make(parser.EnvFile)
	for k, v := range source {
		result[k] = v
	}
	for _, s := range steps {
		switch s.Action {
		case ActionAdd, ActionUpdate:
			result[s.Key] = s.Value
		case ActionRemove:
			delete(result, s.Key)
		}
	}
	return result
}

// Render serialises an EnvFile back to .env file content.
func Render(env parser.EnvFile) string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, env[k]))
	}
	return sb.String()
}
