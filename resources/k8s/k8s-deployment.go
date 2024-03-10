package k8s

import (
	"log"

	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	v1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ContainerSpec struct {
	Name         string
	Image        string
	Command      []string
	Args         []string
	DynamicArgs  pulumi.StringArray
	LivenessPort int
}

func CreateDeployment(ctx *pulumi.Context, appLabel string, cSpec ContainerSpec, ports corev1.ContainerPortArray, volumes corev1.VolumeArray, volumeMounts corev1.VolumeMountArray, envVars corev1.EnvVarArray, tolerations corev1.TolerationArray, securityContext corev1.SecurityContextArgs, dependencies []pulumi.Resource, ns *corev1.Namespace) *appsv1.Deployment {

	selectLabel := pulumi.StringMap{
		"app": pulumi.String(appLabel),
	}

	// label := pulumi.StringMap{
	// 	"app": pulumi.String(appLabel + "-svc"),
	// }

	myDeployment, err := appsv1.NewDeployment(ctx, appLabel, &appsv1.DeploymentArgs{
		Metadata: &v1.ObjectMetaArgs{
			Namespace: ns.Metadata.Name(),
		},
		Spec: appsv1.DeploymentSpecArgs{
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: selectLabel,
			},
			Replicas: pulumi.Int(1),
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels:    selectLabel,
					Namespace: ns.Metadata.Name(),
				},
				Spec: &corev1.PodSpecArgs{
					Volumes: volumes,
					Containers: corev1.ContainerArray{
						corev1.ContainerArgs{
							Name:    pulumi.String(cSpec.Name),
							Image:   pulumi.String(cSpec.Image),
							Command: pulumi.ToStringArray(cSpec.Command),
							Args:    pulumi.ToStringArray(cSpec.Args),
							//SecurityContext: securityContext,
							VolumeMounts: volumeMounts,
							Env:          envVars,
							Ports:        ports,
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
					Tolerations: tolerations,
				},
			},
		},
	}, pulumi.DependsOn(dependencies))

	if err != nil {
		log.Fatal(err.Error())
	}

	return myDeployment
}

func CreateDeployment2(
	ctx *pulumi.Context, appLabel string, cSpec *ContainerSpec, volumes corev1.VolumeArray, volumeMounts corev1.VolumeMountArray, envVars corev1.EnvVarArray,
	tolerations corev1.TolerationArray, securityContext corev1.SecurityContextArgs,
	dependencies []pulumi.Resource, namespace *corev1.Namespace) *appsv1.Deployment {

	appLabels := pulumi.StringMap{
		"app": pulumi.String(appLabel),
	}

	myDeployment, err := appsv1.NewDeployment(ctx, appLabel, &appsv1.DeploymentArgs{
		Metadata: &v1.ObjectMetaArgs{
			Namespace: namespace.Metadata.Name(),
		},
		Spec: appsv1.DeploymentSpecArgs{
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: appLabels,
			},
			Replicas: pulumi.Int(1),
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels:    appLabels,
					Namespace: namespace.Metadata.Name(),
				},
				Spec: &corev1.PodSpecArgs{
					Volumes: volumes,
					Containers: corev1.ContainerArray{
						corev1.ContainerArgs{
							Name:         pulumi.String(cSpec.Name),
							Image:        pulumi.String(cSpec.Image),
							Command:      pulumi.ToStringArray(cSpec.Command),
							Args:         cSpec.DynamicArgs,
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
					Tolerations: tolerations,
				},
			},
		},
	}, pulumi.DependsOn(dependencies))
	if err != nil {
		log.Fatal(err.Error())
	}
	//ctx.Export(appLabel+"_deployment", myDeployment.Metadata.Elem().Name())
	return myDeployment
}
func GetTolerations() corev1.TolerationArray {
	return corev1.TolerationArray{
		corev1.TolerationArgs{
			Key:      pulumi.String("platform"),
			Operator: pulumi.String("Equal"),
			Value:    pulumi.String("lin01"),
			Effect:   pulumi.String("NoSchedule"),
		},
	}
}

func GetSecurityContext() corev1.SecurityContextArgs {
	return corev1.SecurityContextArgs{
		Capabilities: corev1.CapabilitiesArgs{
			Add: pulumi.ToStringArray([]string{"IPC_LOCK"}),
		},
	}
}

func CreateDeploymentCF(ctx *pulumi.Context, appLabel string, cSpec ContainerSpec, ports corev1.ContainerPortArray, volumes corev1.VolumeArray, volumeMounts corev1.VolumeMountArray, envVars corev1.EnvVarArray, tolerations corev1.TolerationArray, securityContext corev1.SecurityContextArgs, dependencies []pulumi.Resource, namespace *corev1.Namespace) *appsv1.Deployment {

	appLabels := pulumi.StringMap{
		"app": pulumi.String(appLabel),
	}

	myDeployment, err := appsv1.NewDeployment(ctx, appLabel, &appsv1.DeploymentArgs{
		Metadata: &v1.ObjectMetaArgs{
			Namespace: namespace.Metadata.Name(),
		},
		Spec: appsv1.DeploymentSpecArgs{
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: appLabels,
			},
			Replicas: pulumi.Int(1),
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels:    appLabels,
					Namespace: namespace.Metadata.Name(),
				},
				Spec: &corev1.PodSpecArgs{
					Volumes: volumes,
					Containers: corev1.ContainerArray{
						corev1.ContainerArgs{
							Name:    pulumi.String(cSpec.Name),
							Image:   pulumi.String(cSpec.Image),
							Command: pulumi.ToStringArray(cSpec.Command),
							Args:    pulumi.ToStringArray(cSpec.Args),
							//SecurityContext: securityContext,
							VolumeMounts: volumeMounts,
							Env:          envVars,
							Ports:        ports,
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
					Tolerations: tolerations,
				},
			},
		},
	}, pulumi.DependsOn(dependencies))

	if err != nil {
		log.Fatal(err.Error())
	}
	ctx.Export(appLabel+"_deployment", myDeployment.Metadata.Elem().Name())
	return myDeployment
}
