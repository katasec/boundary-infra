package k8s

import (
	"fmt"
	"log"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	v1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateK8Service(ctx *pulumi.Context, appName string, appLabel string, protocol string, ports corev1.ServicePortArray, dependencies []pulumi.Resource, namespace *corev1.Namespace) *corev1.Service {

	myspecArgs := &corev1.ServiceSpecArgs{
		Ports: ports,
		Selector: pulumi.StringMap{
			"app": pulumi.String(appLabel),
		},
		Type: corev1.ServiceSpecTypeLoadBalancer,
	}

	// jsonBytes, _ := json.MarshalIndent(myspecArgs, "", "    ")
	// fmt.Println(string(jsonBytes))

	svcName := fmt.Sprintf("%s-svc", appName)
	service, err := corev1.NewService(ctx, svcName, &corev1.ServiceArgs{
		Metadata: v1.ObjectMetaArgs{
			Namespace: namespace.Metadata.Name(),
			Annotations: pulumi.ToStringMap(map[string]string{
				"service.beta.kubernetes.io/azure-load-balancer-internal": "true",
			}),
		},
		Spec: myspecArgs,
	}, pulumi.DependsOn(dependencies))

	if err != nil {

		log.Fatal(err.Error())
	}

	//ctx.Export(svcName, service)
	return service
}
