package db

import (
	"github.com/cldcvr/terrarium/api/pkg/pb/terrariumpb"
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type TFModuleAttribute struct {
	Model

	ModuleID                       uuid.UUID `gorm:"uniqueIndex:module_attribute_unique"`
	ModuleAttributeName            string    `gorm:"uniqueIndex:module_attribute_unique"`
	Description                    string
	RelatedResourceTypeAttributeID uuid.UUID
	Optional                       bool
	Computed                       bool

	Module            *TFModule            `gorm:"foreignKey:ModuleID"`
	ResourceAttribute *TFResourceAttribute `gorm:"foreignKey:RelatedResourceTypeAttributeID"` // Resource attribute with relates to this module attribute
}

type TFModuleAttributes []TFModuleAttribute

func (ma *TFModuleAttribute) GetConnectedModuleOutputs() TFModuleAttributes {
	if ma.ResourceAttribute != nil {
		return ma.ResourceAttribute.GetConnectedModuleOutputs()
	}

	return nil
}

func PopulateModuleAttrMappingsFilter(enable bool) FilterOption {
	if !enable {
		return NoOpFilter
	}

	return func(g *gorm.DB) *gorm.DB {
		return g.
			Where("computed = false").                                                                         // only pick input attributes, omit output attributes
			Preload("ResourceAttribute.OutputMappings.OutputAttribute.RelatedModuleAttrs", "computed = true"). // pick only module output attributes as valid references
			Preload("ResourceAttribute.OutputMappings.OutputAttribute.RelatedModuleAttrs.Module")              // load mapping of the input attribute to another module's output attribute
	}
}

func ModuleAttrSearchFilter(query string) FilterOption {
	if query == "" {
		return NoOpFilter
	}

	return func(g *gorm.DB) *gorm.DB {
		q := "%" + query + "%"
		return g.Where("module_attribute_name ILIKE ?", q)
	}
}

func ModuleAttrByIDsFilter(moduleId uuid.UUID, ids ...uuid.UUID) FilterOption {
	return func(g *gorm.DB) *gorm.DB {
		if len(ids) > 0 {
			g = g.Where("id in (?)", ids)
		}

		return g.Where("module_id = ?", moduleId)
	}
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set existing record ID in the given object
func (db *gDB) CreateTFModuleAttribute(e *TFModuleAttribute) (uuid.UUID, error) {
	return createOrUpdate(db.g(), e, []string{"module_id", "module_attribute_name"})
}

// QueryTFModuleAttributes based on the given filters
func (db *gDB) QueryTFModuleAttributes(filterOps ...FilterOption) (result TFModuleAttributes, err error) {
	q := db.g().Model(&TFModuleAttribute{})

	for _, filer := range filterOps {
		q = filer(q)
	}

	err = q.Order("module_attribute_name ASC").Find(&result).Error
	if err != nil {
		return nil, eris.Wrap(err, "query module attribute")
	}

	return
}

func (dbAttr TFModuleAttribute) ToProto() *terrariumpb.ModuleAttribute {
	resp := &terrariumpb.ModuleAttribute{
		Name:        dbAttr.ModuleAttributeName,
		Description: dbAttr.Description,
	}

	if dbAttr.Module != nil {
		resp.ParentModule = dbAttr.Module.ToProto()
	}

	outMoAttrs := dbAttr.GetConnectedModuleOutputs()
	if len(outMoAttrs) > 0 {
		resp.OutputModuleAttributes = outMoAttrs.ToProto()
	}

	return resp
}

func (dbAttrs TFModuleAttributes) ToProto() []*terrariumpb.ModuleAttribute {
	resp := make([]*terrariumpb.ModuleAttribute, len(dbAttrs))

	for i, dbAttr := range dbAttrs {
		resp[i] = dbAttr.ToProto()
	}

	return resp
}
