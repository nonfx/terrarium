// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package generate

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/cldcvr/terrarium/src/cli/internal/constants"
	"github.com/cldcvr/terrarium/src/pkg/metadata/app"
	"github.com/cldcvr/terrarium/src/pkg/metadata/platform"
	tfwriter "github.com/cldcvr/terrarium/src/pkg/tf/writer"
	"github.com/rotisserie/eris"
)

func blocksToPull(g platform.Graph, components ...string) []platform.BlockID {
	blockIDs := []platform.BlockID{}
	for _, comp := range components {
		bID := platform.NewBlockID(platform.BlockType_ModuleCall, platform.ComponentPrefix+comp)
		blockIDs = append(blockIDs, bID)
	}

	return blockIDs
}
func writeTF(g platform.Graph, destDir string, apps app.Apps, tfModule *tfconfig.Module, profileName string) (blockCount int, err error) {
	appDeps := apps.GetUniqueDependencyTypes()
	blocks := blocksToPull(g, appDeps...)

	log.Info("found dependencies", "dependencies", appDeps)

	locals, fileIndex, count, err := processBlocks(g, blocks, tfModule)
	if err != nil {
		return count, err
	}

	err = copyFilesFromIndex(tfModule.Path, destDir, fileIndex)
	if err != nil {
		return count, err
	}

	if len(locals) > 0 {
		err = writeLocalsToFile(locals, destDir, apps)
		if err != nil {
			return count, err
		}
	}

	if profileName != "" {
		if err := copyProfileConfigurationFile(tfModule.Path, profileName, destDir); err != nil {
			return count, err
		}
	}

	return count, nil
}

func processBlocks(g platform.Graph, blocks []platform.BlockID, tfModule *tfconfig.Module) (locals map[string]interface{}, fileIndex map[string][][2]int, blockCount int, err error) {
	locals = map[string]interface{}{}
	fileIndex = map[string][][2]int{}

	err = g.Walk(blocks, func(bID platform.BlockID) error {
		compType, compName := bID.ParseComponent()
		if compName != "" && compType == platform.BlockType_Local {
			// skip component inputs as they needs to be generated separately

			localVarName := platform.ComponentPrefix + compName
			locals[localVarName] = map[string]interface{}{}

			blockCount++
			return nil
		}

		b, found := bID.GetBlock(tfModule)
		if !found {
			return nil
		}

		relFilePath, err := getRelativeFilePath(tfModule.Path, b.GetPos().Filename)
		if err != nil {
			return eris.Wrap(err, "invalid updated module file path")
		}

		updateFileIndex(fileIndex, relFilePath, b)

		blockCount++
		return nil
	})

	return locals, fileIndex, blockCount, err
}

func getRelativeFilePath(basePath, targetPath string) (string, error) {
	return filepath.Rel(basePath, targetPath)
}

func updateFileIndex(fileIndex map[string][][2]int, relFilePath string, b platform.ParsedBlock) {
	pos := b.GetPos()
	fileIndex[relFilePath] = append(fileIndex[relFilePath], [2]int{pos.Line, pos.EndLine})

	if parentPosGetter, ok := b.(platform.BlockParentPosGetter); ok && parentPosGetter.GetParentPos() != nil {
		pPos := parentPosGetter.GetParentPos()
		fileIndex[relFilePath] = append(fileIndex[relFilePath], [2]int{pPos.Line, pPos.Line})       // add first line
		fileIndex[relFilePath] = append(fileIndex[relFilePath], [2]int{pPos.EndLine, pPos.EndLine}) // add last line
	}
}

func copyFilesFromIndex(srcDir, destDir string, fileIndex map[string][][2]int) error {
	for file, ranges := range fileIndex {
		err := copyLines(srcDir, destDir, file, ranges...)
		if err != nil {
			return eris.Wrapf(err, "failed to copy lines from file: %s", file)
		}
	}
	return nil
}

func writeLocalsToFile(locals map[string]interface{}, destDir string, apps app.Apps) error {
	for k := range locals {
		compName := strings.TrimPrefix(k, platform.ComponentPrefix)
		locals[k] = apps.GetDependenciesByType(compName).GetInputs()
	}

	localsFile, err := os.Create(path.Join(destDir, "tr_gen_locals.tf"))
	if err != nil {
		return eris.Wrapf(err, "error creating generated locals file")
	}
	defer localsFile.Close()

	return tfwriter.WriteLocals(locals, localsFile)
}

func readAllLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func isInRange(lineNum int, ranges ...[2]int) bool {
	for _, r := range ranges {
		if lineNum >= r[0] && lineNum <= r[1] {
			return true
		}
	}

	return false
}

// copyLines copies specific line ranges from the source file to the destination file.
// The rest of the lines in the destination file remain unchanged.
// If the destination file has fewer lines than expected, it inserts empty lines.
func copyLines(srcDir, destDir, relFile string, ranges ...[2]int) error {
	srcFile, destFile := filepath.Join(srcDir, relFile), filepath.Join(destDir, relFile)

	err := os.MkdirAll(filepath.Dir(destFile), constants.ReadWriteExecutePermissions)
	if err != nil {
		return eris.Wrapf(err, "failed to create directory for %s", destFile)
	}

	srcLines, err := readAllLines(srcFile)
	if err != nil {
		return eris.Wrap(err, "failed to read lines to copy")
	}

	destLines, err := readAllLines(destFile)
	if err != nil && !os.IsNotExist(err) {
		return eris.Wrap(err, "failed to read destination file")
	}

	output, err := os.Create(destFile)
	if err != nil {
		return eris.Wrap(err, "failed to create destination file")
	}
	defer output.Close()

	writer := bufio.NewWriter(output)

	for i, line := range srcLines {
		lineNum := i + 1
		if isInRange(lineNum, ranges...) {
			writer.WriteString(line + "\n")
		} else if lineNum <= len(destLines) {
			writer.WriteString(destLines[lineNum-1] + "\n")
		} else {
			writer.WriteString("\n")
		}
	}

	return writer.Flush()
}

func copyProfileConfigurationFile(moduleDirPath string, profileName string, codeDestDirPath string) error {
	sourcePath, err := getProfileVariableInputSourceFile(moduleDirPath, profileName)
	if err != nil {
		return eris.Wrapf(err, "could not retrieve configuration file for platform profile '%s'", profileName)
	}

	destPath, err := getProfileVariableInputDestFile(codeDestDirPath)
	if err != nil {
		return eris.Wrapf(err, "could not retrieve configuration target path for platform profile '%s'", profileName)
	}

	if copyFile(sourcePath, destPath); err != nil {
		return eris.Wrapf(err, "could not copy platform '%s' profile configuration", profileName)
	}

	return nil
}

func getProfileVariableInputSourceFile(moduleDirPath string, profileName string) (filePath string, err error) {
	profileSourceFileName := fmt.Sprintf("%s.%s", profileName, "tfvars")
	profileSourceFilePath := path.Join(moduleDirPath, profileSourceFileName)

	sourceFileStat, err := os.Stat(profileSourceFilePath)
	if os.IsNotExist(err) {
		return profileSourceFilePath, eris.Errorf("platform must define configuration for profile '%s' in terraform input file '%s'", profileName, profileSourceFilePath)
	} else if err != nil {
		return profileSourceFilePath, eris.Wrapf(err, "could not open platform profile configuration file '%s'", profileSourceFilePath)
	} else if !sourceFileStat.Mode().IsRegular() {
		return profileSourceFilePath, eris.Errorf("platform profile configuration path '%s' must point to a regular terraform input file", profileSourceFilePath)
	}

	return profileSourceFilePath, err
}

func getProfileVariableInputDestFile(destDirPath string) (filePath string, err error) {
	// any existing profile configuration file will be replaced - i.e. there is ever only one profile configuration applied
	profileDestFilePath := path.Join(destDirPath, "tr_gen_profile.auto.tfvars")

	_, err = os.Stat(profileDestFilePath)
	if err != nil && !os.IsNotExist(err) {
		return profileDestFilePath, eris.Wrapf(err, "could not access profile configuration target path '%s'", profileDestFilePath)
	}

	return profileDestFilePath, nil
}

func copyFile(srcPath string, dstPath string) error {
	source, err := os.Open(srcPath)
	if err != nil {
		return eris.Wrapf(err, "could not open copy source file '%s'", srcPath)
	}
	defer source.Close()

	destination, err := os.Create(dstPath)
	if err != nil {
		return eris.Wrapf(err, "could not create copy destination file '%s'", srcPath)
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		return eris.Wrapf(err, "failed to copy data from '%s' to '%s'", srcPath, dstPath)
	}

	return nil
}
