// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

#pragma once
#include "opentelemetry/exporters/otlp/otlp_http_exporter_factory.h"
#include "opentelemetry/context/propagation/global_propagator.h"
#include "opentelemetry/context/propagation/text_map_propagator.h"
#include "opentelemetry/exporters/ostream/span_exporter_factory.h"
#include "opentelemetry/nostd/shared_ptr.h"
#include "opentelemetry/sdk/trace/simple_processor_factory.h"
#include "opentelemetry/sdk/trace/tracer_context.h"
#include "opentelemetry/sdk/trace/tracer_context_factory.h"
#include "opentelemetry/sdk/trace/tracer_provider_factory.h"
#include "opentelemetry/trace/propagation/http_trace_context.h"
#include "opentelemetry/trace/provider.h"
#include "opentelemetry/sdk/resource/resource.h"
#include "opentelemetry/semconv/resource_attributes.h"

#include <grpcpp/grpcpp.h>
#include <cstring>
#include <cstdlib>
#include <iostream>
#include <vector>
#include <map>

using grpc::ClientContext;
using grpc::ServerContext;

namespace
{
class GrpcClientCarrier : public opentelemetry::context::propagation::TextMapCarrier
{
public:
  GrpcClientCarrier(ClientContext *context) : context_(context) {}
  GrpcClientCarrier() = default;
  virtual opentelemetry::nostd::string_view Get(
      opentelemetry::nostd::string_view key) const noexcept override
  {
    return "";
  }

  virtual void Set(opentelemetry::nostd::string_view key,
                   opentelemetry::nostd::string_view value) noexcept override
  {
    std::cout << " Client ::: Adding " << key << " " << value << "\n";
    context_->AddMetadata(key.data(), value.data());
  }

  ClientContext *context_;
};

class GrpcServerCarrier : public opentelemetry::context::propagation::TextMapCarrier
{
public:
  GrpcServerCarrier(ServerContext *context) : context_(context) {}
  GrpcServerCarrier() = default;
  virtual opentelemetry::nostd::string_view Get(
      opentelemetry::nostd::string_view key) const noexcept override
  {
    auto it = context_->client_metadata().find(key.data());
    if (it != context_->client_metadata().end())
    {
      return it->second.data();
    }
    return "";
  }

  virtual void Set(opentelemetry::nostd::string_view key,
                   opentelemetry::nostd::string_view value) noexcept override
  {
   // Not required for server
  }

  ServerContext *context_;
};

opentelemetry::sdk::resource::Resource createResource()
{
  std::map<std::string, std::string> resource_attributes;
  
  // Set default service name
  resource_attributes[opentelemetry::semconv::resource::kServiceName] = "currency";
  
  // Check for environment variables and override defaults
  const char* service_name = std::getenv("OTEL_SERVICE_NAME");
  if (service_name) {
    resource_attributes[opentelemetry::semconv::resource::kServiceName] = service_name;
  }
  
  const char* service_version = std::getenv("OTEL_SERVICE_VERSION");
  if (service_version) {
    resource_attributes[opentelemetry::semconv::resource::kServiceVersion] = service_version;
  }
  
  const char* cms_workspace = std::getenv("CMS_WORKSPACE");
  if (cms_workspace) {
    resource_attributes["acs.cms.workspace"] = cms_workspace;
  }
  
  const char* deployment_environment = std::getenv("OTEL_DEPLOYMENT_ENVIRONMENT");
  if (deployment_environment) {
    resource_attributes[opentelemetry::semconv::resource::kDeploymentEnvironment] = deployment_environment;
  }
  
  const char* service_namespace = std::getenv("OTEL_SERVICE_NAMESPACE");
  if (service_namespace) {
    resource_attributes[opentelemetry::semconv::resource::kServiceNamespace] = service_namespace;
  }
  
  const char* service_instance_id = std::getenv("OTEL_SERVICE_INSTANCE_ID");
  if (service_instance_id) {
    resource_attributes[opentelemetry::semconv::resource::kServiceInstanceId] = service_instance_id;
  }

  return opentelemetry::sdk::resource::Resource::Create(resource_attributes);
}

void initTracer()
{
  auto resource = createResource();
  auto exporter = opentelemetry::exporter::otlp::OtlpHttpExporterFactory::Create();
  auto processor =
      opentelemetry::sdk::trace::SimpleSpanProcessorFactory::Create(std::move(exporter));
  std::vector<std::unique_ptr<opentelemetry::sdk::trace::SpanProcessor>> processors;
  processors.push_back(std::move(processor));

  auto context =
      opentelemetry::sdk::trace::TracerContextFactory::Create(std::move(processors), resource);
  std::shared_ptr<opentelemetry::trace::TracerProvider> provider =
      opentelemetry::sdk::trace::TracerProviderFactory::Create(std::move(context));

  // Set the global trace provider
  opentelemetry::trace::Provider::SetTracerProvider(provider);

  // set global propagator
  opentelemetry::context::propagation::GlobalTextMapPropagator::SetGlobalPropagator(
      opentelemetry::nostd::shared_ptr<opentelemetry::context::propagation::TextMapPropagator>(
          new opentelemetry::trace::propagation::HttpTraceContext()));
}

opentelemetry::nostd::shared_ptr<opentelemetry::trace::Tracer> get_tracer(std::string tracer_name)
{
  auto provider = opentelemetry::trace::Provider::GetTracerProvider();
  return provider->GetTracer(tracer_name);
}

} // namespace
