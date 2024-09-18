package coverage

import (
	_ "embed"
	"encoding/json"
	"log"
	"strings"
)

//go:embed tf.json
var coverageJson string

var cov []Coverage

func init() {
	_ = json.Unmarshal([]byte(coverageJson), &cov)
	for i := range cov {
		cov[i].IdPattern = strings.ReplaceAll(cov[i].IdPattern, "/{}", "")
	}
	if len(cov) <= 10 {
		log.Printf("[WARN] Coverage report for DEVELOPMENT is loaded. Please use the released binaries in production.")
	}
}

func GetApiVersion(idPattern string) string {
	for _, r := range cov {
		if r.Operation != "PUT" {
			continue
		}
		if !strings.EqualFold(idPattern, r.IdPattern) {
			continue
		}
		return r.ApiVersion
	}
	return ""
}

func GetPutCoverage(props []string, idPattern string) ([]string, []string) {
	return getCoverage(props, "PUT", idPattern)
}

func GetGetCoverage(props []string, idPattern string) ([]string, []string) {
	return getCoverage(props, "GET", idPattern)
}

func getCoverage(props []string, operation, idPattern string) ([]string, []string) {
	for _, r := range cov {
		if r.Operation != operation {
			continue
		}
		if !strings.EqualFold(idPattern, r.IdPattern) {
			continue
		}
		propsSet := make(map[string]bool)
		propsSet["name"] = true
		for _, prop := range r.Properties {
			parts := strings.Split(prop.Name, "/")
			for i := range parts {
				if index := strings.Index(parts[i], "{"); index != -1 {
					parts[i] = parts[i][0:index]
				}
			}
			propsSet[strings.Join(parts, ".")] = true
		}
		covered := make([]string, 0)
		uncovered := make([]string, 0)
		for _, prop := range props {
			if propsSet[prop] {
				covered = append(covered, prop)
			} else {
				uncovered = append(uncovered, prop)
			}
		}
		return covered, uncovered
	}
	return []string{}, props
}
