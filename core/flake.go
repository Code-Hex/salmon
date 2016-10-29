package core

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var PluginTemplate = `
var run = map[string]func(...string) (string, error){
{{range .}}    "{{.Key}}": plugin.Run{{.Value}},
{{end}}
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
	if err := flake.MigrateExecutorGo(); err != nil {
		return err
	}

	return nil
}

func (flake *Flake) findSourceCode(fileInfos []os.FileInfo) {
	regex := regexp.MustCompile(`func Run([A-Za-z]+)`)

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
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
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

func (flake *Flake) MigrateExecutorGo() error {
	regex := regexp.MustCompile(`var run = map`)

	f, err := os.OpenFile("executor.go", os.O_RDWR, 0666)
	if err != nil {
		return errors.Wrapf(err, "Could not open executor.go")
	}
	defer f.Close()

	var src string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		if regex.MatchString(text) {
			break
		}
		src += text + "\n"
	}

	f.Seek(0, io.SeekStart)
	f.Truncate(0)

	tmpl := template.New("PluginTemplate for executor.go")
	template.Must(tmpl.Parse(src + PluginTemplate))

	tmpl.Execute(f, flake.ptemplates)

	out, err := exec.Command("goimports", "executor.go").Output()
	if err != nil {
		return errors.Wrapf(err, "Failed goimports")
	}

	f.Seek(0, io.SeekStart)
	f.Truncate(0)

	f.Write(out)

	return nil
}
