package generate

type BlockType int

const (
	BlockType_Unknown = iota
	BlockType_Module
	BlockType_Resource
	BlockType_Data
	BlockType_Variable
	BlockType_Output
	BlockType_Local
)

type BlockLocation struct {
	File      string
	LineStart int
	LineEnd   int
}

type TFModule struct{}

type TFResource struct{}

type TFBlock struct {
	Type        BlockType
	Name        string
	Location    BlockLocation
	ParsedValue interface{}
	DependsOn   TFBlocks
}

type TFBlocks []*TFBlock

func (b TFBlocks) Parse(platformDir string) {

}

func (b TFBlocks) GetByTypeId(bType BlockType, id string) *TFBlock {

	return nil
}

func (b TFBlocks) ProcessAppDependencies(deps AppDependencies) {

}

func (b TFBlocks) Render(outputDir string) error {

	return nil
}
