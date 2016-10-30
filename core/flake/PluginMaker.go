package flake

import (
	"html/template"
	"os"
	"strings"

	"github.com/pkg/errors"
)

var PluginMakerTemplate = `package plugin

// Detail{{.Name}} is description {{.LowerName}} command
var Detail{{.Name}} = ""

// Run{{.Name}} is root function of {{.LowerName}} command
func Run{{.Name}}(args ...string) (string, error) {
    return "", nil
}

`

type PMtemplate struct {
	Name      string
	LowerName string
}

func PluginMaker(name string) error {

	// Check to exist plugin directory on the current directory
	if _, err := os.Stat("plugin"); err != nil {
		return errors.Wrapf(err, "Does not exist plugin directory")
	}

	// Check the "...".go on the plugin directory
	if _, err := os.Stat("plugin/" + name + ".go"); err == nil {
		return errors.Errorf("plugin/%s.go is already exist", name)
	}

	tmpl := template.New("PluginMakerTemplate for executor.go")
	pmtmpl := PMtemplate{Name: strings.Title(name), LowerName: name}
	template.Must(tmpl.Parse(PluginMakerTemplate))

	f, err := os.Create("plugin/" + name + ".go")
	if err != nil {
		return errors.Wrapf(err, "Failed to create %s.go", name)
	}
	defer f.Close()
	tmpl.Execute(f, pmtmpl)

	os.Stdout.WriteString("Created plugin/" + name + ".go\n")

	return nil
}
