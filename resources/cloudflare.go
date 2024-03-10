package resources

import (
	"boundary/resources/k8s"
	"boundary/utils"
	"fmt"
	"log"
	"strings"

	"encoding/base64"

	"github.com/pulumi/pulumi-cloudflare/sdk/v4/go/cloudflare"

	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
)

func CreateArgoTunnel(ctx *pulumi.Context, accountId string, tunnelName string) (*cloudflare.ArgoTunnel, pulumi.StringOutput) {

	// Create a random password for the runnel
	tunnelSecret, err := random.NewRandomString(ctx, "random", &random.RandomStringArgs{
		Length:  pulumi.Int(44),
		Special: pulumi.Bool(false),
	})
	if err != nil {
		log.Fatalf(err.Error())
	}

	rndTunnelPostfix, err := random.NewRandomString(ctx, "rndTunnelPostfix", &random.RandomStringArgs{
		Length:  pulumi.Int(6),
		Special: pulumi.Bool(false),
	})
	utils.ExitOnError(err)

	rndTunnelName := pulumi.Sprintf("%s-%s", tunnelName, rndTunnelPostfix.Result)

	// Create an Argo tunnel with the random password
	tunnel, err := cloudflare.NewArgoTunnel(ctx, tunnelName, &cloudflare.ArgoTunnelArgs{
		AccountId: pulumi.String(accountId),
		Name:      rndTunnelName,
		Secret:    tunnelSecret.Result,
	}) //, pulumi.DeleteBeforeReplace(true)
	if err != nil {
		log.Fatal(err.Error())
	}

	return tunnel, tunnelSecret.Result
}

func CreateCloudFlareRecord(ctx *pulumi.Context, name string, zoneId string, value pulumi.StringOutput, recordType string) *cloudflare.Record {
	record, err := cloudflare.NewRecord(ctx, name, &cloudflare.RecordArgs{
		ZoneId:  pulumi.String(zoneId),
		Name:    pulumi.String(name),
		Value:   value,
		Type:    pulumi.String(recordType),
		Ttl:     pulumi.Int(1),
		Proxied: pulumi.Bool(true),
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	return record
}

func NewCloudFlareTunnelSecret(ctx *pulumi.Context, tunnelId pulumi.StringOutput, tunnelName string, tunnelSecret pulumi.StringOutput, accountId string, dependencies []pulumi.Resource, namespace *corev1.Namespace) *corev1.Secret {

	// Generate secret in json form that's required by cloudflared
	value := pulumi.Sprintf("{\"AccountTag\":\"%s\",\"TunnelSecret\":\"%s\",\"TunnelID\":\"%s\",\"TunnelName\":\"%s\"}", accountId, tunnelSecret, tunnelId, tunnelName)

	// Encode to base 64 before writing K8s secret
	encodedValue := value.ApplyT(func(value string) string {
		return base64.StdEncoding.EncodeToString([]byte(value))
	}).(pulumi.StringOutput)

	// Generate K8s secret name based on tunnel name
	secretName := fmt.Sprintf("%s-tunnel-credentials", tunnelName)

	// Create Secret
	secret, err := corev1.NewSecret(ctx, secretName, &corev1.SecretArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name:      pulumi.String(secretName),
			Namespace: namespace.Metadata.Name(),
		},
		Data: pulumi.StringMap{
			"value": encodedValue,
		},
	}, pulumi.DependsOn(dependencies))

	if err != nil {
		log.Fatal(err.Error())
	}

	return secret
}

func DeployCFtunnel(ctx *pulumi.Context, tag string, secret *corev1.Secret, namespace *corev1.Namespace) *appsv1.Deployment {
	// Use role to determine applabel
	appLabel := fmt.Sprintf("boundary-%v-pulumi", tag)

	// Install SelfSigned certs as trusted root
	//cloudflareCmd = "while true; do echo Sleeping...; sleep 5; done"
	controllerUrl := fmt.Sprintf("https://%v.%v:9200", ControllerARecord, DnsZone)

	// Setup arguments for cloudflared tunnel comannds
	cloudflareCmd := fmt.Sprintf("cloudflared tunnel --no-autoupdate run --url %v -f --credentials-file /etc/cloudflared/credentials.json %v", controllerUrl, CloudflareTunnelName)

	args := []string{strings.Join([]string{
		"cp /tls/boundary-controller.crt /usr/local/share/ca-certificates/",
		fmt.Sprintf("export TUNNEL_URL=%s", controllerUrl),
		"export TUNNEL_LOGLEVEL=debug",
		"export TUNNEL_TRANSPORT_LOGLEVEL=debug",
		"export TUNNEL_METRICS=localhost:6000",
		"update-ca-certificates",
		cloudflareCmd,
	}, ";")}

	// Use boundary container from hashicorp
	cSpec := k8s.ContainerSpec{
		Name: fmt.Sprintf("argo-%s", tag),
		//Image: "ghcr.io/katasec/cloudflared:2022.9.0",
		Image: fmt.Sprintf("%s:%s", cloudflareImage, cloudflareImageVerion),
		Command: []string{
			"/bin/sh",
			"-exc",
			"--",
		},
		Args:         args,
		LivenessPort: 6000,
	}

	// Define volume based on k8s secret
	volumeName := "tunnel-credentials"
	volumes, volumeMounts := k8s.GetBoundaryMounts()
	myvolume := &corev1.VolumeArgs{
		Name: pulumi.String(volumeName),
		Secret: &corev1.SecretVolumeSourceArgs{
			SecretName: secret.Metadata.Name(),
			Items: &corev1.KeyToPathArray{
				&corev1.KeyToPathArgs{
					Key:  pulumi.String("value"),
					Path: pulumi.String("credentials.json"),
				},
			},
		},
	}

	// Define mount point for k8s secret volume
	volumeMount := corev1.VolumeMountArgs{
		Name:      pulumi.String(volumeName),
		MountPath: pulumi.String("/etc/cloudflared"),
		ReadOnly:  pulumi.Bool(true),
	}
	volumes = append(volumes, myvolume)
	volumeMounts = append(volumeMounts, volumeMount)

	// Define ENV vars for pod
	envVars := getEnvVars()

	// Define Tolerations for pod placement
	tolerations := k8s.GetTolerations()

	// Define security context
	securityContext := k8s.GetSecurityContext()

	ports := corev1.ContainerPortArray{}

	deployment := k8s.CreateDeploymentCF(ctx, appLabel, cSpec, ports, volumes, volumeMounts, envVars, tolerations, securityContext, []pulumi.Resource{secret}, namespace)

	return deployment
}
