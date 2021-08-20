package flufikinfo

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"regexp"
)

type FlufikDependency struct {
	Name     string
	Version  string
	Operator string
}

func (flufikDep *FlufikDependency) FlufikUnmarshalYaml(v *yaml.Node) error {
	var versionString string

	if err := v.Decode(&versionString); err != nil {
		return fmt.Errorf("unexpected version type string: %w", err)
	}

	reg, err := regexp.Compile("[<=>]+")
	if err != nil {
		return fmt.Errorf("unexpected regex format: %w", err)
	}

	pos := reg.FindStringIndex(versionString)
	if pos == nil {
		flufikDep.Name = versionString
		flufikDep.Operator = ""
		flufikDep.Version = ""
	} else {
		flufikDep.Name = versionString[:pos[0]]
		flufikDep.Operator = versionString[pos[0]:pos[1]]
		flufikDep.Version = versionString[pos[1]:]
	}
	return nil
}

func (flufikDep *FlufikDependency) FlufikRPMFormat() string {
	return fmt.Sprintf("%s%s%s", flufikDep.Name, flufikDep.Operator, flufikDep.Version)
}

func (flufikDep *FlufikDependency) FlufikDEBFormat() string {
	flufikOperator := flufikDep.Operator
	if flufikDep.Operator == "<" || flufikDep.Operator == ">" {
		flufikOperator += flufikDep.Operator
	}
	return fmt.Sprintf("%s (%s %s)", flufikDep.Name, flufikDep.Operator, flufikDep.Version)
}
