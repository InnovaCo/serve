package processor

import (
	"github.com/Jeffail/gabs"
)

type AnchorMerger struct{}

func (t AnchorMerger) Process(tree *gabs.Container) (*gabs.Container, error) {
	return ProcessAll(tree, func (ktype string, output *gabs.Container, value interface{}, key interface{}) error {
		if ktype == "map" && key == "<<" {
			if mmap, ok := value.(map[string]interface{}); ok {
				for k, v := range mmap {
					if !output.Exists(k) {
						output.Set(v, k)
					}
				}

				return output.Delete("<<")
			}
		}

		return nil
	})
}