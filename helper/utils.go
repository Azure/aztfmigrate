package helper

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	ctyJson "github.com/zclconf/go-cty/cty/json"
)

// IsArrayWithSameValue returns true if all elements in `arr` are identical
func IsArrayWithSameValue(arr []string) bool {
	for _, x := range arr {
		if x != arr[0] {
			return false
		}
	}
	return true
}

// Prefix returns the longest common prefix of all elements in `arr`
func Prefix(arr []string) string {
	if len(arr) == 0 {
		return ""
	}
	index := 0
	for index = 0; index < len(arr[0]); index++ {
		match := true
		for i := 1; i < len(arr); i++ {
			if index >= len(arr[i]) || arr[i][index] != arr[0][index] {
				match = false
				break
			}
		}
		if !match {
			break
		}
	}
	return arr[0][0:index]
}

// Suffix returns the longest common suffix of all elements in `arr`
func Suffix(arr []string) string {
	if len(arr) == 0 {
		return ""
	}
	index := 0
	for index = 1; index <= len(arr[0]); index++ {
		match := true
		for i := 1; i < len(arr); i++ {
			if index > len(arr[i]) || arr[i][len(arr[i])-index] != arr[0][len(arr[0])-index] {
				match = false
				break
			}
		}
		if !match {
			break
		}
	}
	return arr[0][len(arr[0])-index+1:]
}

func ListHclFiles(workingDirectory string) []fs.DirEntry {
	res := make([]fs.DirEntry, 0)
	files, err := os.ReadDir(workingDirectory)
	if err != nil {
		return res
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".tf") {
			res = append(res, file)
		}
	}
	return res
}

// GetTokensForExpression convert a literal value to hclwrite.Tokens
func GetTokensForExpression(expression string) hclwrite.Tokens {
	syntaxTokens, diags := hclsyntax.LexConfig([]byte(expression), "main.tf", hcl.InitialPos)
	if diags.HasErrors() {
		return nil
	}
	res := make([]*hclwrite.Token, 0)
	for _, token := range syntaxTokens {
		res = append(res, &hclwrite.Token{
			Type:  token.Type,
			Bytes: token.Bytes,
		})
	}
	return res
}

// ParseHclArray parse `attrValue` to an array, example `attrValue` `["a", "b", 0]` will return ["\"a\"", "\"b\"", "0"]
func ParseHclArray(attrValue string) []string {
	if strings.HasPrefix(attrValue, "[") && strings.HasSuffix(attrValue, "]") {
		arr := strings.Split(attrValue[1:len(attrValue)-1], ",")
		for i := range arr {
			arr[i] = strings.TrimSpace(arr[i])
		}
		return arr
	}
	return nil
}

// ToHclSearchReplace generates hcl expression from `input`
func ToHclSearchReplace(input interface{}, search []string, replacement []string) (string, bool) {
	found := false
	switch value := input.(type) {
	case []interface{}:
		if len(value) == 0 {
			return "[]", false
		}
		res := make([]string, 0)
		for _, element := range value {
			config, ok := ToHclSearchReplace(element, search, replacement)
			found = found || ok
			res = append(res, config)
		}
		return fmt.Sprintf("[\n%s\n]", strings.Join(res, ",\n")), found
	case map[string]interface{}:
		if len(value) == 0 {
			return "{}", found
		}
		attrs := make([]string, 0)
		for k, v := range value {
			if v == nil {
				attrs = append(attrs, fmt.Sprintf("%s = null", quotedKey(k)))
				continue
			}
			config, ok := ToHclSearchReplace(v, search, replacement)
			found = found || ok
			attrs = append(attrs, fmt.Sprintf("%s = %s", quotedKey(k), config))
		}
		return fmt.Sprintf("{\n%s\n}", strings.Join(attrs, "\n")), found
	case string:
		for i := range search {
			if search[i] == value {
				return replacement[i], true
			}
		}
		return fmt.Sprintf(`"%s"`, strings.ReplaceAll(value, "\"", "\\\"")), false
	default:
		return fmt.Sprintf("%v", value), false
	}
}

func GetValueFromExpression(tokens hclwrite.Tokens) interface{} {
	expression, _ := hclsyntax.ParseExpression(tokens.Bytes(), "", hcl.InitialPos)
	if value, dialog := expression.Value(&hcl.EvalContext{}); dialog == nil || !dialog.HasErrors() {
		if data, err := ctyJson.Marshal(value, value.Type()); err == nil {
			var input interface{}
			if err = json.Unmarshal(data, &input); err == nil {
				return input
			}
		}
	}
	return nil
}

func quotedKey(input string) string {
	if len(input) == 0 {
		return input
	}
	if strings.Contains(input, ".") || strings.Contains(input, "/") || input[0] == '$' || input[0] >= '0' && input[0] <= '9' {
		return fmt.Sprintf("\"%s\"", input)
	}
	return input
}
