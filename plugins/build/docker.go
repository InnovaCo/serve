package build

import (
	"github.com/InnovaCo/serve/manifest"
	"github.com/InnovaCo/serve/utils"
)

func init() {
	manifest.PluginRegestry.Add("build.docker", BuildDocker{})
}

type BuildDocker struct{}

func (p BuildDocker) Run(data manifest.Manifest) error {
	if err := utils.RunCmd("docker build --pull -t %s .", data.GetString("image")); err != nil {
		return err
	}

	return utils.RunCmd("docker push %s", data.GetString("image"))
}
