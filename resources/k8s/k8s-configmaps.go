package k8s

import (
	"boundary/utils"
	"log"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateConfigMapFromFile(ctx *pulumi.Context, configName string, configKey string, fileName string, config utils.BoundaryConfigFile, namespace *corev1.Namespace) {

	// Use text templating to generate config file from inputs
	fileContent := utils.RenderTemplate(fileName, config)

	// Create configmap from config file
	configKeys := map[string]string{configKey: fileContent}
	CreateConfigMap(ctx, configName, configKeys, namespace)
}

func CreateConfigMap(ctx *pulumi.Context, configName string, configKeys map[string]string, namespace *corev1.Namespace) {

	_, err := corev1.NewConfigMap(ctx, configName, &corev1.ConfigMapArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name:      pulumi.String(configName),
			Namespace: namespace.Metadata.Name(),
		},
		Data: pulumi.ToStringMap(configKeys),
	}, nil)

	if err != nil {
		log.Fatal(err.Error())
	}

}

func CreateConfigMapKV(ctx *pulumi.Context, configName string, key string, value pulumi.StringOutput, namespace *corev1.Namespace) {

	_, err := corev1.NewConfigMap(ctx, configName, &corev1.ConfigMapArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name:      pulumi.String(configName),
			Namespace: namespace.Metadata.Name(),
		},
		Data: pulumi.StringMap{
			key: value,
		},
	}, nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	//ctx.Export(configName, configMap)
}
