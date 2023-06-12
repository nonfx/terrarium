{{- define "module_attribute" -}}
  {{- $hasLink := false}}
  {{- $padding := "\n  "}}
  {{- if and (not (eq .ResourceAttribute nil)) (gt (len .ResourceAttribute.OutputMappings) 0) }}
    {{- $resOutputMapping := (index .ResourceAttribute.OutputMappings 0) }}
    {{- if gt (len $resOutputMapping.OutputAttribute.RelatedModuleAttrs) 0 }}
      {{- $hasLink = true }}
      {{- $linkedModule := index $resOutputMapping.OutputAttribute.RelatedModuleAttrs 0}}
      {{- print $padding}}{{.ModuleAttributeName}} = module.{{$linkedModule.Module.ModuleName}}.{{$linkedModule.ModuleAttributeName}}
    {{- end}}
  {{- end}}
  {{- if and (not .Optional) (not $hasLink) }}
    {{- print $padding}}{{- .ModuleAttributeName}} =
  {{- end}}
{{- end -}}

{{- define "module_call" -}}
## TERRARIUM MODULE ## {{.Source}} ## {{.ModuleName}} ## _TAXONOMY_ ##
module "{{.ModuleName}}" {
  source = "{{.Source}}"
  {{- if ne .Version "" -}}
    {{"\n  "}}version = "{{.Version}}"
  {{- end}}

  {{- range .Attributes}}
    {{- template "module_attribute" .}}
  {{- end}}
}

{{end -}}
