package resources

import (
	"boundary/resources/k8s"
	"boundary/utils"
	"fmt"
	"strings"

	batchv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/batch/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	v1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func BoundaryDbInit(ctx *pulumi.Context, namespace *corev1.Namespace) *batchv1.Job {

	volumes, volumeMounts := k8s.GetBoundaryMounts()
	tolerations := k8s.GetTolerations()
	envVars := getEnvVars()

	args := []string{strings.Join([]string{
		"env | grep -i azure",
		"boundary database init -config /boundary/controller.hcl",
	}, ";")}

	boundaryDbInitJob, err := batchv1.NewJob(ctx, "boundary-database-init-pulumi", &batchv1.JobArgs{
		Metadata: &v1.ObjectMetaArgs{
			Namespace: namespace.Metadata.Name(),
		},
		Spec: &batchv1.JobSpecArgs{
			BackoffLimit: pulumi.Int(0),
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: v1.ObjectMetaArgs{
					Labels: pulumi.StringMap{
						dbJobKey: pulumi.String(dbJobValue),
					},
					Namespace: namespace.Metadata.Name(),
				},
				Spec: &corev1.PodSpecArgs{
					Volumes: volumes,
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Command: pulumi.StringArray{
								pulumi.String("/bin/sh"),
								pulumi.String("-xc"),
								pulumi.String("--"),
							},
							Args:  utils.NewPulumiStringArray(args),
							Image: pulumi.String(fmt.Sprintf("%s:%s", boundaryImage, boundaryImageVersion)),
							Name:  pulumi.String("boundary-database-init-pulumi"),
							//SecurityContext: k8s.GetSecurityContext(),
							VolumeMounts: volumeMounts,
							Env:          envVars,
							Resources: &corev1.ResourceRequirementsArgs{
								Limits: pulumi.StringMap{
									"cpu":    pulumi.String("1"),
									"memory": pulumi.String("512Mi"),
								},
								Requests: pulumi.StringMap{
									"memory": pulumi.String("256Mi"),
									"cpu":    pulumi.String("300m"),
								},
							},
						},
					},
					RestartPolicy: pulumi.String("Never"),
					Tolerations:   tolerations,
				},
			},
			Completions: pulumi.Int(1),
		},
	}, pulumi.DependsOn([]pulumi.Resource{
		PostgreSQLFirewallRules,
		myvault,
		RootKey,
		WorkerKey,
		RecoveryKey,
	}))

	utils.ExitOnError(err)

	return boundaryDbInitJob
}
