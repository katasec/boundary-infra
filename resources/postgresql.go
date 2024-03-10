package resources

import (
	utils "boundary/utils"
	"log"
	"strconv"

	"github.com/pulumi/pulumi-azure/sdk/v5/go/azure/postgresql"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreatePostgreSQL(ctx *pulumi.Context, rgName string, location string, serverName string, login pulumi.String, password pulumi.StringOutput) {

	// Create new flexible server
	PostgresFlexibleServer, err := postgresql.NewFlexibleServer(ctx, serverName, &postgresql.FlexibleServerArgs{
		ResourceGroupName:     ResourceGroup.Name,
		Location:              pulumi.String(location),
		Version:               pulumi.String("12"),
		AdministratorLogin:    login,
		AdministratorPassword: password,
		Zone:                  pulumi.String("1"),
		StorageMb:             pulumi.Int(32768),
		SkuName:               pulumi.String("B_Standard_B2s"),
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	// Construct PostgresUrl from output
	//escapedPassword := pulumi.String(url.QueryEscape(string(password)))
	escapedPassword := password
	PostgresUrl = pulumi.Sprintf("postgres://%s:%s@%s:5432/postgres?sslmode=require", login, escapedPassword, PostgresFlexibleServer.Fqdn)

	// PGCrypto Extensions
	_, err = postgresql.NewFlexibleServerConfiguration(ctx, "azure.extensions", &postgresql.FlexibleServerConfigurationArgs{
		ServerId: PostgresFlexibleServer.ID(),
		Name:     pulumi.String("azure.extensions"),
		Value:    pulumi.String("PGCRYPTO"),
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	// Add my ip to firewall
	ip := utils.HttpGet("https://ifconfig.me")
	CreatePostgreSQLFirewallRule(ctx, PostgresFlexibleServer, "ALLOW_ME", ip, ip)

	// // Add ALL ips to firewall
	// allIp := "0.0.0.0/0"
	// CreatePostgreSQLFirewallRule(ctx, PostgresFlexibleServer, "ALLOW_ALL", allIp, allIp)

	// Add Allowed Ips to PostgreSQL rules
	for i, ip := range allowedIps {
		CreatePostgreSQLFirewallRule(ctx, PostgresFlexibleServer, "AZ_FW_IP_"+strconv.Itoa(i), ip, ip)
	}

}

func CreatePostgreSQLFirewallRule(ctx *pulumi.Context, flexibleServer *postgresql.FlexibleServer, ruleName string, startIp string, endIp string) {
	var err error
	PostgreSQLFirewallRules, err = postgresql.NewFlexibleServerFirewallRule(ctx, ruleName, &postgresql.FlexibleServerFirewallRuleArgs{
		ServerId:       flexibleServer.ID(),
		StartIpAddress: pulumi.String(startIp),
		EndIpAddress:   pulumi.String(endIp),
	})

	if err != nil {
		log.Fatalf(err.Error())
	}

}
