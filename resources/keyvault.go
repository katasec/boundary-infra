package resources

import (
	// github.com/pulumi/pulumi-azure-native/sdk/go/azure

	"log"
	"os"

	utils "boundary/utils"

	keyvault "github.com/pulumi/pulumi-azure-native/sdk/go/azure/keyvault"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type KeyVaultRequest struct {
	Name                      string
	Location                  string
	ResourceGroupName         pulumi.StringOutput
	TenantId                  string
	EnabledForDiskEncryption  bool
	SoftDeleteRetentionInDays int
	EnablePurgeProtection     bool
	EnabledForDeployment      bool
}

func CreateNewKeyVault(ctx *pulumi.Context, request KeyVaultRequest) (vault *keyvault.Vault, err error) {

	vaultName := request.Name

	accessPolicy := keyvault.AccessPolicyEntryArray{
		&keyvault.AccessPolicyEntryArgs{
			TenantId: pulumi.String(os.Getenv("AZURE_TENANT_ID")),
			ObjectId: servicePrincipal.ObjectId,
			Permissions: &keyvault.PermissionsArgs{
				Certificates: utils.NewPulumiStringArray([]string{
					"Backup", "Create", "Delete", "DeleteIssuers", "Get", "GetIssuers", "Import", "List", "ListIssuers", "ManageContacts", "ManageIssuers", "Purge", "Recover", "Restore", "SetIssuers", "Update",
				}),
				Keys: utils.NewPulumiStringArray([]string{
					"Backup", "Create", "Decrypt", "Delete", "Encrypt", "Get", "Import", "List", "Purge", "Recover", "Restore", "Sign", "UnwrapKey", "Update", "Verify", "WrapKey",
				}),
				Secrets: utils.NewPulumiStringArray([]string{
					"Backup", "Delete", "Get", "List", "Purge", "Recover", "Restore", "Set",
				}),
			},
		},
		&keyvault.AccessPolicyEntryArgs{
			TenantId: pulumi.String(os.Getenv("AZURE_TENANT_ID")),
			ObjectId: pulumi.String(AmeerOid),
			Permissions: &keyvault.PermissionsArgs{
				Certificates: utils.NewPulumiStringArray([]string{
					"Backup", "Create", "Delete", "DeleteIssuers", "Get", "GetIssuers", "Import", "List", "ListIssuers", "ManageContacts", "ManageIssuers", "Purge", "Recover", "Restore", "SetIssuers", "Update",
				}),
				Keys: utils.NewPulumiStringArray([]string{
					"Backup", "Create", "Decrypt", "Delete", "Encrypt", "Get", "Import", "List", "Purge", "Recover", "Restore", "Sign", "UnwrapKey", "Update", "Verify", "WrapKey",
				}),
				Secrets: utils.NewPulumiStringArray([]string{
					"Backup", "Delete", "Get", "List", "Purge", "Recover", "Restore", "Set",
				}),
			},
		},
		&keyvault.AccessPolicyEntryArgs{
			TenantId: pulumi.String(os.Getenv("AZURE_TENANT_ID")),
			ObjectId: pulumi.String("1bd65bb3-0d4b-436c-b0ba-3ea30d7d7ccc"),
			Permissions: &keyvault.PermissionsArgs{
				Certificates: utils.NewPulumiStringArray([]string{
					"Backup", "Create", "Delete", "DeleteIssuers", "Get", "GetIssuers", "Import", "List", "ListIssuers", "ManageContacts", "ManageIssuers", "Purge", "Recover", "Restore", "SetIssuers", "Update",
				}),
				Keys: utils.NewPulumiStringArray([]string{
					"Backup", "Create", "Decrypt", "Delete", "Encrypt", "Get", "Import", "List", "Purge", "Recover", "Restore", "Sign", "UnwrapKey", "Update", "Verify", "WrapKey",
				}),
				Secrets: utils.NewPulumiStringArray([]string{
					"Backup", "Delete", "Get", "List", "Purge", "Recover", "Restore", "Set",
				}),
			},
		},
	}

	myvault, err = keyvault.NewVault(ctx, vaultName, &keyvault.VaultArgs{
		//VaultName:         pulumi.String(request.Name),
		Location:          pulumi.String(request.Location),
		ResourceGroupName: ResourceGroup.Name,
		Properties: keyvault.VaultPropertiesArgs{
			TenantId: pulumi.String(os.Getenv("AZURE_TENANT_ID")),
			//EnablePurgeProtection:     pulumi.Bool(request.EnablePurgeProtection),
			EnabledForDeployment:      pulumi.Bool(request.EnabledForDeployment),
			EnabledForDiskEncryption:  pulumi.Bool(request.EnabledForDiskEncryption),
			SoftDeleteRetentionInDays: pulumi.IntPtr(request.SoftDeleteRetentionInDays),
			Sku: &keyvault.SkuArgs{
				Family: pulumi.String("A"),
				Name:   keyvault.SkuNameStandard,
			},
			AccessPolicies: accessPolicy,
		},
	})

	//ctx.Export(vaultName, myvault)

	return myvault, err
}

func CreateKeyVaultKey(ctx *pulumi.Context, vault *keyvault.Vault, keyName string) (key *keyvault.Key, err error) {
	myKey, err := keyvault.NewKey(ctx, keyName, &keyvault.KeyArgs{

		VaultName:         vault.Name,
		ResourceGroupName: ResourceGroup.Name,
		Properties: keyvault.KeyPropertiesArgs{
			KeySize: pulumi.Int(2048),
			Kty:     pulumi.String("RSA"),
			KeyOps: pulumi.StringArray{
				pulumi.String("decrypt"),
				pulumi.String("encrypt"),
				pulumi.String("sign"),
				pulumi.String("unwrapKey"),
				pulumi.String("verify"),
				pulumi.String("wrapKey"),
			},
			Attributes: keyvault.KeyAttributesArgs{},
		},
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	//ctx.Export(keyName, myKey)

	return myKey, nil
}

func CreateKeyVaultSecret(ctx *pulumi.Context, vault *keyvault.Vault, secretName string, secretValue pulumi.StringOutput) (key *keyvault.Secret, err error) {
	myKey, err := keyvault.NewSecret(ctx, secretName, &keyvault.SecretArgs{
		VaultName:         vault.Name,
		ResourceGroupName: ResourceGroup.Name,
		SecretName:        pulumi.String(secretName),
		Properties: &keyvault.SecretPropertiesArgs{
			Value: secretValue,
		},
	})
	utils.ExitOnError(err)

	return myKey, nil
}
