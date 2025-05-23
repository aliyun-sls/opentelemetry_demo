// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

#include "opentelemetry/exporters/otlp/otlp_http_exporter_factory.h"
#include "opentelemetry/logs/provider.h"
#include "opentelemetry/sdk/logs/logger.h"
#include "opentelemetry/sdk/logs/logger_provider_factory.h"
#include "opentelemetry/sdk/logs/batch_log_record_processor_factory.h"
#include "opentelemetry/sdk/logs/logger_context_factory.h"
#include "opentelemetry/exporters/otlp/otlp_http_log_record_exporter_factory.h"
#include "opentelemetry/sdk/resource/resource.h"
#include "opentelemetry/sdk/resource/semantic_conventions.h"

using namespace std;
namespace nostd     = opentelemetry::nostd;
namespace otlp      = opentelemetry::exporter::otlp;
namespace logs      = opentelemetry::logs;
namespace logs_sdk  = opentelemetry::sdk::logs;

namespace
{
  void initLogger() {
    // 获取环境变量中的 workspace
    const char* workspace = std::getenv("ACS_CMS_WORKSPACE");
    std::string workspace_str = workspace ? workspace : "default";

    // 创建默认的 resource
    opentelemetry::sdk::resource::ResourceAttributes attributes;
    attributes["acs.cms.workspace"] = workspace_str;
    auto resource = opentelemetry::sdk::resource::Resource::Create(attributes);

    otlp::OtlpHttpLogRecordExporterOptions loggerOptions;
    auto exporter  = otlp::OtlpHttpLogRecordExporterFactory::Create(loggerOptions);

    // 配置 BatchLogRecordProcessor 选项
    logs_sdk::BatchLogRecordProcessorOptions batch_options;
    batch_options.max_queue_size = 2048;  // 最大队列大小
    batch_options.schedule_delay_millis = std::chrono::milliseconds(5000);  // 调度延迟
    batch_options.max_export_batch_size = 512;  // 最大导出批次大小

    auto processor = logs_sdk::BatchLogRecordProcessorFactory::Create(
        std::move(exporter), batch_options);
    std::vector<std::unique_ptr<logs_sdk::LogRecordProcessor>> processors;
    processors.push_back(std::move(processor));
    auto context = logs_sdk::LoggerContextFactory::Create(std::move(processors), std::move(resource));
    std::shared_ptr<logs::LoggerProvider> provider = logs_sdk::LoggerProviderFactory::Create(std::move(context));
    opentelemetry::logs::Provider::SetLoggerProvider(provider);
  }

  nostd::shared_ptr<opentelemetry::logs::Logger> getLogger(std::string name){
    auto provider = logs::Provider::GetLoggerProvider();
    return provider->GetLogger(name + "_logger", name, OPENTELEMETRY_SDK_VERSION);
  }
}
