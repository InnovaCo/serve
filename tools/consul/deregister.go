package consul

import (
	"fmt"
	"log"
	"regexp"

	"github.com/codegangsta/cli"
	"github.com/hashicorp/consul/api"
)

func DeregisterCommand() cli.Command {
	return cli.Command{
		Name:  "deregister",
		Usage: "Deregister all unhealth services on cluster",
		Subcommands: []cli.Command{
			{
				Name: "service",
				Flags: []cli.Flag{
					cli.BoolFlag{Name: "dry-run"},
				},
				Action: func(c *cli.Context) error {
					consul, _ := api.NewClient(api.DefaultConfig())

					name := c.Args().First()
					if name == "" {
						return fmt.Errorf("Service name is required! Given '%s'", name)
					}
					nameRegex := regexp.MustCompile(name)

					servs, err := consul.Agent().Services()
					if err != nil {
						return err
					}

					for _, srv := range servs {
						if nameRegex.MatchString(srv.ID) {
							log.Println("Deregistering", srv.ID)

							if !c.BoolT("dry-run") {
								err := consul.Agent().ServiceDeregister(srv.ID)
								if err != nil {
									return err
								}
							}
						}
					}

					return nil
				},
			},
			{
				Name:  "unhealth",
				Usage: "Deregister all unhealth services on cluster",
				Flags: []cli.Flag{
					cli.BoolFlag{Name: "dry-run"},
				},
				Action: func(c *cli.Context) error {
					consul, _ := api.NewClient(api.DefaultConfig())

					checks, err := consul.Agent().Checks()
					if err != nil {
						return err
					}

					for _, check := range checks {
						if check.Status == "critical" {
							log.Println("Deregistering", check.ServiceID, "=", check.Status)

							if !c.BoolT("dry-run") {
								err := consul.Agent().ServiceDeregister(check.ServiceID)
								if err != nil {
									return err
								}
							}
						}
					}

					return nil
				},
			},
		},
	}
}
