package all

import (
	"github.com/jinzhenj/api1/pkg/api1"
	"github.com/jinzhenj/api1/pkg/golang"
	"github.com/jinzhenj/api1/pkg/openapi"
	"github.com/jinzhenj/api1/pkg/utils"
)

type Render struct {
	parser        *api1.Parser
	openapiRender *openapi.Render
	golangRender  *golang.Render
}

func NewRender() *Render {
	return &Render{
		parser:        &api1.Parser{},
		openapiRender: &openapi.Render{},
		golangRender:  &golang.Render{},
	}
}

func (r *Render) RenderFiles(files []string) ([]CodeFile, error) {
	var codeFiles []CodeFile
	schema, err := r.parser.ParseFiles(files...)
	if err != nil {
		return nil, err
	}
	openAPI, err := r.openapiRender.Render(schema)
	if err != nil {
		return nil, err
	}
	goFiles, err := r.golangRender.Render(schema)
	if err != nil {
		return nil, err
	}
	codeFiles = append(codeFiles, CodeFile{
		Name:    "doc/api1.json",
		Content: utils.ToJson(schema) + "\n",
	})
	codeFiles = append(codeFiles, CodeFile{
		Name:    "doc/openapi.json",
		Content: utils.ToJson(openAPI) + "\n",
	})
	codeFiles = append(codeFiles, CodeFile{
		Name:    "doc/openapi.go",
		Content: renderOpenAPIGoFile(openAPI),
	})
	for _, goFile := range goFiles {
		codeFiles = append(codeFiles, CodeFile{
			Name:    goFile.Name,
			Content: goFile.Code(),
		})
	}
	return codeFiles, nil
}

func renderOpenAPIGoFile(openAPI *openapi.OpenAPI) string {
	openAPI.Servers = []openapi.Server{{
		Url: "{url}",
		Variables: map[string]openapi.ServerVariable{
			"url": {Default: "{{.BasePath}}"},
		},
	}}
	code := "package doc\n\n"
	code += "const OpenAPI = `" + utils.ToJson(openAPI) + "`\n"
	return code
}
