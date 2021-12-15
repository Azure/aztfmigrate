package helper

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/ms-henglu/azurerm-restapi-to-azurerm/types"
)

func ReplaceResourceBlock(targetAddress string, newBlock *hclwrite.Block) error {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return err
	}
	files, err := ioutil.ReadDir(workingDirectory)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".tf") {
			continue
		}
		src, err := ioutil.ReadFile(filepath.Join(workingDirectory, file.Name()))
		if err != nil {
			return err
		}
		f, diag := hclwrite.ParseConfig(src, file.Name(), hcl.InitialPos)
		if f == nil || diag != nil && diag.HasErrors() || f.Body() == nil {
			continue
		}
		blocks := f.Body().Blocks()
		f.Body().Clear()
		found := false
		for _, block := range blocks {
			if block != nil && block.Type() == "resource" {
				address := strings.Join(block.Labels(), ".")
				if targetAddress == address {
					if newBlock != nil {
						f.Body().AppendBlock(newBlock)
						f.Body().AppendNewline()
					}
					found = true
					continue
				}
			}
			f.Body().AppendBlock(block)
			f.Body().AppendNewline()
		}
		if found {
			if err := os.WriteFile(file.Name(), hclwrite.Format(f.Bytes()), 0644); err != nil {
				log.Printf("[Error] saving configuration %s: %+v", file.Name(), err)
			}
			return nil
		}
	}
	return nil
}

func ReplaceGenericOutputs(outputs []types.Output) error {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return err
	}
	files, err := ioutil.ReadDir(workingDirectory)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".tf") {
			continue
		}
		src, err := ioutil.ReadFile(filepath.Join(workingDirectory, file.Name()))
		if err != nil {
			return err
		}
		f, diag := hclwrite.ParseConfig(src, file.Name(), hcl.InitialPos)
		if f == nil || diag != nil && diag.HasErrors() || f.Body() == nil {
			continue
		}
		for _, block := range f.Body().Blocks() {
			if block != nil {
				replaceOutputs(block, outputs)
			}
		}
		if err := os.WriteFile(file.Name(), hclwrite.Format(f.Bytes()), 0644); err != nil {
			log.Printf("[Error] saving configuration %s: %+v", file.Name(), err)
		}
	}
	return nil
}

func replaceOutputs(block *hclwrite.Block, outputs []types.Output) {
	for attrName, attr := range block.Body().Attributes() {
		attrValue := string(attr.Expr().BuildTokens(nil).Bytes())
		attrValue = strings.TrimSpace(attrValue)
		for _, output := range outputs {
			if attrValue == output.OldName {
				block.Body().SetAttributeRaw(attrName, GetTokensForExpression(output.NewName))
				break
			}
		}
	}
	for index := range block.Body().Blocks() {
		replaceOutputs(block.Body().Blocks()[index], outputs)
	}
}

func UpdateMigratedResourceBlock(resources []types.GenericPatchResource) error {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return err
	}
	files, err := ioutil.ReadDir(workingDirectory)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".tf") {
			continue
		}
		src, err := ioutil.ReadFile(filepath.Join(workingDirectory, file.Name()))
		if err != nil {
			return err
		}
		f, diag := hclwrite.ParseConfig(src, file.Name(), hcl.InitialPos)
		if f == nil || diag != nil && diag.HasErrors() || f.Body() == nil {
			continue
		}
		for _, block := range f.Body().Blocks() {
			if block != nil && block.Type() == "resource" {
				address := strings.Join(block.Labels(), ".")
				for index, r := range resources {
					if r.NewAddress() == address { // TODO: && r.Change.Action != no_op
						recursiveUpdate(block, r.Block, r.Change.Before, r.Change.After)
						resources[index].Migrated = true
						break
					}
				}
			}
		}

		if err := os.WriteFile(file.Name(), hclwrite.Format(f.Bytes()), 0644); err != nil {
			log.Printf("[Error] saving configuration %s: %+v", file.Name(), err)
		}
	}
	return nil
}

