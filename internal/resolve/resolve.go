// Package resolve provides utilities for resolving field values
// from deployed service state against expected manifest values.
package resolve

import (
	"fmt"
	"strings"
)

// ServiceState represents the live state of a deployed service.
type ServiceState struct {
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Replicas    int               `json:"replicas"`
	Environment map[string]string `json:"environment,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// Field extracts the value of a named field from a ServiceState.
// Supported fields: name, image, replicas, env.<KEY>, label.<KEY>.
func Field(state ServiceState, field string) (string, error) {
	switch {
	case field == "name":
		return state.Name, nil
	case field == "image":
		return state.Image, nil
	case field == "replicas":
		return fmt.Sprintf("%d", state.Replicas), nil
	case strings.HasPrefix(field, "env."):
		key := strings.TrimPrefix(field, "env.")
		val, ok := state.Environment[key]
		if !ok {
			return "", fmt.Errorf("environment key %q not found in service %q", key, state.Name)
		}
		return val, nil
	case strings.HasPrefix(field, "label."):
		key := strings.TrimPrefix(field, "label.")
		val, ok := state.Labels[key]
		if !ok {
			return "", fmt.Errorf("label key %q not found in service %q", key, state.Name)
		}
		return val, nil
	default:
		return "", fmt.Errorf("unsupported field %q", field)
	}
}

// AllFields returns a map of all resolvable top-level fields for a ServiceState.
func AllFields(state ServiceState) map[string]string {
	fields := map[string]string{
		"name":     state.Name,
		"image":    state.Image,
		"replicas": fmt.Sprintf("%d", state.Replicas),
	}
	for k, v := range state.Environment {
		fields["env."+k] = v
	}
	for k, v := range state.Labels {
		fields["label."+k] = v
	}
	return fields
}
