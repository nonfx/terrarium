// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

//go:build dbtest
// +build dbtest

package db_test

import (
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/xeipuuv/gojsonschema"
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
	uuidDep1        = uuid.MustParse("6ba7b81a-9dad-11d1-80b4-00c04fd430d2")
	uuidDep2        = uuid.MustParse("32e17e03-05c9-4c62-b244-84ece8682493")
	uuidDepAttr1    = uuid.MustParse("84c189b3-a85d-456a-abc0-eb9985ad3abf")
	uuidDepAttr2    = uuid.MustParse("76fb7d75-51cf-41fc-a390-b51753395b0c")
	uuidDepAttr3    = uuid.MustParse("eeddc51b-0afd-4d55-95bb-bf07a7b52016")
	uuidPlat1       = uuid.MustParse("c8bcaa19-f8de-4234-a0f5-4e930b61b037")
	uuidPlat2       = uuid.MustParse("b9696e52-1f44-4531-9612-f2dcd14f0972")
	uuidPlat1Comp1  = uuid.MustParse("f87cddf1-14f8-45f4-b737-9fe7702a937c")
	uuidPlat1Comp2  = uuid.MustParse("b07f92e8-e2d9-41f1-8fe8-b1bd444ecc5c")
	uuidPlat2Comp1  = uuid.MustParse("55dc2fab-9e32-4662-872d-d8fe8639bcda")
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

			Level1: "mockdata-l1",
			Level2: "mockdata-l2",
			Level3: "mockdata-l3",
			Level4: "mockdata-l4",
			Level5: "mockdata-l5",
			Level6: "mockdata-l6",
			Level7: "mockdata-l7",
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

			Level1: "mockdata-l1",
			Level2: "mockdata-l2",
			Level3: "mockdata-l3.2",
			Level4: "mockdata-l4.2",
			Level5: "mockdata-l5.2",
			Level6: "mockdata-l6.2",
			Level7: "mockdata-l7.2",
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
		Model:       db.Model{ID: uuidDep1},
		Title:       "dependency-1",
		InterfaceID: "dependency-1-interface",
		Description: "this is first test dependency",
		TaxonomyID:  uuidTax1,
		Attributes: db.DependencyAttributes{
			{
				Model:        db.Model{ID: uuidDepAttr1},
				DependencyID: uuidDep1,
				Name:         "dep-1-attr-1",
				Schema: &jsonschema.Node{
					Title:       "Attr 1",
					Description: "attribute 1",
					Type:        gojsonschema.TYPE_NUMBER,
				},
				Computed: false, //input
			},
			{
				Model:        db.Model{ID: uuidDepAttr2},
				DependencyID: uuidDep1,
				Name:         "dep-1-attr-2",
				Schema: &jsonschema.Node{
					Title:       "Attr 2",
					Description: "attribute 2",
					Type:        gojsonschema.TYPE_NUMBER,
				},
				Computed: true, //output
			},
			{
				Model:        db.Model{ID: uuidDepAttr3},
				DependencyID: uuidDep1,
				Name:         "dep-1-attr-3",
				Schema: &jsonschema.Node{
					Title:       "Attr 3",
					Description: "attribute 3",
					Type:        gojsonschema.TYPE_NUMBER,
				},
				Computed: true, //output
			},
		},
	},
	{
		Model:       db.Model{ID: uuidDep2},
		Title:       "dependency-2",
		InterfaceID: "dependency-2-interface",
		Description: "this is second test dependency",
		TaxonomyID:  uuidTax2,
		Attributes:  db.DependencyAttributes{},
	},
}

var platforms = db.Platforms{
	{
		Model:     db.Model{ID: uuidPlat1},
		Title:     "test-platform-1",
		CommitSHA: "2ed744403e50",
		Components: []db.PlatformComponent{
			{
				Model:        db.Model{ID: uuidPlat1Comp1},
				PlatformID:   uuidPlat1,
				DependencyID: uuidDep1,
			},
			{
				Model:        db.Model{ID: uuidPlat1Comp2},
				PlatformID:   uuidPlat1,
				DependencyID: uuidDep2,
			},
		},
	},
	{
		Model:     db.Model{ID: uuidPlat2},
		Title:     "test-platform-2",
		CommitSHA: "c4cf4e16e4f6",
		Components: []db.PlatformComponent{
			{
				Model:        db.Model{ID: uuidPlat2Comp1},
				PlatformID:   uuidPlat2,
				DependencyID: uuidDep1,
			},
		},
	},
}

func saveTestData(t *testing.T, g *gorm.DB) {
	_, err := db.AutoMigrate(g)
	require.NoError(t, err)

	err = g.Save(modules).Error
	require.NoError(t, err)

	err = g.Save(resourceMappings).Error
	require.NoError(t, err)

	err = g.Save(dependencies).Error
	require.NoError(t, err)

	err = g.Save(platforms).Error
	require.NoError(t, err)
}
