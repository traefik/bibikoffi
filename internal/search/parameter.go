package search

import (
	"fmt"
	"time"
)

const gitHubSearchDateLayout = "2006-01-02"

// Parameter search parameter
type Parameter func() string

// NoOp no operation
func NoOp() string {
	return ""
}

// Cond apply conditionally some parameters
func Cond(apply bool, parameters ...Parameter) Parameter {
	if apply {
		return func() string {
			var filter string
			for _, param := range parameters {
				if param != nil {
					filter += param()
				}
			}
			return filter
		}
	}
	return NoOp
}

// WithLabels add a search filter by labels
func WithLabels(labels ...string) Parameter {
	var labelsFilter string
	for _, lbl := range labels {
		labelsFilter += fmt.Sprintf("label:%s ", lbl)
	}
	return func() string {
		return " " + labelsFilter
	}
}

// WithExcludedLabels add a search filter by unwanted labels
func WithExcludedLabels(labels ...string) Parameter {
	var labelsFilter string
	for _, lbl := range labels {
		labelsFilter += fmt.Sprintf("-label:%s ", lbl)
	}
	return func() string {
		return " " + labelsFilter
	}
}

// CreatedBefore issues created before a number of days
func CreatedBefore(days int) Parameter {
	if days == 0 {
		return NoOp
	}

	return func() string {
		return fmt.Sprintf(" created:<%s ", daysToDate(days))
	}
}

// CreatedAfter issues created after a number of days
func CreatedAfter(days int) Parameter {
	if days == 0 {
		return NoOp
	}

	return func() string {
		return fmt.Sprintf(" created:>=%s ", daysToDate(days))
	}
}

// UpdatedBefore issues updated before a number of days
func UpdatedBefore(days int) Parameter {
	if days == 0 {
		return NoOp
	}

	return func() string {
		return fmt.Sprintf(" updated:<%s ", daysToDate(days))
	}
}

// UpdatedAfter issues updated after a number of days
func UpdatedAfter(days int) Parameter {
	if days == 0 {
		return NoOp
	}

	return func() string {
		return fmt.Sprintf(" updated:>=%s ", daysToDate(days))
	}
}

// State issues state
func State(v string) Parameter {
	if v == "" {
		return NoOp
	}

	return func() string {
		return fmt.Sprintf(" state:%s ", v)
	}
}

// Lock conversation state.
func Lock(v bool) Parameter {
	if v {
		return func() string {
			return " is:locked "
		}
	}

	return func() string {
		return " is:unlocked "
	}
}

func daysToDate(days int) string {
	date := time.Now().Add(-time.Duration(days) * 24 * time.Hour)
	return date.Format(gitHubSearchDateLayout)
}
