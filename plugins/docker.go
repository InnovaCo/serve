package plugins

import (
  "fmt"
  "log"
  "github.com/fatih/color"
  "github.com/InnovaCo/serve/manifest"
  "github.com/InnovaCo/serve/utils"
)

func init() {
  manifest.PluginRegestry.Add("docker", Docker{})
}

type Docker struct{}

func (p Docker) Run(data manifest.Manifest) error {

  operation := data.GetString("operation")
  docker_host := data.GetString(operation + ".host")
  ssh_user := data.GetString(operation + ".sshuser")
  image := data.GetString(operation + ".image")
  name := fmt.Sprintf(" --name=%s-%s", operation, data.GetString(operation + ".name"))
  cmd := "docker run "
  volumes := ""
  envs := ""
  parameters := ""
  ports := ""

  if data.Has(operation + ".volumes") {
    for _, vol := range data.GetArray(operation + ".volumes") {
      volumes += fmt.Sprintf(" --volume=%s:%s:%s", vol.GetString("hostPath"), vol.GetString("containerPath"), vol.GetString("mode"))
    }
  }
 
  if data.Has(operation + ".envs") {
    for key, value := range data.GetMap(operation + ".envs") {
      envs += fmt.Sprintf(" --env='%s=%s'", key, value)
    }
  }
 
  if data.Has(operation + ".ports") {
    for _, publish := range data.GetArray(operation + ".ports") {
      ports += fmt.Sprintf(" --publish=%s:%s", publish.GetString("hostPort"), publish.GetString("containerPort"))
    }
  }

  if data.Has(operation + ".parameters") {
    for _, item := range data.GetArray(operation + ".parameters") {
      parameters += fmt.Sprintf(" %s", item.Unwrap().(string))
    }
  }

  cmd = fmt.Sprintf("%s %s %s %s %s %s %s", cmd, volumes, envs, parameters, ports, name, image)
  log.Println(color.GreenString("<<<<< : (%s) : >>>>>", cmd))

  return utils.RunSingleSshCmd(docker_host, ssh_user, cmd)

}

func (p Docker) Stop(data manifest.Manifest) error {
	return nil
}
