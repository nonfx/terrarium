// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package platform

import (
	"os"
	"path"
	"strings"

	"github.com/cldcvr/terraform-config-inspect/tfconfig"
)

const (
	profileFileSuffix = ".tfvars"
)

func (pArr *Profiles) Parse(platformModule *tfconfig.Module) {
	if list, err := os.ReadDir(platformModule.Path); err == nil {
		for _, item := range list {
			if item.Type().IsRegular() {
				if baseName, found := strings.CutSuffix(item.Name(), profileFileSuffix); found {
					p := pArr.GetByID(baseName)
					if p == nil {
						p = pArr.Append(Profile{ID: baseName})
					}

					if doc, err := GetDoc(path.Join(platformModule.Path, item.Name()), -1, false); err == nil {
						SetValueFromDocIfFound(&p.Title, docCommentTitleArgTag, doc)
						SetValueFromDocIfFound(&p.Description, docCommentDescArgTag, doc)
					}
				}
			}
		}
	}
}

func (pArr Profiles) GetByID(id string) *Profile {
	for i, v := range pArr {
		if v.ID == id {
			return &pArr[i]
		}
	}

	return nil
}

func (pArr *Profiles) Append(p Profile) *Profile {
	(*pArr) = append((*pArr), p)
	return &(*pArr)[len(*pArr)-1]
}
