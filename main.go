package main

import (
	res "boundary/resources"

	"boundary/utils"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// func main() {
// 	//pulumi.Run(DeployBoundary)
// }

func main() {

	pulumi.Run(func(ctx *pulumi.Context) error {

		// Verify required config params before starting deployment
		//CheckConfig(ctx)

		// Deploy Boundary Infra
		res.DeployBoundary(ctx)

		return nil
	})
}

func CheckConfig(ctx *pulumi.Context) *BoundaryInfraConfig {
	return &BoundaryInfraConfig{
		CloudFlare: &CloudFlareConfig{
			Email:  utils.GetConfString(ctx, "cloudflare.email"),
			ApiKey: utils.GetConfString(ctx, "cloudflare.api.key"),
		},
	}

}

type BoundaryInfraConfig struct {
	CloudFlare *CloudFlareConfig
	Dns        *PrivateDnsConfig
}

type CloudFlareConfig struct {
	Email      string
	ApiKey     string
	AccountId  string
	ZoneId     string
	Dnszone    string
	TunnelName string
}

type PrivateDnsConfig struct {
	WrokerArecord     string
	ControllerARecord string
	PrivateDnsZone    string
	ResourceGroup     string
}

type AzureAdConfig struct {
	AppDisplayName string
}
