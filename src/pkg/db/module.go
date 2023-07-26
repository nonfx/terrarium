package db

import (
	"strconv"
	"strings"

	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

type Version string

type TFModule struct {
	Model

	ModuleName  string
	Source      string  `gorm:"uniqueIndex:module_unique"`
	Version     Version `gorm:"uniqueIndex:module_unique"`
	Description string
	TaxonomyID  uuid.UUID `gorm:"default:null"`
	Namespace   string    `gorm:"default:farm_repo"`

	Taxonomy   *Taxonomy           `gorm:"foreignKey:TaxonomyID"`
	Attributes []TFModuleAttribute `gorm:"foreignKey:ModuleID"` // Attributes of the module
}

type TFModules []TFModule

func (m TFModule) GetInputAttributesWithMappings() (input TFModuleAttributes) {
	input = make(TFModuleAttributes, 0, len(m.Attributes))
	for _, a := range m.Attributes {
		if a.Computed {
			continue
		}

		if len(a.GetConnectedModuleOutputs()) > 0 {
			input = append(input, a)
		}
	}

	input = slices.Clip(input)
	return
}

func PopulateModuleMappingsFilter(enable bool) FilterOption {
	if !enable {
		return NoOpFilter
	}

	return func(g *gorm.DB) *gorm.DB {
		return g.
			Preload("Attributes", "computed = false").                                                                    // only pick input attributes, omit output attributes
			Preload("Attributes.ResourceAttribute.OutputMappings.OutputAttribute.RelatedModuleAttrs", "computed = true"). // pick only module output attributes as valid references
			Preload("Attributes.ResourceAttribute.OutputMappings.OutputAttribute.RelatedModuleAttrs.Module")              // load mapping of the input attribute to another module's output attribute
	}
}

func ModuleSearchFilter(query string) FilterOption {
	if query == "" {
		return NoOpFilter
	}

	return func(g *gorm.DB) *gorm.DB {
		q := "%" + query + "%"
		return g.Where("source ILIKE ? OR module_name ILIKE ?", q, q)
	}
}

func ModuleNamespaceFilter(namespace []string) FilterOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("namespace IN ?", namespace)
	}
}

func ModuleByIDsFilter(ids ...uuid.UUID) FilterOption {
	if len(ids) == 0 {
		return NoOpFilter
	}

	return func(g *gorm.DB) *gorm.DB {
		return g.Where("id in (?)", ids)
	}
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set existing record ID in the given object
func (db *gDB) CreateTFModule(e *TFModule) (uuid.UUID, error) {
	return createOrUpdate(db.g(), e, []string{"source", "version"})
}

// QueryTFModules based on the given filters
func (db *gDB) QueryTFModules(filterOps ...FilterOption) (result TFModules, err error) {
	q := db.g().Model(&TFModule{})

	for _, filer := range filterOps {
		q = filer(q)
	}

	err = q.Order("source ASC, version DESC").Find(&result).Error
	if err != nil {
		return nil, eris.Wrap(err, "query module")
	}

	return
}

// FindOutputMappingsByModuleID DEPRECATED fetch the terraform module(s) by id along with it's input-attribute and the mappings to another module's outputs.
// The returned module contains attributes (input only) populated along with the mapping of each module also populated.
//
// Path:
//
// - TFModules.Attributes = array of input attribute objects
//
// - TFModules.Attributes.ResourceAttribute = Resource attribute related to the current module attribute
//
// - TFModules.Attributes.ResourceAttribute.OutputMappings = Array of resource output attribute mappings
//
// - Attributes.ResourceAttribute.OutputMappings.OutputAttribute.RelatedModuleAttrs = Array of other module's attributes related to the resource output attribute
//
// - Attributes.ResourceAttribute.OutputMappings.OutputAttribute.RelatedModuleAttrs.Module = Module object of the module attribute being mapped.
func (db *gDB) FindOutputMappingsByModuleID(ids ...uuid.UUID) (result TFModules, err error) {
	err = db.g().
		Where("id in (?)", ids).
		Preload("Attributes", "computed = false").                                                                    // only pick input attributes, omit output attributes
		Preload("Attributes.ResourceAttribute.OutputMappings.OutputAttribute.RelatedModuleAttrs", "computed = true"). // pick only module output attributes as valid references
		Preload("Attributes.ResourceAttribute.OutputMappings.OutputAttribute.RelatedModuleAttrs.Module").             // load mapping of the input attribute to another module's output attribute
		Find(&result).
		Error

	return
}

func (m TFModule) ToProto() *terrariumpb.Module {
	inp := m.GetInputAttributesWithMappings()
	return &terrariumpb.Module{
		Id:              m.ID.String(),
		TaxonomyId:      m.TaxonomyID.String(),
		ModuleName:      m.ModuleName,
		Source:          m.Source,
		Version:         string(m.Version),
		Description:     m.Description,
		InputAttributes: inp.ToProto(),
		Namespace:       m.Namespace,
	}
}

func (mArr TFModules) ToProto() []*terrariumpb.Module {
	res := make([]*terrariumpb.Module, len(mArr))
	for i, m := range mArr {
		res[i] = m.ToProto()
	}

	return res
}

// Compare returns -1 if v1 is less then v2, 0 if both are equal and 1 if v1 is greater then v2.
func (v1 Version) Compare(v2 Version) int {
	parts1 := strings.Split(string(v1), ".")
	parts2 := strings.Split(string(v2), ".")

	for i := 0; i < len(parts1) || i < len(parts2); i++ {
		num1 := 0
		num2 := 0

		if i < len(parts1) {
			num1, _ = strconv.Atoi(parts1[i])
		}

		if i < len(parts2) {
			num2, _ = strconv.Atoi(parts2[i])
		}

		if num1 < num2 {
			return -1
		} else if num1 > num2 {
			return 1
		}
		// if this parts is equal, compare the next part
	}

	return 0
}
