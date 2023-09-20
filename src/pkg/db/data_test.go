// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

//go:build dbtest
// +build dbtest

package db_test

import (
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var (
	uuidMod1        = uuid.MustParse("f47ac10b-58cc-4372-a567-0e02b2c3d479")
	uuidTax1        = uuid.MustParse("e6fb062d-74d6-4491-80bc-5d2c8e6d9ebb")
	uuidMod1Attr1   = uuid.MustParse("a1d0c6e8-3b5d-4cc5-bfef-3dadb3e8f211")
	uuidMod1Attr2   = uuid.MustParse("7e57d004-2b97-44e7-8f03-66e975e08855")
	uuidMod1Attr3   = uuid.MustParse("5a8dd3ad-a524-479f-8f7a-4388c9e6d3c2")
	uuidMod2        = uuid.MustParse("d3c1d35c-47a9-4837-add4-0e30db0f2f3b")
	uuidTax2        = uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	uuidMod2Attr1   = uuid.MustParse("8e281e37-7e17-4b35-bdd3-83933f2badcc")
	uuidMod2Attr2   = uuid.MustParse("2d7c7a8a-df6b-4c1b-8ba8-d9db4d709456")
	uuidMod2Attr3   = uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	uuidRes1Attr1   = uuid.MustParse("6ba7b811-9dad-11d1-80b4-00c04fd430c9")
	uuidRes1        = uuid.MustParse("6ba7b812-9dad-11d1-80b4-00c04fd430ca")
	uuidProv1       = uuid.MustParse("6ba7b813-9dad-11d1-80b4-00c04fd430cb")
	uuidRes2Attr2   = uuid.MustParse("6ba7b814-9dad-11d1-80b4-00c04fd430cc")
	uuidRes1Attr3   = uuid.MustParse("6ba7b815-9dad-11d1-80b4-00c04fd430cd")
	uuidRes2        = uuid.MustParse("6ba7b816-9dad-11d1-80b4-00c04fd430ce")
	uuidRes2Attr1   = uuid.MustParse("6ba7b817-9dad-11d1-80b4-00c04fd430cf")
	uuidRes1Attr2   = uuid.MustParse("6ba7b818-9dad-11d1-80b4-00c04fd430d0")
	uuidRes2Attr3   = uuid.MustParse("6ba7b819-9dad-11d1-80b4-00c04fd430d1")
	uuidResMapping1 = uuid.MustParse("12345678-1234-5678-1234-567812345678")
	uuidResMapping2 = uuid.MustParse("87654321-8765-4321-8765-432187654321")
)

var modules = db.TFModules{
	db.TFModule{
		Model: db.Model{ID: uuidMod1},

		ModuleName:  "module-1",
		Source:      "module-1-source-1",
		Version:     "1.1",
		Description: "this is first test module",
		Namespace:   "unit-test",
		TaxonomyID:  uuidTax1,

		Taxonomy: &db.Taxonomy{
			Model: db.Model{ID: uuidTax1},

			Level1: "module-1-l1",
			Level2: "module-1-l2",
			Level3: "module-1-l3",
			Level4: "module-1-l4",
			Level5: "module-1-l5",
			Level6: "module-1-l6",
			Level7: "module-1-l7",
		},

		Attributes: []db.TFModuleAttribute{
			{
				Model: db.Model{ID: uuidMod1Attr1},

				ModuleID:                       uuidMod1,
				ModuleAttributeName:            "module-1-attr-1",
				Description:                    "first attribute of the first module",
				Optional:                       false,
				Computed:                       false,
				RelatedResourceTypeAttributeID: uuidRes1Attr1,

				ResourceAttribute: &db.TFResourceAttribute{
					Model:          db.Model{ID: uuidRes1Attr1},
					AttributePath:  "res-1-attr-1",
					Computed:       false,
					ResourceTypeID: uuidRes1,
					ProviderID:     uuidProv1,

					ResourceType: db.TFResourceType{
						Model:        db.Model{ID: uuidRes1},
						ProviderID:   uuidProv1,
						ResourceType: "res-1",
						Provider: db.TFProvider{
							Model: db.Model{ID: uuidProv1},
							Name:  "mocked-unit-test",
						},
					},
				},
			},
			{
				Model: db.Model{ID: uuidMod1Attr2},

				ModuleID:                       uuidMod1,
				ModuleAttributeName:            "module-1-attr-2",
				Description:                    "second attribute of the first module",
				Optional:                       true,
				Computed:                       false,
				RelatedResourceTypeAttributeID: uuidRes2Attr2,

				ResourceAttribute: &db.TFResourceAttribute{
					Model:          db.Model{ID: uuidRes2Attr2},
					AttributePath:  "res-2-attr-2",
					Optional:       true,
					Computed:       false,
					ResourceTypeID: uuidRes2,
					ProviderID:     uuidProv1,

					ResourceType: db.TFResourceType{
						Model:        db.Model{ID: uuidRes2},
						ProviderID:   uuidProv1,
						ResourceType: "res-2",
					},
				},
			},
			{
				Model: db.Model{ID: uuidMod1Attr3},

				ModuleID:                       uuidMod1,
				ModuleAttributeName:            "module-1-attr-3",
				Description:                    "first output attribute of the first module",
				Computed:                       true,
				RelatedResourceTypeAttributeID: uuidRes1Attr3,

				ResourceAttribute: &db.TFResourceAttribute{
					Model:          db.Model{ID: uuidRes1Attr3},
					AttributePath:  "res-1-attr-3",
					Optional:       true,
					Computed:       true,
					ResourceTypeID: uuidRes1,
					ProviderID:     uuidProv1,
				},
			},
		},
	},
	db.TFModule{
		Model: db.Model{ID: uuidMod2},

		ModuleName:  "module-2",
		Source:      "module-2-source-1",
		Version:     "1.1",
		Description: "this is second test module",
		Namespace:   "unit-test",
		TaxonomyID:  uuidTax2,

		Taxonomy: &db.Taxonomy{
			Model: db.Model{ID: uuidTax2},

			Level1: "module-2-l1",
			Level2: "module-2-l2",
			Level3: "module-2-l3",
			Level4: "module-2-l4",
			Level5: "module-2-l5",
			Level6: "module-2-l6",
			Level7: "module-2-l7",
		},

		Attributes: []db.TFModuleAttribute{
			{
				Model: db.Model{ID: uuidMod2Attr1},

				ModuleID:                       uuidMod2,
				ModuleAttributeName:            "module-2-attr-1",
				Description:                    "first attribute of the second module",
				Optional:                       false,
				Computed:                       false,
				RelatedResourceTypeAttributeID: uuidRes2Attr1,

				ResourceAttribute: &db.TFResourceAttribute{
					Model:          db.Model{ID: uuidRes2Attr1},
					AttributePath:  "res-2-attr-1",
					Optional:       false,
					Computed:       false,
					ResourceTypeID: uuidRes2,
					ProviderID:     uuidProv1,
				},
			},
			{
				Model: db.Model{ID: uuidMod2Attr2},

				ModuleID:                       uuidMod2,
				ModuleAttributeName:            "module-2-attr-2",
				Description:                    "second attribute of the second module",
				Optional:                       true,
				Computed:                       false,
				RelatedResourceTypeAttributeID: uuidRes1Attr2,

				ResourceAttribute: &db.TFResourceAttribute{
					Model:          db.Model{ID: uuidRes1Attr2},
					AttributePath:  "res-1-attr-2",
					Computed:       false,
					ResourceTypeID: uuidRes1,
					ProviderID:     uuidProv1,
				},
			},
			{
				Model: db.Model{ID: uuidMod2Attr3},

				ModuleID:                       uuidMod2,
				ModuleAttributeName:            "module-2-attr-3",
				Description:                    "first output attribute of the second module",
				Computed:                       true,
				RelatedResourceTypeAttributeID: uuidRes2Attr3,

				ResourceAttribute: &db.TFResourceAttribute{
					Model:          db.Model{ID: uuidRes2Attr3},
					AttributePath:  "res-2-attr-3",
					Optional:       true,
					Computed:       true,
					ResourceTypeID: uuidRes2,
					ProviderID:     uuidProv1,
				},
			},
		},
	},
}

var resourceMappings = []db.TFResourceAttributesMapping{
	{
		Model:             db.Model{ID: uuidResMapping1},
		InputAttributeID:  uuidRes1Attr1,
		OutputAttributeID: uuidRes2Attr3,
	},
	{
		Model:             db.Model{ID: uuidResMapping2},
		InputAttributeID:  uuidRes2Attr1,
		OutputAttributeID: uuidRes1Attr3,
	},
}

var dependencies = []db.Dependency{
	{
		Model:       db.Model{ID: uuid.MustParse("6ba7b81a-9dad-11d1-80b4-00c04fd430d2")},
		Title:       "dependency-1",
		InterfaceID: "dependency-1-interface",
		Description: "this is first test dependency",
		TaxonomyID:  uuidTax1,
	},
}

func saveTestData(t *testing.T, g *gorm.DB) {
	t.Helper()

	_, err := db.AutoMigrate(g)
	require.NoError(t, err)

	err = g.Save(modules).Error
	require.NoError(t, err)

	err = g.Save(resourceMappings).Error
	require.NoError(t, err)

	err = g.Save(dependencies).Error
	require.NoError(t, err)
}
