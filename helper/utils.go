package helper

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
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

func ListHclFiles() []fs.FileInfo {
	res := make([]fs.FileInfo, 0)
	workingDirectory, err := os.Getwd()
	if err != nil {
		return res
	}
	files, err := ioutil.ReadDir(workingDirectory)
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
	f, dialog := hclwrite.ParseConfig([]byte(fmt.Sprintf("%s=%s", "temp", expression)), "", hcl.InitialPos)
	if dialog == nil || !dialog.HasErrors() && f != nil {
		return f.Body().GetAttribute("temp").Expr().BuildTokens(nil)
	}
	return nil
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
