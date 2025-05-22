// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

#pragma once
#include "opentelemetry/exporters/otlp/otlp_http_exporter_factory.h"
#include "opentelemetry/context/propagation/global_propagator.h"
#include "opentelemetry/context/propagation/text_map_propagator.h"
#include "opentelemetry/exporters/ostream/span_exporter_factory.h"
#include "opentelemetry/nostd/shared_ptr.h"
#include "opentelemetry/sdk/trace/batch_span_processor_factory.h"
#include "opentelemetry/sdk/trace/tracer_context.h"
#include "opentelemetry/sdk/trace/tracer_context_factory.h"
#include "opentelemetry/sdk/trace/tracer_provider_factory.h"
#include "opentelemetry/trace/propagation/http_trace_context.h"
#include "opentelemetry/trace/provider.h"
#include "opentelemetry/sdk/resource/resource.h"
#include "opentelemetry/sdk/resource/resource_detector.h"
#include "opentelemetry/sdk/resource/semantic_conventions.h"

#include <grpcpp/grpcpp.h>
#include <cstring>
#include <iostream>
#include <vector>
#include <cstdlib>
#include <unistd.h>
#include <netdb.h>
#include <arpa/inet.h>

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

std::string get_hostname() {
    char hostname[256];
    if (gethostname(hostname, sizeof(hostname)) == 0) {
        return std::string(hostname);
    }
    return "unknown";
}

std::string get_ip_address() {
    char hostname[256];
    if (gethostname(hostname, sizeof(hostname)) == 0) {
        struct hostent *host = gethostbyname(hostname);
        if (host != nullptr) {
            struct in_addr addr;
            memcpy(&addr, host->h_addr_list[0], sizeof(struct in_addr));
            return std::string(inet_ntoa(addr));
        }
    }
    return "unknown";
}

void initTracer()
{
  opentelemetry::exporter::otlp::OtlpHttpExporterOptions opts;
  auto exporter = opentelemetry::exporter::otlp::OtlpHttpExporterFactory::Create(opts);
  
  // 配置 BatchSpanProcessor 选项
  opentelemetry::sdk::trace::BatchSpanProcessorOptions batch_options;
  batch_options.max_queue_size = 2048;  // 最大队列大小
  batch_options.schedule_delay_millis = std::chrono::milliseconds(5000);  // 调度延迟
  batch_options.max_export_batch_size = 512;  // 最大导出批次大小
  
  auto processor = opentelemetry::sdk::trace::BatchSpanProcessorFactory::Create(
      std::move(exporter), batch_options);
      
  std::vector<std::unique_ptr<opentelemetry::sdk::trace::SpanProcessor>> processors;
  processors.push_back(std::move(processor));

  // 获取环境变量中的 workspace
  const char* workspace = std::getenv("ACS_CMS_WORKSPACE");
  std::string workspace_str = workspace ? workspace : "default";

  // 创建默认的 resource
  opentelemetry::sdk::resource::ResourceAttributes attributes;
  attributes[opentelemetry::sdk::resource::SemanticConventions::HOST_NAME] = get_hostname();
  attributes[opentelemetry::sdk::resource::SemanticConventions::HOST_IP] = get_ip_address();
  attributes["acs_cms_workspace"] = workspace_str;
  auto resource = opentelemetry::sdk::resource::Resource::Create(attributes);

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
