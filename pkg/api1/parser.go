package api1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const CommentSign = "#"
const SemSign = "@"

const tID = "[A-Za-z][0-9A-Za-z_]*?"
const tType = "[A-Za-z\\[][0-9A-Za-z_\\[\\]\\?]*?"

var reComment = compile("^([^%s]*?)\\s*(%s.*)$", CommentSign, CommentSign)
var reSemComment = compile("^\\s*%s([0-9A-Za-z_.]+?)(?::(json|ya?ml))?(?:|\\s*([|])|\\s+([^|].*?))\\s*$", SemSign)
var reGroup = compile("^group\\s+(%s)$", tID)
var reScalar = compile("^scalar\\s+(%s)$", tID)
var reBlockStart = compile("^(enum|struct|interface)\\s+(%s)\\s*\\{$", tID)
var reEnumOption = compile("^(%s)$", tID)
var reStructField = compile("^(%s)\\s*:\\s*(%s)$", tID, tType)
var reFunInline = compile("^(%s)\\s*\\((.*?)\\)(?:\\s*:\\s*(%s))?$", tID, tType)
var reFunStart = compile("^(%s)\\s*\\(\\s*(.*?)$", tID)
var reFunEnd = compile("^(.*?)\\s*\\)(?:\\s*:\\s*(%s))?$", tType)
var reParam = compile("^(%s)\\s*:\\s*(%s)(?:\\s*=\\s*(.+?))?$", tID, tType)
var reCommaSep = compile("\\s*,\\s*")

type Parser struct {
	comments     []string
	postComments []string
}

func (p *Parser) Parse(contents ...string) (*Schema, error) {
	schema := &Schema{}
	for _, content := range contents {
		subSchema, err := p.parse(content)
		if err != nil {
			return nil, err
		}
		schema.Groups = append(schema.Groups, subSchema.Groups...)
	}
	schema.MergeGroupIfaces()
	if err := schema.Check(); err != nil {
		return nil, err
	}
	return schema, nil
}

func (p *Parser) ParseFiles(files ...string) (*Schema, error) {
	schema := &Schema{}
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, errors.Errorf("Read file [%s] failed: %v", file, err)
		}
		subSchema, err := p.parse(string(content))
		if err != nil {
			return nil, errors.Errorf("Parse file [%s] failed: %v", file, err)
		}
		schema.Groups = append(schema.Groups, subSchema.Groups...)
	}
	schema.MergeGroupIfaces()
	if err := schema.Check(); err != nil {
		return nil, err
	}
	return schema, nil
}

