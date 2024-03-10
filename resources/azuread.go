package resources

import (
	"fmt"
	"log"

	"github.com/pulumi/pulumi-azuread/sdk/v5/go/azuread"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateAzureServicePrincipal(ctx *pulumi.Context) error {

	var err error
	// Create a new registered Azure AD app

	app, err := azuread.NewApplication(ctx, azureadAppDisplayName, &azuread.ApplicationArgs{
		DisplayName: pulumi.String(azureadAppDisplayName),

		// GroupMembershipClaims: pulumi.StringArray{
		// 	pulumi.String("All"),
		// },
		// RequiredResourceAccesses: azuread.ApplicationRequiredResourceAccessArray{
		// 	&azuread.ApplicationRequiredResourceAccessArgs{

		// 		ResourceAppId: pulumi.String("00000003-0000-0000-c000-000000000000"), // Microsoft Graph
		// 		ResourceAccesses: azuread.ApplicationRequiredResourceAccessResourceAccessArray{

		// 			&azuread.ApplicationRequiredResourceAccessResourceAccessArgs{
		// 				Id:   pulumi.String("df021288-bdef-4463-88db-98f22de89214"), // User.Read.All
		// 				Type: pulumi.String("Role"),
		// 			},

		// 			&azuread.ApplicationRequiredResourceAccessResourceAccessArgs{
		// 				Id:   pulumi.String("b4e74841-8e56-480b-be8b-910348b18b4c"), // User.ReadWrite
		// 				Type: pulumi.String("Scope"),
		// 			},
		// 		},
		// 	},
		// },
		// Web: &azuread.ApplicationWebArgs{
		// 	RedirectUris: pulumi.StringArray{
		// 		pulumi.String(fmt.Sprintf("%s/v1/auth-methods/oidc:authenticate:callback", BoundaryUrl)),
		// 	},
		// 	ImplicitGrant: &azuread.ApplicationWebImplicitGrantArgs{
		// 		AccessTokenIssuanceEnabled: pulumi.Bool(true),
		// 		IdTokenIssuanceEnabled:     pulumi.Bool(true),
		// 	},
		// },
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx.Export("boundaryApplicationId", app.ApplicationId)

	// Create a service principal in the App
	servicePrincipal, err = azuread.NewServicePrincipal(ctx, "boundary_service_principal", &azuread.ServicePrincipalArgs{
		ApplicationId: app.ApplicationId,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx.Export("boundaryServicePrincipal", servicePrincipal.ID())

	// Create password for the service principal
	servicePrincipalPassword, err = azuread.NewServicePrincipalPassword(ctx, "boundary_service_principal_password", &azuread.ServicePrincipalPasswordArgs{
		ServicePrincipalId: servicePrincipal.ID(),
		EndDateRelative:    pulumi.StringPtr("8760h"),
		DisplayName:        pulumi.StringPtr(fmt.Sprintf("%v-password", azureadAppDisplayName)),
	})

	ctx.Export("boundaryServicePrincipalPassword", servicePrincipalPassword.Value)
	if err != nil {
		log.Fatal(err.Error())
	}

	return nil
}
