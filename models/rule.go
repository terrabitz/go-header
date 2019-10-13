package models

import (
	"io/ioutil"
	"regexp"

	"github.com/denis-tingajkin/go-header/messages"
)

//Rule means rule for matching files
type Rule struct {
	//Template means license header for files
	Template string `yaml:"template"`
	//TemplatePath means license header for files located to specific folder
	TemplatePath string `yaml:"template-path"`
	//PathMatcher means regex for file path
	PathMatcher string `yaml:"path-matcher"`
	//AuthorMatcher means author regex for authors
	AuthorMatcher string `yaml:"author-matcher"`
	//ExcludePathMatcher means regex pattern to exclude files
	ExcludePathMatcher string `yaml:"exclude-path-matcher"`
	authorMatcher      *regexp.Regexp
	pathMatcher        *regexp.Regexp
	excludePathMatcher *regexp.Regexp
}

func (r *Rule) loadTemplate() error {
	if r.Template == "" && r.TemplatePath != "" {
		bytes, err := ioutil.ReadFile(r.TemplatePath)
		if err != nil {
			return messages.CanNotLoadTemplateFromFile(err)
		}
		r.Template = string(bytes)
	}
	if r.Template == "" {
		return messages.TemplateNotProvided()
	}
	return nil
}

func (r *Rule) Compile() messages.ErrorList {
	result := messages.NewErrorList()
	var err error
	if r.PathMatcher != "" {
		if r.pathMatcher, err = regexp.Compile(r.PathMatcher); err != nil {
			result.Append(err)
		}
	}
	if r.AuthorMatcher != "" {
		if r.authorMatcher, err = regexp.Compile(r.AuthorMatcher); err != nil {
			result.Append(err)
		}
	}
	if r.ExcludePathMatcher != "" {
		if r.excludePathMatcher, err = regexp.Compile(r.ExcludePathMatcher); err != nil {
			result.Append(err)
		}
	}
	return result
}

func (r Rule) Match(s *Source) bool {
	if r.pathMatcher != nil {
		if !r.pathMatcher.MatchString(s.Path) {
			return false
		}
	}
	if r.excludePathMatcher != nil {
		if r.excludePathMatcher.MatchString(s.Path) {
			return false
		}
	}
	if r.authorMatcher != nil {
		if !r.authorMatcher.MatchString(s.Author) {
			return false
		}
	}
	return true
}
