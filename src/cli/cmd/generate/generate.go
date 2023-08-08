package generate

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/cldcvr/terrarium/src/pkg/metadata/platform"
	"github.com/rotisserie/eris"
)

func getTFBlockPos(bId platform.BlockID, m *tfconfig.Module) *tfconfig.SourcePos {
	t, bn := bId.Parse()

	switch t {
	case platform.BlockType_ModuleCall:
		b := m.ModuleCalls[bn]
		if b != nil {
			return &b.Pos
		}
	case platform.BlockType_Resource:
		b := m.ManagedResources[bn]
		if b != nil {
			return &b.Pos
		}
	case platform.BlockType_Data:
		b := m.DataResources[bn]
		if b != nil {
			return &b.Pos
		}
	case platform.BlockType_Variable:
		b := m.Variables[bn]
		if b != nil {
			return &b.Pos
		}
	case platform.BlockType_Output:
		b := m.Outputs[bn]
		if b != nil {
			return &b.Pos
		}
	}

	return nil
}

func fetchBlockString(pos *tfconfig.SourcePos) ([]byte, error) {
	b, err := os.ReadFile(pos.Filename)
	if err != nil {
		return nil, eris.Wrapf(err, "os file (%s) read error", pos.Filename)
	}

	if len(b) <= pos.EndByte {
		return nil, eris.Errorf("block position is not found in the file: %s:%d-%d", pos.Filename, pos.Line, pos.EndLine)
	}

	return b[pos.StartByte:pos.EndByte], nil
}

func blocksToPull(g platform.Graph, components ...string) []platform.BlockID {
	blockIds := []platform.BlockID{}
	for _, comp := range components {
		bId := platform.NewBlockID(platform.BlockType_ModuleCall, platform.ComponentPrefix+comp)
		blockIds = append(blockIds, bId)
	}

	return blockIds
}

func writeTF(g platform.Graph, destDir string, blocks []platform.BlockID, tfModule *tfconfig.Module) (blockCount int, err error) {
	return blockCount, g.Walk(blocks, func(bId platform.BlockID) error {
		pos := getTFBlockPos(bId, tfModule)
		if pos == nil {
			return nil
		}

		err := writeBlockToFile(destDir, tfModule.Path, pos)
		if err != nil {
			return err
		}

		blockCount++
		return nil
	})
}

func writeBlockToFile(destRootDir, srcDir string, tfBlockPos *tfconfig.SourcePos) error {
	relFilePath, err := filepath.Rel(srcDir, tfBlockPos.Filename)
	if err != nil {
		return eris.Wrapf(err, "failed to get relative path of file: %s from dir: %s", tfBlockPos.Filename, srcDir)
	}

	destFilePath := filepath.Join(destRootDir, relFilePath)
	destDirPath := filepath.Dir(destFilePath)

	err = os.MkdirAll(destDirPath, 0755)
	if err != nil {
		return eris.Wrapf(err, "failed to create directory at %s", destDirPath)
	}

	fileBytes, err := os.ReadFile(destFilePath)
	if os.IsNotExist(err) {
		err = nil
		fileBytes = []byte{}
	}

	nl := []byte("\n")

	lines := bytes.Split(fileBytes, nl)

	if len(lines) < tfBlockPos.Line-1 {
		lines = append(lines, make([][]byte, tfBlockPos.Line-len(lines))...)
	}

	beforeLines := lines[:tfBlockPos.Line-1]
	afterLines := [][]byte{{}}
	if len(lines) > tfBlockPos.EndLine {
		afterLines = lines[tfBlockPos.EndLine:]
	}

	blockCode, err := fetchBlockString(tfBlockPos)
	if err != nil {
		return eris.Wrapf(err, "failed to read hcl block")
	}

	curBlockLines := bytes.Split(blockCode, nl)

	finalLines := append(beforeLines, curBlockLines...)
	finalLines = append(finalLines, afterLines...)

	err = os.WriteFile(destFilePath, bytes.Join(finalLines, nl), 0644)
	if err != nil {
		return eris.Wrapf(err, "failed to write file at %s", destFilePath)
	}

	return nil
}
