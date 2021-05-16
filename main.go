package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html"
	"html/template"
	"os"
	"sort"

	"gopkg.in/yaml.v2"
)

type resource struct {
	ID   string
	Type string `yaml:"Type,omitempty" json:"Type,omitempty"`
}

type output struct {
	ID          string
	Description string `yaml:"Description,omitempty" json:"Description,omitempty"`
}

type parameter struct {
	ID          string
	Description string `yaml:"Description,omitempty" json:"Description,omitempty"`
	Type        string `yaml:"Type,omitempty" json:"Type,omitempty"`
	Default     string `yaml:"Default,omitempty" json:"Default,omitempty"`
}

type cfnTemplate struct {
	Description string               `yaml:"Description,omitempty" json:"Description,omitempty"`
	Parameters  map[string]parameter `yaml:"Parameters,omitempty" json:"Parameters,omitempty"`
	Resources   map[string]resource  `yaml:"Resources,omitempty" json:"Resources,omitempty"`
	Outputs     map[string]output    `yaml:"Outputs,omitempty" json:"Outputs,omitempty"`
}

type templateValue struct {
	Filename    string
	Description string
	Parameters  []parameter
	Resources   []resource
	Outputs     []output
}

//go:embed doc.template
var docTemplate string

func newTemplateValue(filename string, raw cfnTemplate) templateValue {
	temp := templateValue{
		Filename:    filename,
		Description: raw.Description,
		Parameters:  []parameter{},
		Resources:   []resource{},
		Outputs:     []output{},
	}

	for id, param := range raw.Parameters {
		param.ID = id
		temp.Parameters = append(temp.Parameters, param)
	}
	sort.Slice(temp.Parameters, func(i, j int) bool {
		return temp.Parameters[i].ID < temp.Parameters[j].ID
	})

	for id, rs := range raw.Resources {
		rs.ID = id
		temp.Resources = append(temp.Resources, rs)
	}
	sort.Slice(temp.Resources, func(i, j int) bool {
		return temp.Resources[i].ID < temp.Resources[j].ID
	})

	for id, out := range raw.Outputs {
		out.ID = id
		temp.Outputs = append(temp.Outputs, out)
	}
	sort.Slice(temp.Outputs, func(i, j int) bool {
		return temp.Outputs[i].ID < temp.Outputs[j].ID
	})

	return temp
}

type argument struct {
	output   string
	filename string
}

func parseArgs(args []string) (argument, error) {
	arg := argument{}
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.StringVar(&arg.output, "o", "", "File name. Write to file instead of stdout")
	flags.Usage = func() {
		fmt.Fprintf(os.Stdout, "cfn-doc is a tool for generating a document of CloudFormation template.\n\n")
		fmt.Fprintf(os.Stdout, "Usage: \n")
		fmt.Fprintf(os.Stdout, "  cfn-doc [OPTION] <TEMPLATE FILE>\n\n")
		fmt.Fprintf(os.Stdout, "Options: \n")
		flags.PrintDefaults()
	}
	err := flags.Parse(args[1:])
	if err != nil {
		return argument{}, err
	}

	arg.filename = flags.Args()[0]
	return arg, nil
}

func generateDoc(filename string, cfn cfnTemplate) ([]byte, error) {
	tplValue := newTemplateValue(filename, cfn)
	tpl, err := template.New("").Parse(docTemplate)
	if err != nil {
		return nil, err
	}

	out := new(bytes.Buffer)
	err = tpl.Execute(out, tplValue)
	return out.Bytes(), err
}

func readCFnTemplate(filename string) (cfnTemplate, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return cfnTemplate{}, err
	}

	raw := cfnTemplate{}
	err = yaml.Unmarshal(b, &raw)
	if !errors.Is(err, new(yaml.TypeError)) {
		return raw, err
	}

	err = json.Unmarshal(b, &raw)
	return raw, err
}

func main() {
	arg, err := parseArgs(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cfn, err := readCFnTemplate(arg.filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	rawDoc, err := generateDoc(arg.filename, cfn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	doc := []byte(html.UnescapeString(string(rawDoc)))

	if arg.output == "" {
		fmt.Println(string(doc))
		return
	}

	if err := os.WriteFile(arg.output, doc, 0644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
