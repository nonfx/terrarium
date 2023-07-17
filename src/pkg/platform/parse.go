package platform

import (
	"fmt"

	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/rotisserie/eris"
	"golang.org/x/exp/maps"
)

func Parse(dir string) error {
	module, diags := tfconfig.LoadModule(dir, &tfconfig.ResolvedModulesSchema{})
	if diags.HasErrors() {
		return eris.Wrap(diags.Err(), "terraform code has errors")
	}

	newModule := tfconfig.NewModule(dir) //(path.Join(dir, ".terrarium"))
	newModule.ModuleCalls = map[string]*tfconfig.ModuleCall{
		// COMPONENT_PREFIX + "postgres": nil,
	}
	ExtractBlocksInto(module, newModule)

	return nil
}

// ExtractBlocksInto extracts the block keys with nil value in destModule from srcModule along with it's dependencies recursively.
func ExtractBlocksInto(srcModule, destModule *tfconfig.Module) {
	setRecurse := false
	for k, v := range destModule.ModuleCalls {
		if v != nil {
			continue // value already set
		}

		srcVal, available := srcModule.ModuleCalls[k]
		if !available || srcVal == nil {
			delete(destModule.ModuleCalls, k)
			continue
		}

		destModule.ModuleCalls[k] = srcVal
		setRecurse = true

		dependencies := GetDependenciesFromInputs(srcVal.Inputs)
		ModuleExtend(destModule, dependencies, false)
	}

	for k, v := range destModule.ManagedResources {
		if v != nil {
			continue // value already set
		}

		srcVal, available := srcModule.ManagedResources[k]
		if !available || srcVal == nil {
			delete(destModule.ManagedResources, k)
			continue
		}

		destModule.ManagedResources[k] = srcVal
		setRecurse = true

		dependencies := GetDependenciesFromInputs(srcVal.Inputs)
		ModuleExtend(destModule, dependencies, false)
	}

	for k, v := range destModule.DataResources {
		if v != nil {
			continue // value already set
		}

		srcVal, available := srcModule.DataResources[k]
		if !available || srcVal == nil {
			delete(destModule.DataResources, k)
			continue
		}

		destModule.DataResources[k] = srcVal
		setRecurse = true

		dependencies := GetDependenciesFromInputs(srcVal.Inputs)
		ModuleExtend(destModule, dependencies, false)
	}

	for k, v := range destModule.ProviderConfigs {
		if v != nil {
			continue // value already set
		}

		srcVal, available := srcModule.ProviderConfigs[k]
		if !available || srcVal == nil {
			delete(destModule.ProviderConfigs, k)
			continue
		}

		destModule.ProviderConfigs[k] = srcVal
		setRecurse = true
	}

	for k, v := range destModule.RequiredProviders {
		if v != nil {
			continue // value already set
		}

		srcVal, available := srcModule.RequiredProviders[k]
		if !available || srcVal == nil {
			delete(destModule.RequiredProviders, k)
			continue
		}

		destModule.RequiredProviders[k] = srcVal
		setRecurse = true
	}

	for k, v := range destModule.Variables {
		if v != nil {
			continue // value already set
		}

		srcVal, available := srcModule.Variables[k]
		if !available || srcVal == nil {
			delete(destModule.Variables, k)
			continue
		}

		destModule.Variables[k] = srcVal
		setRecurse = true
	}

	if setRecurse {
		ExtractBlocksInto(srcModule, destModule)
	} else {
		maps.Copy(destModule.Outputs, srcModule.Outputs)
		KeepRelevantOutputs(destModule)
	}
}

func MapExtend[M ~map[K]V, K comparable, V any](dest, src M, override bool) {
	for srcK, srcV := range src {
		if _, isSet := dest[srcK]; !isSet || override {
			dest[srcK] = srcV
		}
	}
}

func ModuleExtend(dest, src *tfconfig.Module, override bool) {
	MapExtend(dest.Locals, src.Locals, override)
	MapExtend(dest.Variables, src.Variables, override)
	MapExtend(dest.Outputs, src.Outputs, override)
	MapExtend(dest.Inputs, src.Inputs, override)
	MapExtend(dest.RequiredProviders, src.RequiredProviders, override)
	MapExtend(dest.ProviderConfigs, src.ProviderConfigs, override)
	MapExtend(dest.ManagedResources, src.ManagedResources, override)
	MapExtend(dest.DataResources, src.DataResources, override)
	MapExtend(dest.ModuleCalls, src.ModuleCalls, override)
}

// GetDependencies returned Module object has keys for the dependencies being set with null values.
func GetDependenciesFromInputs(inputs map[string]tfconfig.AttributeReference) (dependencies *tfconfig.Module) {
	dependencies = tfconfig.NewModule(".")
	for _, v := range inputs {
		switch v.Type() {
		case "":
		case "module":
			dependencies.ModuleCalls[v.Name()] = nil
		case "local":
			dependencies.Locals[v.Name()] = nil
		case "var":
			dependencies.Variables[v.Name()] = nil
		default:
			resID := fmt.Sprintf("%s.%s", v.Type(), v.Name())
			dependencies.ManagedResources[resID] = nil
			dependencies.DataResources["data."+resID] = nil
		}
	}

	return
}

func KeepRelevantOutputs(m *tfconfig.Module) {
	for k, v := range m.Outputs {
		keep := true
		switch v.Value.Type() {
		case "", "local":
		case "module":
			_, keep = m.ModuleCalls[v.Value.Name()]
		case "resource":
			_, keep = m.ManagedResources[v.Value.Name()]
		case "data":
			_, keep = m.DataResources[v.Value.Name()]
		case "var":
			_, keep = m.Variables[v.Value.Name()]
		default:
			resID := fmt.Sprintf("%s.%s", v.Value.Type(), v.Value.Name())
			_, keep1 := m.ManagedResources[resID]
			_, keep2 := m.DataResources["data."+resID]
			keep = keep1 || keep2
		}

		if !keep {
			delete(m.Outputs, k)
		}
	}
}
