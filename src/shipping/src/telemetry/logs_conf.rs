// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

use opentelemetry::{global, KeyValue};
use opentelemetry_sdk::logs::Config;
use opentelemetry_sdk::Resource;
use opentelemetry_otlp::WithExportConfig;
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt, Registry};
use tracing_opentelemetry::OpenTelemetryLayer;

pub fn init_logger() -> Result<(), Box<dyn std::error::Error>> {
    // 配置 OTLP HTTP 导出器
    let otlp_exporter = opentelemetry_otlp::new_exporter()
        .http()
        .with_endpoint("http://localhost:4318/v1/logs");

    // 创建日志配置
    let config = Config::default().with_resource(get_resource_attr());

    // 创建日志处理器
    let logger_provider = opentelemetry_sdk::logs::LoggerProvider::builder()
        .with_config(config)
        .with_batch_exporter(otlp_exporter, opentelemetry_sdk::runtime::Tokio)
        .build();

    // 设置全局日志处理器
    global::set_logger_provider(logger_provider);

    // 创建 tracing 层
    let layer = OpenTelemetryLayer::new(tracing_opentelemetry::layer().with_tracer(
        opentelemetry::global::tracer("shipping-service"),
    ));

    // 初始化 tracing 订阅者
    Registry::default()
        .with(layer)
        .try_init()?;

    Ok(())
}

fn get_resource_attr() -> Resource {
    use std::env;
    
    // 获取 workspace 环境变量
    let workspace = env::var("ACS_CMS_WORKSPACE").unwrap_or_else(|_| "default".to_string());
    
    // 创建带有 workspace 的 resource
    Resource::new(vec![
        KeyValue::new("acs_cms_workspace", workspace),
    ])
}
