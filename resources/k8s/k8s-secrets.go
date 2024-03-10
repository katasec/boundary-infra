package k8s

// func NewCloudFlareTunnelSecret(ctx *pulumi.Context, tunnelId string, tunnelName string, tunnelSecret string, accountId string) *corev1.Secret {

// 	value := fmt.Sprintf("{\"AccountTag\":\"%s\",\"TunnelSecret\":\"%s\",\"TunnelID\":\"%s\",\"TunnelName\":\"%s\"}", accountId, tunnelSecret, tunnelId, tunnelName)

// 	secretName := fmt.Sprintf("%s-tunnel-credentials", tunnelName)

// 	secret, err := corev1.NewSecret(ctx, secretName, &corev1.SecretArgs{
// 		Metadata: v1.ObjectMetaArgs{
// 			Name: pulumi.String(secretName),
// 		},
// 		Data: pulumi.StringMap{
// 			"value": pulumi.String(value),
// 		},
// 	})

// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	return secret
// }
