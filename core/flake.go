package core

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var PluginTemplate = `
var run = map[string]func(...string) (string, error){
	{{range .}}{{.Key}}: plugin.Run{{.Value}},{{end}}
}
`

type Ptemplate struct {
	Key   string
	Value string
}

type Ptemplates []Ptemplate

type Flake struct {
	plugin     bool
	ptemplates Ptemplates
	command    *cobra.Command
}

func FlakeNew() *Flake {
	flake := &Flake{
		command: &cobra.Command{
			Use:   "flake",
			Short: "run with cli mode",
			Long:  "run with cli mode",
		},
	}

	flake.command.RunE = flake.FlakeCmdRun

	// Register flags on flake sub command
	flake.command.Flags().BoolVarP(&flake.plugin, "plugin-register", "", false, "register plugin")

	return flake
}

func (flake *Flake) FlakeCmdRun(cmd *cobra.Command, args []string) error {
	if flake.plugin {
		return flake.pluginRegister()
	}
	return nil
}

func (flake *Flake) pluginRegister() error {

	fileInfos, err := ioutil.ReadDir("plugin")
	if err != nil {
		return errors.Wrapf(err, "Could not open plugin directory")
	}

	flake.findSourceCode(fileInfos)
	tmpl := template.New("PluginTemplate for executor.go")
	template.Must(tmpl.Parse(PluginTemplate))
	tmpl.Execute(os.Stdout, flake.ptemplates)

	return nil
}

func (flake *Flake) findSourceCode(fileInfos []os.FileInfo) {
	regex := regexp.MustCompile(`func Run\([A-Za-z]+\)`)

	for _, fi := range fileInfos {
		f, err := os.Open("plugin/" + fi.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}
		flake.grepRunInPlugins(regex, f)
		f.Close()
	}
}

func (flake *Flake) grepRunInPlugins(re *regexp.Regexp, f *os.File) {
	// ここに open して Run[A-Za-z]+()を探して,
	// map へ "[a-z]+":"Run[A-Za-z]+" を保存する
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		os.Stdout.Write([]byte(text))
		if re.MatchString(text) {
			captured := re.FindStringSubmatch(text)
			name := captured[1]
			flake.ptemplates = append(flake.ptemplates, Ptemplate{
				Key:   strings.ToLower(name),
				Value: name,
			})

			return
		}
	}
}
