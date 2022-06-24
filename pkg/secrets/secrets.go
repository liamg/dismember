package secrets

import (
	"fmt"
	"regexp"
	"strings"

	squealer "github.com/owenrumney/squealer/pkg/config"
)

type Pattern struct {
	Regex  *regexp.Regexp
	Name   string
	Source string
}

func (p *Pattern) String() string {
	if p.Name == "" {
		return p.Regex.String()
	}
	return fmt.Sprintf("%s (pattern from %s)", p.Name, p.Source)
}

func Patterns() []Pattern {
	var all []Pattern
	for _, rule := range squealer.DefaultConfig().Rules {
		if rule.FileFilter != "" {
			continue
		}
		r, err := regexp.Compile(rule.Rule)
		if err != nil {
			continue
		}
		all = append(all, Pattern{
			Regex:  r,
			Name:   strings.ReplaceAll(rule.Description, "Check for ", ""),
			Source: "https://github.com/owenrumney/squealer",
		})
	}
	return all
}