func recursiveUpdate(old *hclwrite.Block, new *hclwrite.Block, before interface{}, after interface{}) {
	beforeMap, ok1 := before.(map[string]interface{})
	afterMap, ok2 := after.(map[string]interface{})
	if !ok1 || !ok2 {
		return
	}
	attrs := make(map[string]bool)
	for attrName := range new.Body().Attributes() {
		attrs[attrName] = true
	}
	for attrName := range old.Body().Attributes() {
		attrs[attrName] = true
	}

	for attrName := range attrs {
		if !reflect.DeepEqual(beforeMap[attrName], afterMap[attrName]) {
			// add
			if beforeMap[attrName] != nil && afterMap[attrName] == nil {
				old.Body().SetAttributeRaw(attrName, new.Body().GetAttribute(attrName).BuildTokens(nil))
				continue
			}
			// delete
			if beforeMap[attrName] == nil && afterMap[attrName] != nil {
				old.Body().RemoveAttribute(attrName)
				continue
			}

			// update
			old.Body().SetAttributeRaw(attrName, new.Body().GetAttribute(attrName).Expr().BuildTokens(nil))
		}
	}

	blocks := make(map[string]bool)
	for _, block := range new.Body().Blocks() {
		blocks[block.Type()] = true
	}
	for _, block := range old.Body().Blocks() {
		blocks[block.Type()] = true
	}
	for blockName := range blocks {
		if !reflect.DeepEqual(beforeMap[blockName], afterMap[blockName]) {
			oldBlocks := make([]*hclwrite.Block, 0)
			for _, block := range old.Body().Blocks() {
				if block.Type() == blockName {
					oldBlocks = append(oldBlocks, block)
				}
			}
			newBlocks := make([]*hclwrite.Block, 0)
			for _, block := range new.Body().Blocks() {
				if block.Type() == blockName {
					newBlocks = append(newBlocks, block)
				}
			}

			// add
			if len(newBlocks) != 0 && len(oldBlocks) == 0 {
				for _, block := range newBlocks {
					old.Body().AppendBlock(block)
				}
				continue
			}
			// delete
			if len(newBlocks) == 0 && len(oldBlocks) != 0 {
				for _, block := range oldBlocks {
					old.Body().RemoveBlock(block)
				}
				continue
			}

			// update
			if len(newBlocks) == len(oldBlocks) {
				beforeArr := make([]interface{}, 0)
				afterArr := make([]interface{}, 0)
				if beforeMap[blockName] != nil {
					temp, _ := beforeMap[blockName].([]interface{})
					beforeArr = temp
				}
				if afterMap[blockName] != nil {
					temp, _ := afterMap[blockName].([]interface{})
					afterArr = temp
				}
				if len(beforeArr) != len(afterArr) && len(beforeArr) != len(oldBlocks) {
					log.Fatal()
				}
				for index := range newBlocks {
					recursiveUpdate(oldBlocks[index], newBlocks[index], beforeArr[index], afterArr[index])
				}
			} else {
				for _, block := range oldBlocks {
					old.Body().RemoveBlock(block)
				}
				for _, block := range newBlocks {
					old.Body().AppendBlock(block)
				}
			}
		}
	}
}

func InjectReference(block *hclwrite.Block, refs []types.Reference) *hclwrite.Block {
	for attrName, attr := range block.Body().Attributes() {
		attrValue := string(attr.Expr().BuildTokens(nil).Bytes())
		attrValue = strings.TrimSpace(attrValue)
		if strings.HasPrefix(attrValue, "[") && strings.HasSuffix(attrValue, "]") {
			arr := strings.Split(attrValue[1:len(attrValue)-1], ",")
			found := false
			for i, v := range arr {
				for _, ref := range refs {
					if strings.TrimSpace(v) == ref.GetStringValue() {
						arr[i] = ref.Name
						found = true
						break
					}
				}
			}
			if found {
				newValue := fmt.Sprintf("[%s]", strings.Join(arr, ", "))
				block.Body().SetAttributeRaw(attrName, GetTokensForExpression(newValue))
				continue
			}
		}
		for _, ref := range refs {
			if ref.Value == nil {
				continue
			}
			if ref.GetStringValue() == attrValue {
				block.Body().SetAttributeRaw(attrName, GetTokensForExpression(ref.Name))
				break
			}
		}
	}
	for index := range block.Body().Blocks() {
		InjectReference(block.Body().Blocks()[index], refs)
	}
	return block
}

