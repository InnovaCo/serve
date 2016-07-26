package consul

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"github.com/hashicorp/consul/api"

	"github.com/InnovaCo/serve/utils"
)

var upstreamNameRegex = regexp.MustCompile("[^\\w]+")

func NginxTemplateContextCommand() cli.Command {
	return cli.Command{
		Name:  "nginx-template-context",
		Usage: "Collect and return data for consul-template",
		Action: func(c *cli.Context) error {
			consul, _ := api.NewClient(api.DefaultConfig())

			upstreams := make(map[string]map[string]map[string]interface{})
			services := make(map[string]map[string][]map[string]string)
			staticServers := make(map[string]map[string]string)
			duplicates := make(map[string]string)

			allServicesRoutes, _, err := consul.KV().List("services/routes/", nil)
			if err != nil {
				return fmt.Errorf("Error on load routes: %s", err)
			}

			for _, kv := range allServicesRoutes {
				name := strings.TrimPrefix(kv.Key, "services/routes/")

				instances, _, err := consul.Health().Service(name, "", true, nil)
				if err != nil {
					return fmt.Errorf("Error on get service `%s` health: %s", name, err)
				}

				if len(instances) == 0 {
					break
				}

				routes := make([]map[string]string, 0)
				if err := json.Unmarshal(kv.Value, &routes); err != nil {
					return fmt.Errorf("Error on parse route json: %s. Serive `%s`, json: %s", err, name, string(kv.Value))
				}

				upstream := upstreamNameRegex.ReplaceAllString("serve_"+name, "_")
				if instances[0].Service.Port == 0 {
					upstream += "_static"
				}

				for _, route := range routes {
					host, ok := route["host"]
					if !ok {
						return fmt.Errorf("Host is required for routing! Service `%s`", name)
					}

					location, ok := route["location"]
					if !ok {
						location = "/"
					}

					parts := strings.Split(name, "/")
					category := strings.Join(parts[0:len(parts)-1], "/")
					packageName := parts[len(parts)-1]

					routeUpstream := upstream
					if ups, ok := route["upstream"]; ok && ups == "static" && !strings.HasSuffix(upstream, "_static") {
						routeUpstream = upstream + "_static"
					}

					staticHost := ""
					if strings.HasSuffix(routeUpstream, "_static") {
						staticHost = upstreamNameRegex.ReplaceAllString(category, "-") + "." + upstreamNameRegex.ReplaceAllString(packageName, "-") + ".static"

						staticServers[staticHost] = map[string]string{
							"category": category,
							"package":  packageName,
						}
					}

					for _, inst := range instances {
						putUpstream(routeUpstream, inst, upstreams)
					}

					delete(route, "host")
					delete(route, "location")
					delete(route, "upstream")

					if _, ok := services[host]; !ok {
						services[host] = make(map[string][]map[string]string, 0)
					}

					if _, ok := services[host][location]; !ok {
						services[host][location] = make([]map[string]string, 0)
					}

					routeKeys := "-"
					routeValues := "-"
					for k, v := range route {
						routeKeys += "${" + k + "}-"
						routeValues += v + "-"
					}

					if exists, ok := duplicates[host+location+routeKeys+routeValues]; !ok {
						duplicates[host+location+routeKeys+routeValues] = name
					} else {
						fmt.Fprintln(os.Stderr, color.RedString("Service with the same routes already exists! exists: %s, skipped: %s", exists, name))
						continue
					}

					services[host][location] = append(services[host][location], map[string]string{
						"upstream":    routeUpstream,
						"staticHost":  staticHost,
						"routeKeys":   routeKeys,
						"routeValues": routeValues,
						"sortIndex":   strconv.Itoa(len(route)),
					})
				}
			}

			// sort routes by sort index
			for _, hh := range services {
				for _, ll := range hh {
					sort.Sort(utils.BySortIndex(ll))
				}
			}

			out, _ := json.MarshalIndent(map[string]interface{}{
				"upstreams":     upstreams,
				"services":      services,
				"staticServers": staticServers,
			}, "", "  ")

			fmt.Fprintln(os.Stdout, string(out))
			return nil
		},
	}
}

func putUpstream(upstream string, inst *api.ServiceEntry, upstreams map[string]map[string]map[string]interface{}) {
	port := inst.Service.Port
	if port == 0 || strings.HasSuffix(upstream, "_static") {
		port = 83
	}

	if _, ok := upstreams[upstream]; !ok {
		upstreams[upstream] = make(map[string]map[string]interface{}, 0)
	}

	address := inst.Node.Address
	if inst.Service.Address != "" {
		address = inst.Service.Address
	}

	upstreams[upstream][fmt.Sprintf("%s:%d", address, port)] = map[string]interface{}{
		"address": address,
		"port":    port,
	}
}
