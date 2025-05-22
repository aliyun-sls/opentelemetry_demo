// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

#pragma once
#include "opentelemetry/exporters/otlp/otlp_http_exporter_factory.h"
#include "opentelemetry/exporters/otlp/otlp_http_exporter_options.h"
#include "opentelemetry/metrics/provider.h"
#include "opentelemetry/sdk/metrics/export/periodic_exporting_metric_reader.h"
#include "opentelemetry/sdk/metrics/meter.h"
#include "opentelemetry/sdk/metrics/meter_provider.h"
#include "opentelemetry/sdk/resource/resource.h"
#include "opentelemetry/sdk/resource/semantic_conventions.h"

// namespaces
namespace common        = opentelemetry::common;
namespace metrics_api   = opentelemetry::metrics;
namespace metric_sdk    = opentelemetry::sdk::metrics;
namespace nostd         = opentelemetry::nostd;
namespace otlp_exporter = opentelemetry::exporter::otlp;
namespace resource      = opentelemetry::sdk::resource;

namespace
{
  void initMeter() 
  {
    // Build MetricExporter
    otlp_exporter::OtlpHttpExporterOptions otlpOptions;
    auto exporter = otlp_exporter::OtlpHttpExporterFactory::Create(otlpOptions);

    // 从环境变量获取 workspace
    const char* workspace_env = std::getenv("ACS_CMS_WORKSPACE");
    std::string workspace = workspace_env ? workspace_env : "default";

    // 创建带有默认 resource 的 MeterProvider
    resource::ResourceAttributes attributes = {
        {"acs_cms_workspace", workspace}
    };
    auto resource = resource::Resource::Create(attributes);

    // Build MeterProvider and Reader
    metric_sdk::PeriodicExportingMetricReaderOptions options;
    std::unique_ptr<metric_sdk::MetricReader> reader{
        new metric_sdk::PeriodicExportingMetricReader(std::move(exporter), options) };
    auto provider = std::shared_ptr<metrics_api::MeterProvider>(new metric_sdk::MeterProvider(resource));
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