func (p *Parser) parse(content string) (*Schema, error) {
	var err error
	var group *ApiGroup = nil
	var parsingBlock bool = false
	var enumType *EnumType = nil
	var structType *StructType = nil
	var iface *Iface = nil
	var fun *Fun = nil

	p.flushComments()

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if m := reComment.FindStringSubmatch(line); m != nil {
			line = m[1]
			comment := m[2]
			if line == "" {
				p.addComment(comment)
			} else {
				p.addPostComment(comment)
			}
		}
		if line == "" {
			continue
		}

		if group == nil {
			if m := reGroup.FindStringSubmatch(line); m != nil {
				group = &ApiGroup{}
				group.Name = m[1]
				group.HasComments, err = p.flushComments()
				if err != nil {
					return nil, err
				}
				continue
			}
			return nil, errors.Errorf("first line must be group, invalid line: [%s]", line)
		}

		if !parsingBlock {
			if m := reScalar.FindStringSubmatch(line); m != nil {
				scalarType := &ScalarType{}
				scalarType.Name = m[1]
				scalarType.HasComments, err = p.flushComments()
				if err != nil {
					return nil, err
				}
				group.ScalarTypes = append(group.ScalarTypes, *scalarType)
				continue
			}
			if m := reBlockStart.FindStringSubmatch(line); m != nil {
				switch m[1] {
				case "enum":
					enumType = &EnumType{}
					enumType.Name = m[2]
					enumType.HasComments, err = p.flushComments()
					if err != nil {
						return nil, err
					}
				case "struct":
					structType = &StructType{}
					structType.Name = m[2]
					structType.HasComments, err = p.flushComments()
					if err != nil {
						return nil, err
					}
				case "interface":
					iface = &Iface{}
					iface.Name = m[2]
					iface.HasComments, err = p.flushComments()
					if err != nil {
						return nil, err
					}
				}
				parsingBlock = true
				continue
			}
			return nil, errors.Errorf("invalid line [%s] for parsing.", line)
		}

		if line == "}" {
			if enumType != nil {
				p.flushPostCommentsTo(&enumType.HasComments)
				group.EnumTypes = append(group.EnumTypes, *enumType)
				enumType = nil
			} else if structType != nil {
				p.flushPostCommentsTo(&structType.HasComments)
				group.StructTypes = append(group.StructTypes, *structType)
				structType = nil
			} else {
				p.flushPostCommentsTo(&iface.HasComments)
				group.Ifaces = append(group.Ifaces, *iface)
				iface = nil
			}
			parsingBlock = false
			continue
		}
		if enumType != nil {
			option, err := p.parseEnumOption(line)
			if err != nil {
				return nil, err
			}
			enumType.Options = append(enumType.Options, *option)
			continue
		}
		if structType != nil {
			field, err := p.parseStructField(line)
			if err != nil {
				return nil, err
			}
			structType.Fields = append(structType.Fields, *field)
			continue
		}

		// assert: iface != nil
		if fun == nil {
			if m := reFunInline.FindStringSubmatch(line); m != nil {
				fun := &Fun{}
				fun.Name = m[1]
				fun.HasComments, err = p.flushComments()
				if err != nil {
					return nil, err
				}
				params, err := p.parseParams(m[2])
				if err != nil {
					return nil, err
				}
				fun.Params = params
				fun.Type = p.parseType(m[3])
				iface.Funs = append(iface.Funs, *fun)
				continue
			}
			if m := reFunStart.FindStringSubmatch(line); m != nil {
				c, err := p.flushComments()
				if err != nil {
					return nil, err
				}
				fun = &Fun{}
				fun.Name = m[1]
				params, err := p.parseParams(m[2])
				if err != nil {
					return nil, err
				}
				if len(params) > 0 {
					fun.Params = params
					fun.Comments = c.Comments
					fun.SemComments = c.SemComments
					params[len(params)-1].PostComments = c.PostComments
				} else {
					fun.HasComments = c
				}
				continue
			}
			return nil, errors.Errorf("invalid line [%s] for parsing.", line)
		}
		if m := reFunEnd.FindStringSubmatch(line); m != nil {
			params, err := p.parseParams(m[1])
			if err != nil {
				return nil, err
			}
			if len(params) > 0 {
				postComments := params[len(params)-1].PostComments
				params[len(params)-1].PostComments = nil
				fun.PostComments = append(fun.PostComments, postComments...)
				fun.Params = append(fun.Params, params...)
			} else {
				p.flushPostCommentsTo(&fun.HasComments)
			}
			fun.Type = p.parseType(m[2])
			iface.Funs = append(iface.Funs, *fun)
			fun = nil
			continue
		}
		params, err := p.parseParams(line)
		if err != nil {
			return nil, err
		}
		fun.Params = append(fun.Params, params...)
	}

	if parsingBlock {
		return nil, errors.New("object in parsing is not finished yet, need close brace.")
	}
	p.flushPostCommentsTo(&group.HasComments)
	return &Schema{Groups: []ApiGroup{*group}}, nil
}

func (p *Parser) parseEnumOption(line string) (*EnumOption, error) {
	var err error
	if m := reEnumOption.FindStringSubmatch(line); m != nil {
		option := &EnumOption{}
		option.Name = m[1]
		option.HasComments, err = p.flushComments()
		if err != nil {
			return nil, err
		}
		return option, nil
	}
	return nil, fmt.Errorf("invalid line [%s] for enum option", line)
}

func (p *Parser) parseStructField(line string) (*StructField, error) {
	var err error
	if m := reStructField.FindStringSubmatch(line); m != nil {
		field := &StructField{}
		field.Name = m[1]
		field.HasComments, err = p.flushComments()
		if err != nil {
			return nil, err
		}
		field.Type = p.parseType(m[2])
		return field, nil
	}
	return nil, fmt.Errorf("invalid line [%s] for struct field", line)
}

