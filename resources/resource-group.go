package resources

import (
	"log"

	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateResourceGroup(ctx *pulumi.Context, rgName string, location string) (err error) {

	ResourceGroup, err = resources.NewResourceGroup(ctx, rgName, &resources.ResourceGroupArgs{
		Location: pulumi.String(location),
	})

	if err != nil {
		log.Fatalf(err.Error())
	}

	return nil
}
