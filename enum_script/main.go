package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	scale_codec "github.com/crypto2lab/scale-codec"
)

const defaultPackage = "main"
const outputExt = ".go"
const inputExt = ".scale"

var alphabet string = "abcdefghijklmnopqrstuvwxyz"

func isScaleFile(filename string) bool {
	return filepath.Ext(filename) == inputExt
}

func removeExtension(filename string) string {
	ext := filepath.Ext(filename)
	return filename[:len(filename)-len(ext)]
}

func parseMap(list string, f func(int, string) string) []string {
	output := make([]string, len(list))
	for idx, item := range list {
		output[idx] = f(idx, string(item))
	}
	return output
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Error: expected only two argument: scale file and package")
	}

	outputPackage := defaultPackage
	cliProvidedPackage := strings.TrimSpace(os.Args[2])
	if cliProvidedPackage != "" {
		outputPackage = cliProvidedPackage
	}

	finfo, err := os.Stat(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if finfo.IsDir() {
		log.Fatalf("Error: directories are not currently supported")
	}

	if !isScaleFile(finfo.Name()) {
		log.Fatalf("Error: expected a .scale file")
	}

	contents, err := os.ReadFile(finfo.Name())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Parsing %v file\n", finfo.Name())
	result := scale_codec.ParseEnum(finfo.Name(), bytes.NewReader(contents))
	if result != 0 {
		log.Fatalf("Error: failed to parse %s file", finfo.Name())
	}

	generatedEnums := parseEnumsDefinition(outputPackage, scale_codec.Enums)
	outputFile := strings.Join([]string{removeExtension(finfo.Name()), outputExt}, "")
	err = os.WriteFile(outputFile, []byte(generatedEnums), os.ModePerm)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("file generated: %s\n", outputFile)
}

func parseEnumsDefinition(pacakge string, enums []scale_codec.Enum) string {
	type enumDefinition struct {
		EnumName string
		Variants []string
	}

	enumTemplate, err := template.New("enums_definitions").Parse(EnumDefinitionTemplate)
	if err != nil {
		log.Fatalf("Parsing template error: %v", err)
	}

	enumsDefinitions := new(strings.Builder)
	for _, enum := range enums {
		variantsName := make([]string, len(enum.Variants))
		for idx, variant := range enum.Variants {
			variantsName[idx] = variant.Name
		}

		value := enumDefinition{
			EnumName: enum.Name,
			Variants: variantsName,
		}

		err := enumTemplate.Execute(enumsDefinitions, value)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		enumsDefinitions.WriteRune('\n')
	}

	fileTemplate, err := template.New("enum_file_template").Parse(EnumFileTemplate)
	if err != nil {
		log.Fatalf("Parsing template error: %v", err)
	}

	type fileTemplateValue struct {
		Package                 string
		GenericTupleDefinitions string
		EnumsDefinitions        string
		VariantsDefinitions     string
	}

	value := fileTemplateValue{
		Package:                 pacakge,
		GenericTupleDefinitions: parseGenericTupleDefinitions(scale_codec.GenericTuple),
		EnumsDefinitions:        enumsDefinitions.String(),
		VariantsDefinitions:     parseVariantsDefinitions(enums),
	}

	fileBuffer := new(strings.Builder)
	err = fileTemplate.Execute(fileBuffer, value)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return fileBuffer.String()
}

func parseGenericTupleDefinitions(genericTuples map[string]int) string {
	builder := &strings.Builder{}

	for structName, arity := range genericTuples {
		t, err := template.New("generic_tuple_definitions").Parse(GenericTupleStructTemplate)
		if err != nil {
			log.Fatalf("Parsing template error: %v", err)
		}

		type genericTuple struct {
			GenericTupleName        string
			GenericArity            string
			GenericTupleFields      string
			GenericNames            string
			Fields                  []string
			UnmarshalFuncSignatures string
			FuncsAndFields          map[string]string
		}

		generics := alphabet[0:arity]
		genericArity := parseMap(generics, func(_ int, s string) string {
			return fmt.Sprintf("%s scale_codec.Marshaler",
				strings.ToUpper(string(s)))
		})

		genericNames := parseMap(generics, func(_ int, s string) string {
			return strings.ToUpper(string(s))
		})

		genericFields := parseMap(generics, func(idx int, s string) string {
			return fmt.Sprintf("\tF%d %s", idx, strings.ToUpper(string(s)))
		})

		unmarshalFuncSignatures := parseMap(generics, func(i int, s string) string {
			return fmt.Sprintf("func%s func (io.Reader) (%s, error)",
				strings.ToUpper(string(s)),
				strings.ToUpper(string(s)))
		})

		unmarshalFuncsNames := parseMap(generics, func(i int, s string) string {
			return fmt.Sprintf("func%s", strings.ToUpper(string(s)))
		})

		fieldsNames := make([]string, arity)
		for idx := range fieldsNames {
			fieldsNames[idx] = fmt.Sprintf("F%d", idx)
		}

		funcsAndFields := make(map[string]string, len(unmarshalFuncsNames))
		for idx, funcName := range unmarshalFuncsNames {
			funcsAndFields[funcName] = fieldsNames[idx]
		}

		value := genericTuple{
			GenericTupleName:        structName,
			GenericArity:            strings.Join(genericArity, ","),
			GenericTupleFields:      strings.Join(genericFields, "\n"),
			GenericNames:            strings.Join(genericNames, ","),
			Fields:                  fieldsNames,
			UnmarshalFuncSignatures: strings.Join(unmarshalFuncSignatures, ","),
			FuncsAndFields:          funcsAndFields,
		}

		err = t.Execute(builder, value)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		builder.WriteRune('\n')
	}

	return builder.String()
}

