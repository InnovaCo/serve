package plugins

import (
	"github.com/InnovaCo/serve/manifest"
	"github.com/InnovaCo/serve/utils"
)

func init() {
	manifest.PluginRegestry.Add("outdated", Outdated{})
}

type Outdated struct{}

func (p Outdated) Run(data manifest.Manifest) error {
	consul, err := utils.ConsulClient(data.GetString("consul-address"))
	if err != nil {
		return err
	}

	fullName := data.GetString("full-name")

	if err := utils.MarkAsOutdated(consul, fullName, 0); err != nil {
		return err
	}

	if err := utils.DelConsulKv(consul, "services/routes/"+fullName); err != nil {
		return err
	}

	return nil
}
