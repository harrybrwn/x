package nerdfont

import (
	"html/template"
	"io"
)

type GenerateTemplateData struct {
	Package  string
	Glyphs   []Glyph
	Metadata GlyphsMetadata

	// bool flags

	MappingFunction bool
	GenList         bool
}

func Generate(out io.Writer, data *GenerateTemplateData) error {
	tmpl, err := newTempl()
	if err != nil {
		return err
	}
	return tmpl.Execute(out, data)
}

func newTempl() (*template.Template, error) {
	return template.New("nerdfonts").Funcs(template.FuncMap{
		"ConstName": func(name string) string {
			return classToConstName(name)
		},
	}).Parse(templateString)
}

const templateString = `// Code generated. DO NOT EDIT.

// Package {{ .Package }} holds nerd font glyfs.
//
// Nerdfonts version {{ .Metadata.Version }} released at {{ .Metadata.Date }}
package {{ .Package }}

const (
{{- range .Glyphs }}
	// {{ .ConstName }} maps to "{{ .ID }}" (0x{{ .Code }}).{{- if .Alias }} Alias of "{{ .Alias }}".
	{{ .ConstName }} = {{ .Alias | ConstName }}
	{{- else }}
	{{ .ConstName }} = "{{ .Icon }}"
	{{- end }}
{{- end }}
)
{{- if .MappingFunction }}
func IconFromClassName(name string) (string, bool) {
	switch name {
{{- range .Glyphs }}
	case "{{ .ID }}":
		return {{ .ConstName }}, true
{{- end }}
	default:
		return "", false
	}
}
{{ end -}}
{{- if .GenList }}

// Glyph represents a nerdfont glyph.
type Glyph struct {
	Icon  string
	Class string
	Code   uint32
}

// Glyphs is a list of all nerdfont glyphs.
var Glyphs = [...]Glyph{
{{- range .Glyphs }}
	{Icon: {{ .ConstName }}, Class: "{{ .ID }}", Code: {{ .HexCode | printf "0x%x" }}},
{{- end }}
}
{{ end -}}
`