func parseVariantsDefinitions(parsedEnums []scale_codec.Enum) string {
	t, err := template.New("variants_definitions").Parse(EnumVariantDefinitionTempate)
	if err != nil {
		log.Fatalf("Parsing template error: %v", err)
	}

	type variant struct {
		EnumName        string
		Name            string
		Type            string
		TypeConstructor string
		UnmarshalSCALE  string
		Index           int
	}

	variantsDefs := new(strings.Builder)
	for _, enum := range parsedEnums {
		for index, vari := range enum.Variants {
			unmarshalScale := defaultUnmarshalSCALE
			if strings.TrimSpace(vari.UnmarshalScale) != "" {
				unmarshalScale = strings.TrimSpace(vari.UnmarshalScale)
			}

			value := variant{
				EnumName:        enum.Name,
				Name:            vari.Name,
				Type:            vari.Type,
				TypeConstructor: vari.TypeConstructor,
				Index:           index,
				UnmarshalSCALE:  unmarshalScale,
			}
			err := t.Execute(variantsDefs, value)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			variantsDefs.WriteRune('\n')
		}
	}

	return variantsDefs.String()
}

const EnumFileTemplate = `// Code generated by scale_codec/enum_script. DO NOT EDIT.
package {{.Package}}

import (
	"bytes"
	"fmt"
	"io"

	scale_codec "github.com/crypto2lab/scale-codec"
)

{{ .GenericTupleDefinitions }}

{{ .EnumsDefinitions }}

{{ .VariantsDefinitions }}`

const EnumDefinitionTemplate = `type {{ .EnumName }} interface {
	scale_codec.Encodable
	Is{{ .EnumName }}()
}

func Unmarshal{{ .EnumName }}(reader io.Reader) ({{ .EnumName }}, error) {
	enumTag := make([]byte, 1)
	n, err := reader.Read(enumTag)
	if err != nil {
		return nil, err
	}

	if n != 1 {
		return nil, fmt.Errorf("%w: got %v", scale_codec.ErrExpectedOneByteRead, n)
	}

	switch enumTag[0] {
	{{ range $i, $a := .Variants }}
	case {{ $a }}Index:
		unmarshaler := New{{ $a }}()
		err := unmarshaler.UnmarshalSCALE(reader)
		if err != nil {
			return nil, err
		}
		return unmarshaler, err
	{{ end }}
	default:
		return nil, fmt.Errorf("unexpected enum tag: %v", enumTag[0])
	}
}`

const defaultUnmarshalSCALE = "return i.Inner.UnmarshalSCALE(reader)"
const EnumVariantDefinitionTempate = `var {{ .Name }}Index byte = {{ .Index }}

var _ {{ .EnumName }} = (*{{ .Name }})(nil)

type {{ .Name }} struct {
	Inner {{ .Type }}
}

func New{{ .Name }}() *{{ .Name }} {
	return &{{ .Name }}{
		Inner: {{ .TypeConstructor }},
	}
}

func ({{ .Name }}) Is{{ .EnumName }}() {}

func (i {{ .Name }}) MarshalSCALE() ([]byte, error) {
	innerEncode, err := i.Inner.MarshalSCALE()
	if err != nil {
		return nil, err
	}

	idx := {{ .Name }}Index
	return bytes.Join([][]byte{[]byte{idx}, innerEncode}, nil), nil
}

func (i *{{ .Name }}) UnmarshalSCALE(reader io.Reader) error {
	{{ .UnmarshalSCALE }}
}`

var GenericTupleStructTemplate = `type {{ .GenericTupleName }}[{{ .GenericArity }}] struct {
{{ .GenericTupleFields }}
}

func (t *{{ .GenericTupleName }}[{{ .GenericNames }}]) MarshalSCALE() (output []byte, err error) {
	output = make([]byte, 0)
	var enc []byte
	{{ range $i, $a := .Fields }}
	enc, err = t.{{ $a }}.MarshalSCALE()
	if err != nil {
		return nil, err
	}
	output = append(output, enc...)
	{{ end }}
	return output, nil
}

func (t *{{ .GenericTupleName }}[{{ .GenericNames }}]) UnmarshalSCALE(reader io.Reader, {{ .UnmarshalFuncSignatures }}) (err error) {
	{{ range $func, $field := .FuncsAndFields }}
	t.{{ $field }}, err =  {{ $func }}(reader)
	if err != nil {
		return err
	}
	{{ end }}
	return nil
}

func Unmarshal{{.GenericTupleName}}FromRawBytes[{{ .GenericArity }}](
	{{ .UnmarshalFuncSignatures }}) func(io.Reader) (*{{ .GenericTupleName }}[{{ .GenericNames }}], error) {
	return func(reader io.Reader) (*{{ .GenericTupleName }}[{{ .GenericNames }}], error) {
		tuple := new({{.GenericTupleName}}[{{ .GenericNames }}])
		err := tuple.UnmarshalSCALE(reader,{{ range $func, $field := .FuncsAndFields }}
			{{ $func }},{{ end }})
		
		if err != nil {
			return nil, err
		}
		return tuple, nil
	}
}`
