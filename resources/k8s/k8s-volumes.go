package k8s

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func GetBoundaryMounts() (volumes corev1.VolumeArray, volumeMounts corev1.VolumeMountArray) {

	/*
		The below volumes are used to store boundary SSL certs (keys, crt files)
		as well as boundary config files (worker and controller)
	*/

	// Define volumes created from the config maps
	volumes = corev1.VolumeArray{
		corev1.VolumeArgs{
			Name: pulumi.String("boundary-config-controller"),
			ConfigMap: corev1.ConfigMapVolumeSourceArgs{
				Name: pulumi.String("boundary-config-controller"),
			},
		},
		corev1.VolumeArgs{
			Name: pulumi.String("boundary-config-worker"),
			ConfigMap: corev1.ConfigMapVolumeSourceArgs{
				Name: pulumi.String("boundary-config-worker"),
			},
		},
		corev1.VolumeArgs{
			Name: pulumi.String("boundary-crt-worker"),
			ConfigMap: corev1.ConfigMapVolumeSourceArgs{
				Name: pulumi.String("boundary-crt-worker"),
			},
		},
		corev1.VolumeArgs{
			Name: pulumi.String("boundary-crt-controller"),
			ConfigMap: corev1.ConfigMapVolumeSourceArgs{
				Name: pulumi.String("boundary-crt-controller"),
			},
		},
		corev1.VolumeArgs{
			Name: pulumi.String("boundary-key"),
			ConfigMap: corev1.ConfigMapVolumeSourceArgs{
				Name: pulumi.String("boundary-key"),
			},
		},
	}

	// Define volume mounts for the volumes above
	volumeMounts = corev1.VolumeMountArray{
		corev1.VolumeMountArgs{
			Name:      pulumi.String("boundary-config-controller"),
			MountPath: pulumi.String("/boundary/controller.hcl"),
			SubPath:   pulumi.String("controller.hcl"),
			ReadOnly:  pulumi.Bool(false),
		},
		corev1.VolumeMountArgs{
			Name:      pulumi.String("boundary-config-worker"),
			MountPath: pulumi.String("/boundary/worker.hcl"),
			SubPath:   pulumi.String("worker.hcl"),
			ReadOnly:  pulumi.Bool(false),
		},
		corev1.VolumeMountArgs{
			Name:      pulumi.String("boundary-crt-controller"),
			MountPath: pulumi.String("/tls/boundary-controller.crt"),
			SubPath:   pulumi.String("boundary-controller.crt"),
			ReadOnly:  pulumi.Bool(false),
		},
		corev1.VolumeMountArgs{
			Name:      pulumi.String("boundary-crt-worker"),
			MountPath: pulumi.String("/tls/boundary-worker.crt"),
			SubPath:   pulumi.String("boundary-worker.crt"),
			ReadOnly:  pulumi.Bool(false),
		},
		corev1.VolumeMountArgs{
			Name:      pulumi.String("boundary-key"),
			MountPath: pulumi.String("/tls/boundary.key"),
			SubPath:   pulumi.String("boundary.key"),
			ReadOnly:  pulumi.Bool(false),
		},
	}

	return volumes, volumeMounts
}
