package types

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/Azure/aztfmigrate/helper"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// GetResourceBlock searches tf files in working directory and return `targetAddress` block
func GetResourceBlock(workingDirectory, targetAddress string) (*hclwrite.Block, error) {
	for _, file := range helper.ListHclFiles(workingDirectory) {
		// #nosec G304
		src, err := os.ReadFile(filepath.Join(workingDirectory, file.Name()))
		if err != nil {
			return nil, err
		}
		f, diag := hclwrite.ParseConfig(src, file.Name(), hcl.InitialPos)
		if f == nil || diag != nil && diag.HasErrors() || f.Body() == nil {
			continue
		}
		for _, block := range f.Body().Blocks() {
			if block != nil && block.Type() == "resource" {
				address := strings.Join(block.Labels(), ".")
				if targetAddress == address {
					return block, nil
				}
			}
		}
	}
	return nil, nil
}

// ReplaceResourceBlock searches tf files in working directory and replace `targetAddress` block with `newBlock`
func ReplaceResourceBlock(workingDirectory, targetAddress string, newBlocks []*hclwrite.Block) error {
	for _, file := range helper.ListHclFiles(workingDirectory) {
		// #nosec G304
		src, err := os.ReadFile(filepath.Join(workingDirectory, file.Name()))
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
					f.Body().AppendUnstructuredTokens(CommentOutBlock(block))
					f.Body().AppendNewline()
					for _, newBlock := range newBlocks {
						if newBlock == nil {
							continue
						}
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
			if err := os.WriteFile(filepath.Join(workingDirectory, file.Name()), hclwrite.Format(f.Bytes()), 0600); err != nil {
				log.Printf("[Error] saving configuration %s: %+v", file.Name(), err)
			}
			return nil
		}
	}
	return nil
}

// ReplaceGenericOutputs searches tf files in working directory and replace generic resource's output with new address
func ReplaceGenericOutputs(workingDirectory string, outputs []Output) error {
	for _, file := range helper.ListHclFiles(workingDirectory) {
		// #nosec G304
		src, err := os.ReadFile(filepath.Join(workingDirectory, file.Name()))
		if err != nil {
			return err
		}
		f, diag := hclwrite.ParseConfig(src, file.Name(), hcl.InitialPos)
		if f == nil || diag != nil && diag.HasErrors() || f.Body() == nil {
			continue
		}
		for _, block := range f.Body().Blocks() {
			if block.Type() == "removed" || block.Type() == "import" {
				continue
			}
			if block != nil {
				replaceOutputs(block, outputs)
			}
		}
		if err := os.WriteFile(filepath.Join(workingDirectory, file.Name()), hclwrite.Format(f.Bytes()), 0600); err != nil {
			log.Printf("[Error] saving configuration %s: %+v", file.Name(), err)
		}
	}
	return nil
}

func replaceOutputs(block *hclwrite.Block, outputs []Output) {
	for attrName, attr := range block.Body().Attributes() {
		attrValue := string(attr.Expr().BuildTokens(nil).Bytes())
		for _, output := range outputs {
			attrValue = strings.ReplaceAll(attrValue, output.OldName, output.NewName)
		}
		block.Body().SetAttributeRaw(attrName, helper.GetTokensForExpression(attrValue))
	}
	for index := range block.Body().Blocks() {
		replaceOutputs(block.Body().Blocks()[index], outputs)
	}
}

// UpdateMigratedResourceBlock searches tf files in working directory and update generic patch resource's target
func UpdateMigratedResourceBlock(workingDirectory string, resources []AzapiUpdateResource) error {
	for _, file := range helper.ListHclFiles(workingDirectory) {
		// #nosec G304
		src, err := os.ReadFile(filepath.Join(workingDirectory, file.Name()))
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
				for _, r := range resources {
					if r.NewAddress(nil) == address { // TODO: && r.Change.Action != no_op
						recursiveUpdate(block, r.Block, r.Change.Before, r.Change.After)
						break
					}
				}
			}
		}

		if err := os.WriteFile(filepath.Join(workingDirectory, file.Name()), hclwrite.Format(f.Bytes()), 0600); err != nil {
			log.Printf("[Error] saving configuration %s: %+v", file.Name(), err)
		}
	}
	return nil
}

func recursiveUpdate(old *hclwrite.Block, new *hclwrite.Block, before interface{}, after interface{}) {
	// user can't use patch resource to add item to some array, so we don't need to deal with before or after is an array
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
				old.Body().SetAttributeRaw(attrName, new.Body().GetAttribute(attrName).Expr().BuildTokens(nil))
				continue
			}
			// delete
			if beforeMap[attrName] == nil && afterMap[attrName] != nil {
				old.Body().RemoveAttribute(attrName)
				continue
			}

			// update
			if new.Body().GetAttribute(attrName) == nil {
				continue
			}
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

// InjectReference replaces `block`'s literal value with reference provided by `refs`
func InjectReference(block *hclwrite.Block, refs []Reference) *hclwrite.Block {
	search := make([]string, 0)
	replacement := make([]string, 0)
	for _, ref := range refs {
		if stringValue, ok := ref.Value.(string); ok {
			search = append(search, stringValue)
			replacement = append(replacement, ref.Name)
		}
	}
	for attrName, attr := range block.Body().Attributes() {
		if input := helper.GetValueFromExpression(attr.Expr().BuildTokens(nil)); input != nil {
			if output, found := helper.ToHclSearchReplace(input, search, replacement); found {
				block.Body().SetAttributeRaw(attrName, helper.GetTokensForExpression(output))
			}
		}
	}
	for index := range block.Body().Blocks() {
		InjectReference(block.Body().Blocks()[index], refs)
	}
	return block
}

// GetValuePropMap returns a map from literal value to reference
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

// CombineBlock combines `blocks` and update its result to `output` block, and return the difference in a map(key is attribute name, value is a list of attribute values)
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
		case helper.IsArrayWithSameValue(values):
			output.Body().SetAttributeRaw(attrName, blocks[0].Body().GetAttribute(attrName).Expr().BuildTokens(nil))
		case isForEach:
			output.Body().SetAttributeRaw(attrName, helper.GetTokensForExpression("each.value."+attrName))
			attrValueMap[attrName] = tokens
		default:
			output.Body().SetAttributeRaw(attrName, helper.GetTokensForExpression(fmt.Sprintf("%s${count.index}%s", helper.Prefix(values), helper.Suffix(values))))
			attrValueMap[attrName] = tokens
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

// GetForEachConstants converts a map of difference to hcl object
func GetForEachConstants(instances []Instance, items map[string][]hclwrite.Tokens) string {
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

func CommentOutBlock(block *hclwrite.Block) hclwrite.Tokens {
	file := hclwrite.NewEmptyFile()
	file.Body().AppendBlock(block)
	content := string(file.Bytes())
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lines[i] = fmt.Sprintf("# %s", line)
	}
	return hclwrite.Tokens{
		&hclwrite.Token{
			Type:         hclsyntax.TokenComment,
			Bytes:        []byte(strings.Join(lines, "\n")),
			SpacesBefore: 0,
		},
	}
}
