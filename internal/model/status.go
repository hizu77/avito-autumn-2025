//go:generate go-enum --output-suffix=.generated

package model

// Status is a PR status ("open" or "merged").
// ENUM(open, merged)
type Status string
