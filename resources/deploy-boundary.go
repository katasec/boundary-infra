package resources

import (
	"boundary/utils"
	"fmt"

	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func DeployBoundary(ctx *pulumi.Context) error {

	// Init locals
	location := "EastAsia"
	rgName := "rg-ea-boundary"

	/*
		Boundary Worker and Contoller need TLS for communication.
		We'll generate self-sign certs for them to use to comminicate
		with each other
	*/

	// Create private key for cert generation
	CreatePrivateKey(ctx)

	// Create controller cert with private key
	controllerFQDN := fmt.Sprintf("%s.%s", ControllerARecord, DnsZone)
	Controller_cert = CreateSelfSignedCert(ctx, "boundary_controller_cert", controllerFQDN, "Acme ltd.")

	// Create worker cert with private key
	workerFQDN := fmt.Sprintf("%s.%s", WorkerARecord, DnsZone)
	Worker_cert = CreateSelfSignedCert(ctx, "boundary_worker_cert", workerFQDN, "Acme ltd.")

	// Create an Azure SP for boundary worker & controller to access vault.
	CreateAzureServicePrincipal(ctx)

	// Create a resource group for boundary resources such Azure Key Vault, PostreSQL etc.
	CreateResourceGroup(ctx, rgName, location)

	// Create Keyvault for Boundary with appropriate access policy
	kvRequest := &KeyVaultRequest{
		Name:                      "kv-boundary-p1",
		Location:                  location,
		ResourceGroupName:         ResourceGroup.Name,
		EnabledForDiskEncryption:  true,
		SoftDeleteRetentionInDays: 7,
		EnabledForDeployment:      true,
	}
	vault, _ := CreateNewKeyVault(ctx, *kvRequest)

	// Create boundary encryption keys for root, worker and recovery
	RootKey, _ = CreateKeyVaultKey(ctx, vault, "root")
	WorkerKey, _ = CreateKeyVaultKey(ctx, vault, "worker")
	RecoveryKey, _ = CreateKeyVaultKey(ctx, vault, "recovery")

	// Create PostgreSQL DB for Boundary
	PostgresAdministratorLogin = pulumi.String("boundary")

	// Create a random password for the runnel
	pgsqlPassword, err := random.NewRandomString(ctx, "pgsqlPassword", &random.RandomStringArgs{
		Length:  pulumi.Int(44),
		Special: pulumi.Bool(false),
	})
	utils.ExitOnError(err)

	CreateKeyVaultSecret(ctx, vault, "pg-admin-password", pgsqlPassword.Result)
	CreateKeyVaultSecret(ctx, vault, "ngrok-auth-token", pulumi.String(Ngrok_Auth_Token).ToStringOutput())

	CreatePostgreSQL(ctx, rgName, location, "bdry-pulumi", PostgresAdministratorLogin, pgsqlPassword.Result)

	// Deploy Boundary in an AKS cluster
	DeployToK8s(ctx)

	return nil
}