func GetValuePropMap(block *hclwrite.Block, prefix string) map[string]string {
	res := make(map[string]string)
	if block == nil {
		return res
	}
	for attrName, attr := range block.Body().Attributes() {
		res[strings.TrimSpace(string(attr.Expr().BuildTokens(nil).Bytes()))] = prefix + "." + attrName
	}
	blocksMap := make(map[string][]*hclwrite.Block)
	for _, block := range block.Body().Blocks() {
		if len(blocksMap[block.Type()]) == 0 {
			blocksMap[block.Type()] = make([]*hclwrite.Block, 0)
		}
		blocksMap[block.Type()] = append(blocksMap[block.Type()], block)
	}
	for blockType, arr := range blocksMap {
		for index, block := range arr {
			propValueMap := GetValuePropMap(block, fmt.Sprintf("%s.%s.%d", prefix, blockType, index))
			for k, v := range propValueMap {
				res[k] = v
			}
		}
	}
	return res
}

func CombineBlock(blocks []*hclwrite.Block, output *hclwrite.Block, isForEach bool) map[string][]hclwrite.Tokens {
	attrNameSet := make(map[string]bool)
	for _, b := range blocks {
		for attrName := range b.Body().Attributes() {
			attrNameSet[attrName] = true
		}
	}
	attrValueMap := make(map[string][]hclwrite.Tokens)
	for attrName := range attrNameSet {
		values := make([]string, len(blocks))
		tokens := make([]hclwrite.Tokens, len(blocks))
		for i, b := range blocks {
			if b == nil {
				values[i] = "null"
				tokens[i] = nil
				continue
			}
			attr := b.Body().GetAttribute(attrName)
			if attr == nil {
				values[i] = "null"
				tokens[i] = nil
			} else {
				tokens[i] = b.Body().GetAttribute(attrName).Expr().BuildTokens(nil)
				values[i] = strings.TrimSpace(string(tokens[i].Bytes()))
			}
		}
		switch {
		case isArrayWithSameValue(values):
			output.Body().SetAttributeRaw(attrName, blocks[0].Body().GetAttribute(attrName).Expr().BuildTokens(nil))
			break
		case isForEach:
			output.Body().SetAttributeRaw(attrName, GetTokensForExpression("each.value."+attrName))
			attrValueMap[attrName] = tokens
			break
		default:
			output.Body().SetAttributeRaw(attrName, GetTokensForExpression(fmt.Sprintf("%s${count.index}%s", prefix(values), suffix(values))))
			attrValueMap[attrName] = tokens
			break
		}
	}

	blockNameSet := make(map[string]bool)
	for _, b := range blocks {
		for _, nb := range b.Body().Blocks() {
			blockNameSet[nb.Type()] = true
		}
	}
	for blockName := range blockNameSet {
		nestedBlocks := make([]*hclwrite.Block, len(blocks))
		for i, b := range blocks {
			if nestedBlock := b.Body().FirstMatchingBlock(blockName, []string{}); nestedBlock != nil {
				nestedBlocks[i] = nestedBlock
			}
		}
		outputNestedBlock := output.Body().AppendNewBlock(blockName, []string{})
		tempMap := CombineBlock(nestedBlocks, outputNestedBlock, isForEach)
		for k, v := range tempMap {
			attrValueMap[k] = v
		}
	}
	return attrValueMap
}

func GetTokensForExpression(reference string) hclwrite.Tokens {
	f, dialog := hclwrite.ParseConfig([]byte(fmt.Sprintf("%s=%s", "temp", reference)), "", hcl.InitialPos)
	if dialog == nil || !dialog.HasErrors() && f != nil {
		return f.Body().GetAttribute("temp").Expr().BuildTokens(nil)
	}
	return nil
}

func GetForEachConstants(instances []types.Instance, items map[string][]hclwrite.Tokens) string {
	config := ""
	i := 0
	for _, instance := range instances {
		item := ""
		for key := range items {
			item += fmt.Sprintf("%s = %s", key, string(items[key][i].Bytes()))
		}
		config += fmt.Sprintf("%s = {\n%s\n}\n", instance.Index, item)
		i++
	}
	config = fmt.Sprintf("{\n%s}\n", config)
	return config
}

func isArrayWithSameValue(arr []string) bool {
	for _, x := range arr {
		if x != arr[0] {
			return false
		}
	}
	return true
}

func prefix(arr []string) string {
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

func suffix(arr []string) string {
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
