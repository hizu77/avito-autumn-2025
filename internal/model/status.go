//go:generate go-enum --output-suffix=.generated

package model

// Status is a PR status ("open" or "merged").
// ENUM(Open=OPEN, Merged=MERGED)
type Status string
