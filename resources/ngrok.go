package resources

import (
	"boundary/resources/k8s"
	"fmt"

	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func DeployNgrokTunnel(
	ctx *pulumi.Context, tag string, remoteAddr string, workerSvc *corev1.Service, ngrokToken string, region string,
	namespace *corev1.Namespace,
) *appsv1.Deployment {

	// svcIpValue := workerSvc.Spec.ClusterIP().ApplyT(func(ip *string) string {
	// 	if ip != nil {
	// 		return *ip
	// 	}
	// 	return ""
	// })

	ctx.Export("LoadBalancer", workerSvc.Status.LoadBalancer())
	ctx.Export("LoadBalancerIP", workerSvc.Spec.LoadBalancerIP())

	// svcIpValue := workerSvc.Status.LoadBalancer().Ingress().Index(pulumi.Int(0)).Ip().ApplyT(func(ip *string) string {
	// 	if ip != nil {
	// 		return *ip
	// 	}
	// 	return WorkerFQDN
	// })

	svcPortValue := workerSvc.Spec.Ports().Index(pulumi.Int(0)).TargetPort().ApplyT(func(port interface{}) int {
		return int(port.(float64))
	})

	workerSvcUrl := pulumi.Sprintf("%s:%d", WorkerFQDN, svcPortValue)

	myArgs := pulumi.StringArray{
		pulumi.String("tcp"),
		pulumi.Sprintf("--region=%s", region),
		pulumi.Sprintf("--remote-addr=%s", remoteAddr),
		pulumi.String("--log=stdout"),
		pulumi.String("--log-format=logfmt"),
		workerSvcUrl,
	}

	applabel := "ngrok-boundary-worker-pulumi"
	cSpec := &k8s.ContainerSpec{
		Name:        applabel,
		Image:       fmt.Sprintf("%s:%s", ngrokImage, ngrokImageVerion),
		DynamicArgs: myArgs,
	}

	// fqdn := postgreSQLserver.Fq
	envVars := &corev1.EnvVarArray{
		corev1.EnvVarArgs{
			Name:  pulumi.String("NGROK_AUTHTOKEN"),
			Value: pulumi.String(ngrokToken),
		},
	}

	tolerations := k8s.GetTolerations()
	securityContext := k8s.GetSecurityContext()
	_ = k8s.CreateDeployment2(ctx, applabel, cSpec, nil, nil, *envVars, tolerations, securityContext, []pulumi.Resource{
		workerSvc,
	}, namespace)

	return nil
}
