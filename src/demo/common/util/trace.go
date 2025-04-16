package util

import (
	"github.com/aliyun-sls/opentelemetry-go-provider-sls/provider"
	"os"
	"sls-mall-go/common/config"
)

func InitTrace() {
	project := config.TraceProject
	endpoint := config.TraceEndpoint
	instance := config.TraceInstance
	resourceAttributes := map[string]string{}
	resourceAttributes["telemetry.sdk.language"] = "go"
	resourceAttributes["telemetry.sdk.name"] = "opentelemetry"
	resourceAttributes["telemetry.sdk.version"] = "1.14.0"
	resourceAttributes["k8s.container.name"] = os.Getenv("CONTAINER_NAME")
	resourceAttributes["k8s.deployment.name"] = os.Getenv("DEPLOYMENT_NAME")
	resourceAttributes["k8s.namespace.name"] = os.Getenv("POD_NAMESPACE")
	resourceAttributes["k8s.node.name"] = os.Getenv("NODE_NAME")
	resourceAttributes["k8s.pod.name"] = os.Getenv("POD_NAME")
	slsConfig, err := provider.NewConfig(provider.WithServiceName(config.ServiceName),
		provider.WithServiceNamespace("sls-mall-go"),
		provider.WithServiceVersion("v0.0.1"),
		provider.WithTraceExporterEndpoint(endpoint),
		provider.WithMetricExporterEndpoint(endpoint),
		provider.WithSLSConfig(project, instance, config.AccessKeyID, config.AccessKeySecret),
		provider.WithResourceAttributes(resourceAttributes))

	Chk(err)
	err = provider.Start(slsConfig)
	Chk(err)
}
