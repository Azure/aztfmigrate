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
				f, dialog := hclwrite.ParseConfig([]byte(fmt.Sprintf("%s=%s", attrName, output.NewName)), "", hcl.InitialPos)
				if dialog == nil || !dialog.HasErrors() && f != nil {
					block.Body().SetAttributeRaw(attrName, f.Body().GetAttribute(attrName).Expr().BuildTokens(nil))
				}
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
		for _, ref := range refs {
			if ref.Value == nil {
				continue
			}
			if ref.GetStringValue() == attrValue {
				f, dialog := hclwrite.ParseConfig([]byte(fmt.Sprintf("%s=%s", attrName, ref.Name)), "", hcl.InitialPos)
				if dialog == nil || !dialog.HasErrors() && f != nil {
					block.Body().SetAttributeRaw(attrName, f.Body().GetAttribute(attrName).Expr().BuildTokens(nil))
				}
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
