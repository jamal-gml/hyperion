package identity

import (
	"os"
	"strings"
)

type Identity struct{ ID, Project, Role, Provider string }

func FromEnv() Identity {
	id := os.Getenv("HYPERION_AGENT_ID")
	if id == "" {
		id = "local-assistant"
	}
	project := os.Getenv("HYPERION_PROJECT")
	if project == "" {
		project = "default"
	}
	role := os.Getenv("HYPERION_ROLE")
	if role == "" {
		role = "assistant"
	}
	provider := os.Getenv("HYPERION_PROVIDER")
	if provider == "" {
		provider = "unknown"
	}
	return Identity{ID: clean(id), Project: clean(project), Role: clean(role), Provider: clean(provider)}
}
func clean(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 100 {
		return s[:100]
	}
	return s
}
