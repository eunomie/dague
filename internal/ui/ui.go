package ui

import (
	"sort"

	"github.com/AlecAivazis/survey/v2"
)

func Select(msg string, options []string) (string, error) {
	opts := append([]string{}, options...)
	sort.Strings(opts)
	qs := &survey.Select{
		Message: msg,
		Options: opts,
	}
	var selected string
	err := survey.AskOne(qs, &selected, nil)
	return selected, err
}
