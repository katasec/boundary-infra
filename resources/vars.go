package resources

import (
	"fmt"
	"os"
	"strings"

	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/keyvault"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	"github.com/pulumi/pulumi-azure/sdk/v5/go/azure/postgresql"
	"github.com/pulumi/pulumi-azuread/sdk/v5/go/azuread"
	"github.com/pulumi/pulumi-tls/sdk/v4/go/tls"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	privateKey                    *tls.PrivateKey
	Controller_cert               *tls.SelfSignedCert
	Worker_cert                   *tls.SelfSignedCert
	ResourceGroup                 *resources.ResourceGroup
	PostgresFlexibleServer        *postgresql.FlexibleServer
	PostgresAdministratorLogin    pulumi.String
	PostgresAdministratorPassword pulumi.String
	PostgresUrl                   pulumi.StringOutput
	azureadAppDisplayName         = os.Getenv("AZURE_APP_DISPLAY_NAME")
	servicePrincipal              *azuread.ServicePrincipal
	servicePrincipalPassword      *azuread.ServicePrincipalPassword
	AmeerOid                      = os.Getenv("AZURE_USER_OID")
	myvault                       *keyvault.Vault
	PostgreSQLFirewallRules       *postgresql.FlexibleServerFirewallRule

	RootKey     *keyvault.Key
	WorkerKey   *keyvault.Key
	RecoveryKey *keyvault.Key

	// CLoudflare Provider
	CLoudflareEmail  = os.Getenv("CLOUDFLARE_EMAIL")
	CLoudflareApiKey = os.Getenv("CLOUDFLARE_ZONE_ID")

	// Cloudflare  Inputs
	CloudflareAccountId  = os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	CloudflareZoneId     = os.Getenv("CLOUDFLARE_ZONE_ID")
	CloudflareDnszone    = os.Getenv("CLOUDFLARE_DNS_ZONE")
	CloudflareTunnelName = "boundary"

	// DNS
	WorkerARecord     = "boundary-worker"
	ControllerARecord = "boundary-controller"
	DnsZone           = os.Getenv("DNS_ZONE")
	WorkerFQDN        = os.Getenv("DNS_WORKER_FQDN")
	WorkerDnsRG       = os.Getenv("DNS_WORKER_DNS_RG")

	BoundaryUrl = fmt.Sprintf("https://%s.%s", CloudflareTunnelName, CloudflareDnszone)

	// Boundary Version
	// boundaryVersion = "0.8.1"
	boundaryImage        = os.Getenv("BOUNDARY_IMAGE")
	boundaryImageVersion = os.Getenv("BOUNDARY_IMAGE_VERSION")

	cloudflareImage       = os.Getenv("CLOUDFLARE_IMAGE")
	cloudflareImageVerion = os.Getenv("CLOUDFLARE_IMAGE_VERSION")

	ngrokImage       = os.Getenv("NGROK_IMAGE")
	ngrokImageVerion = os.Getenv("NGROK_IMAGE_VERSION")

	// Tag for databasejob
	dbJobKey   = "resource-identifier"
	dbJobValue = "boundary-database-job"

	// Database - Allowed IPs (Azure Firewall IPs)
	allowedIps = strings.Split(os.Getenv("ALLOWED_IPS"), ",")

	// Ngrok Config
	//NGROK_AUTHTOKEN  <- Token is being referenced from this environment variable
	Ngrok_Auth_Token           = os.Getenv("NGROK_AUTHTOKEN")
	Ngrok_Worker_PublicAddress = os.Getenv("NGROK_WORKER_PUBLICADDRESS")
	Ngrok_Edge_Label           = os.Getenv("NGROK_EDGE_LABEL")
	Ngrok_Region               = os.Getenv("NGROK_REGION")
)