func (p *Parser) parseParams(s string) ([]Param, error) {
	var params []Param
	if s == "" {
		return params, nil
	}
	parts := reCommaSep.Split(s, -1)
	var m []string
	for _, part := range parts {
		if part == "" {
			continue
		}
		if m = reParam.FindStringSubmatch(part); m == nil {
			return nil, fmt.Errorf("invalid string [%s] for params", s)
		}
		param := Param{}
		param.Name = m[1]
		param.Type = p.parseType(m[2])
		params = append(params, param)
	}
	if len(params) > 0 {
		c, err := p.flushComments()
		if err != nil {
			return nil, err
		}
		params[0].Comments = c.Comments
		params[0].SemComments = c.SemComments
		params[len(params)-1].PostComments = c.PostComments
	}
	return params, nil
}

func (p *Parser) parseType(s string) *TypeRef {
	if len(s) == 0 {
		return nil
	}
	nullable := false
	if len(s) > 1 && s[len(s)-1:] == "?" {
		s = s[:len(s)-1]
		nullable = true
	}
	if len(s) > 2 && s[:1] == "[" && s[len(s)-1:] == "]" {
		itemTypeStr := strings.TrimSpace(s[1 : len(s)-1])
		arrayType := TypeRef{}
		arrayType.ItemType = p.parseType(itemTypeStr)
		arrayType.Nullable = nullable
		return &arrayType
	}
	singleType := TypeRef{}
	singleType.Name = s
	singleType.Nullable = nullable
	return &singleType
}

func (p *Parser) addComment(s string) {
	if strings.HasPrefix(s, CommentSign+" ") {
		s = strings.TrimPrefix(s, CommentSign+" ")
	} else {
		s = strings.TrimPrefix(s, CommentSign)
	}
	p.comments = append(p.comments, s)
}

func (p *Parser) addPostComment(s string) {
	if strings.HasPrefix(s, CommentSign+" ") {
		s = strings.TrimPrefix(s, CommentSign+" ")
	} else {
		s = strings.TrimPrefix(s, CommentSign)
	}
	p.postComments = append(p.postComments, s)
}

func (p *Parser) flushComments() (HasComments, error) {
	c := HasComments{}
	var semKey string
	var semVal string
	var semFmt string
	var mlMode bool
	semCommit := func() error {
		if semKey == "" {
			return nil
		}
		var val interface{}
		var err error
		if semFmt == "json" {
			err = json.Unmarshal([]byte(semVal), &val)
		} else if semFmt == "yaml" || semFmt == "yml" {
			err = yaml.Unmarshal([]byte(semVal), &val)
		} else {
			val = semVal
		}
		if err != nil {
			return errors.Wrapf(err, "invalid %s string: [%s]", semFmt, semVal)
		}
		c.AddSemComment(semKey, val)
		semKey = ""
		semVal = ""
		semFmt = ""
		mlMode = false
		return nil
	}
	for _, line := range p.comments {
		if mlMode {
			var matched bool
			if strings.HasPrefix(line, "  ") {
				line = line[2:]
				matched = true
			} else if strings.HasPrefix(line, " ") {
				line = line[1:]
				matched = true
			}
			if matched {
				if len(semVal) > 0 {
					semVal += "\n"
				}
				semVal += line
				continue
			}
			if err := semCommit(); err != nil {
				return HasComments{}, err
			}
		}
		if m := reSemComment.FindStringSubmatch(line); m != nil {
			semKey = m[1]
			semFmt = m[2]
			mlMode = len(m[3]) > 0
			if mlMode {
				continue
			}
			semVal = m[4]
			if err := semCommit(); err != nil {
				return HasComments{}, err
			}
			continue
		}
		c.Comments = append(c.Comments, line)
	}
	if err := semCommit(); err != nil {
		return HasComments{}, err
	}
	c.PostComments = p.postComments
	p.comments = nil
	p.postComments = nil
	return c, nil
}

func (p *Parser) flushPostCommentsTo(c *HasComments) {
	c.PostComments = append(c.PostComments, p.comments...)
	c.PostComments = append(c.PostComments, p.postComments...)
	p.comments = nil
	p.postComments = nil
}

func compile(s string, v ...interface{}) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf(s, v...))
}
