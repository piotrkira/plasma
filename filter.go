package main

import (
	"regexp"
	"strings"
)

type Threat string

var sqlInjectionPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\b(select|delete|update|alter)\b.*`),
}

var xssPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)<script.*?>`),
}

func ContainsSqlInjectionPattern(data string) bool {
	data = strings.ToLower(data)
	for _, pattern := range sqlInjectionPatterns {
		if pattern.MatchString(data) {
			return true
		}
	}
	return false
}

func ContainsXSSPatern(data string) bool {
	data = strings.ToLower(data)
	for _, pattern := range xssPatterns {
		if pattern.MatchString(data) {
			return true
		}
	}
	return false
}

func GetThreat(urlQuery, body string) Threat {
	if ContainsSqlInjectionPattern(urlQuery) {
		return "URL_SQL_INJECTION"
	}
	if ContainsSqlInjectionPattern(body) {
		return "BODY_SQL_INJECTION"
	}
	if ContainsXSSPatern(urlQuery) {
		return "URL_XSS_ATTACK"
	}
	if ContainsXSSPatern(body) {
		return "BODY_XSS_ATTACK"
	}
	return ""
}
