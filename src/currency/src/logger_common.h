// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

#pragma once
#include "opentelemetry/exporters/otlp/otlp_http_exporter_factory.h"
#include "opentelemetry/exporters/otlp/otlp_http_exporter_options.h"
#include "opentelemetry/logs/provider.h"
#include "opentelemetry/sdk/logs/logger.h"
#include "opentelemetry/sdk/logs/logger_provider_factory.h"
#include "opentelemetry/sdk/logs/batch_log_record_processor_factory.h"
#include "opentelemetry/sdk/logs/logger_context_factory.h"
#include "opentelemetry/sdk/resource/resource.h"
#include "opentelemetry/sdk/resource/semantic_conventions.h"

using namespace std;
namespace nostd     = opentelemetry::nostd;
namespace otlp      = opentelemetry::exporter::otlp;
namespace logs      = opentelemetry::logs;
namespace logs_sdk  = opentelemetry::sdk::logs;
namespace resource  = opentelemetry::sdk::resource;

namespace
{
  void initLogger() {
    otlp::OtlpHttpExporterOptions loggerOptions;
    auto exporter  = otlp::OtlpHttpExporterFactory::Create(loggerOptions);
    auto processor = logs_sdk::BatchLogRecordProcessorFactory::Create(std::move(exporter));
    std::vector<std::unique_ptr<logs_sdk::LogRecordProcessor>> processors;
    processors.push_back(std::move(processor));

    // 从环境变量获取 workspace
    const char* workspace_env = std::getenv("ACS_CMS_WORKSPACE");
    std::string workspace = workspace_env ? workspace_env : "default";

    // 创建带有默认 resource 的 LoggerContext
    resource::ResourceAttributes attributes = {
        {"acs_cms_workspace", workspace}
    };
    auto resource = resource::Resource::Create(attributes);

    auto context = logs_sdk::LoggerContextFactory::Create(std::move(processors), std::move(resource));
    std::shared_ptr<logs::LoggerProvider> provider = logs_sdk::LoggerProviderFactory::Create(std::move(context));
    opentelemetry::logs::Provider::SetLoggerProvider(provider);
  }

  nostd::shared_ptr<opentelemetry::logs::Logger> getLogger(std::string name){
    auto provider = logs::Provider::GetLoggerProvider();
    return provider->GetLogger(name + "_logger", name, OPENTELEMETRY_SDK_VERSION);
  }
}
