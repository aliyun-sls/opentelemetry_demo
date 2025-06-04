// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

#include "opentelemetry/exporters/otlp/otlp_http_metric_exporter_factory.h"
#include "opentelemetry/metrics/provider.h"
#include "opentelemetry/sdk/metrics/export/periodic_exporting_metric_reader.h"
#include "opentelemetry/sdk/metrics/meter.h"
#include "opentelemetry/sdk/metrics/meter_provider.h"
#include "opentelemetry/sdk/metrics/view/view_registry.h"
#include "opentelemetry/sdk/resource/resource.h"

#include <cstdlib>
#include <map>

// namespaces
namespace common        = opentelemetry::common;
namespace metrics_api   = opentelemetry::metrics;
namespace metric_sdk    = opentelemetry::sdk::metrics;
namespace nostd         = opentelemetry::nostd;
namespace otlp_exporter = opentelemetry::exporter::otlp;

namespace
{
  opentelemetry::sdk::resource::Resource createMetricResource()
  {
    opentelemetry::sdk::resource::ResourceAttributes resource_attributes;
    
    // Set default service name
    resource_attributes["service.name"] = "currency";
    
    // Check for environment variables and override defaults
    const char* service_name = std::getenv("OTEL_SERVICE_NAME");
    if (service_name) {
      resource_attributes["service.name"] = service_name;
    }
    
    const char* service_version = std::getenv("OTEL_SERVICE_VERSION");
    if (service_version) {
      resource_attributes["service.version"] = service_version;
    }
    
    const char* cms_workspace = std::getenv("CMS_WORKSPACE");
    if (cms_workspace) {
      resource_attributes["acs.cms.workspace"] = cms_workspace;
    }

    const char* deployment_environment = std::getenv("OTEL_DEPLOYMENT_ENVIRONMENT");
    if (deployment_environment) {
      resource_attributes["deployment.environment"] = deployment_environment;
    }
    
    const char* service_namespace = std::getenv("OTEL_SERVICE_NAMESPACE");
    if (service_namespace) {
      resource_attributes["service.namespace"] = service_namespace;
    }
    
    const char* service_instance_id = std::getenv("OTEL_SERVICE_INSTANCE_ID");
    if (service_instance_id) {
      resource_attributes["service.instance.id"] = service_instance_id;
    }

    return opentelemetry::sdk::resource::Resource::Create(resource_attributes);
  }

  void initMeter() 
  {
    auto resource = createMetricResource();
    
    // Build MetricExporter
    otlp_exporter::OtlpHttpMetricExporterOptions otlpOptions;
    auto exporter = otlp_exporter::OtlpHttpMetricExporterFactory::Create(otlpOptions);

    // Build MeterProvider and Reader
    metric_sdk::PeriodicExportingMetricReaderOptions options;
    std::unique_ptr<metric_sdk::MetricReader> reader{
        new metric_sdk::PeriodicExportingMetricReader(std::move(exporter), options) };
    
    // Create MeterProvider with ViewRegistry and Resource
    auto provider = std::shared_ptr<metrics_api::MeterProvider>(
        new metric_sdk::MeterProvider(
            std::unique_ptr<metric_sdk::ViewRegistry>(new metric_sdk::ViewRegistry()),
            resource
        )
    );
    auto p = std::static_pointer_cast<metric_sdk::MeterProvider>(provider);
    p->AddMetricReader(std::move(reader));
    metrics_api::Provider::SetMeterProvider(provider);
  }

  nostd::unique_ptr<metrics_api::Counter<uint64_t>> initIntCounter(std::string name, std::string version)
  {
    std::string counter_name = name + "_counter";
    auto provider = metrics_api::Provider::GetMeterProvider();
    nostd::shared_ptr<metrics_api::Meter> meter = provider->GetMeter(name, version);
    auto int_counter = meter->CreateUInt64Counter(counter_name);
    return int_counter;
  }
}
