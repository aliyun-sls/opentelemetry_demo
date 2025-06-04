// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

#include "opentelemetry/exporters/otlp/otlp_http_exporter_factory.h"
#include "opentelemetry/logs/provider.h"
#include "opentelemetry/sdk/logs/logger.h"
#include "opentelemetry/sdk/logs/logger_provider_factory.h"
#include "opentelemetry/sdk/logs/simple_log_record_processor_factory.h"
#include "opentelemetry/sdk/logs/logger_context_factory.h"
#include "opentelemetry/exporters/otlp/otlp_http_log_record_exporter_factory.h"
#include "opentelemetry/sdk/resource/resource.h"

#include <cstdlib>
#include <map>

using namespace std;
namespace nostd     = opentelemetry::nostd;
namespace otlp      = opentelemetry::exporter::otlp;
namespace logs      = opentelemetry::logs;
namespace logs_sdk  = opentelemetry::sdk::logs;

namespace
{
  opentelemetry::sdk::resource::Resource createLogResource()
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

  void initLogger() {
    auto resource = createLogResource();
    
    otlp::OtlpHttpLogRecordExporterOptions loggerOptions;
    auto exporter  = otlp::OtlpHttpLogRecordExporterFactory::Create(loggerOptions);
    auto processor = logs_sdk::SimpleLogRecordProcessorFactory::Create(std::move(exporter));
    std::vector<std::unique_ptr<logs_sdk::LogRecordProcessor>> processors;
    processors.push_back(std::move(processor));
    auto context = logs_sdk::LoggerContextFactory::Create(std::move(processors), resource);
    std::shared_ptr<logs::LoggerProvider> provider = logs_sdk::LoggerProviderFactory::Create(std::move(context));
    opentelemetry::logs::Provider::SetLoggerProvider(provider);
  }

  nostd::shared_ptr<opentelemetry::logs::Logger> getLogger(std::string name){
    auto provider = logs::Provider::GetLoggerProvider();
    return provider->GetLogger(name + "_logger", name, OPENTELEMETRY_SDK_VERSION);
  }
}
