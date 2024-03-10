package resources

import (
	"boundary/utils"

	tls "github.com/pulumi/pulumi-tls/sdk/v4/go/tls"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreatePrivateKey(ctx *pulumi.Context) {
	// Create Private Key
	privateKey, _ = tls.NewPrivateKey(ctx, "boundary_private_key", &tls.PrivateKeyArgs{
		Algorithm: pulumi.String("RSA"),
	})
}
func CreateSelfSignedCert(ctx *pulumi.Context, certName, dnsName string, org string) (sscert *tls.SelfSignedCert) {

	// Use Private Key to generate cert
	cert, err := tls.NewSelfSignedCert(
		ctx,
		certName,
		&tls.SelfSignedCertArgs{
			AllowedUses: pulumi.StringArray{
				pulumi.String("key_encipherment"),
				pulumi.String("digital_signature"),
				pulumi.String("server_auth"),
				pulumi.String("client_auth"),
			},
			DnsNames: pulumi.StringArray{
				pulumi.String(dnsName),
			},
			Uris: pulumi.StringArray{
				pulumi.String(dnsName),
			},
			Subject: &tls.SelfSignedCertSubjectArgs{
				CommonName:   pulumi.String(dnsName),
				Organization: pulumi.String(org),
			},
			PrivateKeyPem:       privateKey.PrivateKeyPem,
			ValidityPeriodHours: pulumi.Int(8760), // 8760 hours = 1 year
		},
		nil,
	)

	utils.ExitOnError(err)
	return cert
}
