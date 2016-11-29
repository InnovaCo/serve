package upload

import (
	"github.com/InnovaCo/serve/manifest"
	"github.com/InnovaCo/serve/utils"
)

func init() {
	manifest.PluginRegestry.Add("upload.tarball", UploadTarball{})
}

type UploadTarball struct{}

func (p UploadTarball) Run(data manifest.Manifest) error {
	if err := utils.RunCmd("curl -vsSf -o tarball.tar.gz %s", data.GetString("unstable-url")); err != nil {
		return err
	}

	return utils.RunCmd("curl -vsSf -XPUT -T tarball.tar.gz %s", data.GetString("stable-url"))
}
