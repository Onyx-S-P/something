package model

import (
	"gopkg.in/yaml.v2"
)

/*Struct for policy file*/
type Policy struct {
	APIVersion string            `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string            `json:"kind,omitempty" yaml:"kind,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Spec       Spec              `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type Spec struct {
	Tags     []string               `json:"tags,omitempty" yaml:"tags,omitempty"`
	Message  string                 `json:"message,omitempty" yaml:"message,omitempty"`
	Selector map[string]interface{} `json:"selector,omitempty" yaml:"selector,omitempty"`
	File     Filespec               `json:"file,omitempty" yaml:"file,omitempty"`
}

type Filespec struct {
	Severity int `json:"severity,omitempty" yaml:"severity,omitempty"`
	//MatchPatterns    map[string]interface{} `json:"matchPatterns,omitempty" yaml:"matchPatterns,omitempty"`
	//MatchPaths       map[string]interface{} `json:"matchPaths,omitempty" yaml:"matchPaths,omitempty"`
	//MatchDirectories map[string]interface{} `json:"matchDirectories,omitempty" yaml:"matchDirectories,omitempty"`
	Action string `json:"action,omitempty" yaml:"action,omitempty"`
}

/*get Policy Spec in String format*/
/*func (p *Policy) GetSpec() string {
	data, _ := yaml.Marshal(p.Spec)

	return string(data)
}*/
func (p *Policy) GetSpec() []byte {
	data, _ := yaml.Marshal(p.Spec)
	return data
}

func (p *Spec) GetFile() string {
	data, _ := yaml.Marshal(p.File)
	return string(data)
}

/*get Policy Metadata in String format*/
func (p *Policy) GetMetadata() string {
	data, _ := yaml.Marshal(p.Metadata)
	return string(data)
}

func (p *Spec) GetSelector() string {
	data, _ := yaml.Marshal(p.Selector)
	return string(data)
}

/*func (p *Filespec) GetMatchPaths() string {
	data, _ := yaml.Marshal(p.MatchPaths)
	return string(data)
}

func (p *Filespec) GetMatchDirectories() string {
	data, _ := yaml.Marshal(p.MatchDirectories)
	return string(data)
}*/

/*Creates new Policy struct from given data*/
func NewPolicyFrom(data []byte) Policy {
	p := Policy{}
	yaml.Unmarshal(data, &p)
	return p
}

func NewSpecFrom(data []byte) Spec {
	s := Spec{}
	yaml.Unmarshal(data, &s)
	return s
}

func NewFilespec(data []byte) Filespec {
	f := Filespec{}
	yaml.Unmarshal(data, &f)
	return f
}
