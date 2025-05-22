// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

use opentelemetry::{global, KeyValue};
use opentelemetry_sdk::logs::Config;
use opentelemetry_sdk::Resource;
use tracing_subscriber::{layer::SubscriberExt, Registry};
use std::env;
use opentelemetry::trace::TracerProvider;

pub fn init_logger() -> Result<(), Box<dyn std::error::Error>> {
    // 通过环境变量获取 endpoint
    let endpoint = env::var("OTEL_EXPORTER_OTLP_ENDPOINT")
        .unwrap_or_else(|_| "http://localhost:4318".to_string());

    // 配置 OTLP HTTP 导出器
    let otlp_exporter = opentelemetry_otlp::new_exporter()
        .with_endpoint(format!("{}/v1/logs", endpoint));

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
    let tracer_provider = opentelemetry_sdk::trace::TracerProvider::builder()
        .with_batch_exporter(
            opentelemetry_otlp::new_exporter()
                .with_endpoint(format!("{}/v1/traces", endpoint)),
            opentelemetry_sdk::runtime::Tokio,
        )
        .build();
    let tracer = tracer_provider.tracer("shipping-service");
    let telemetry = tracing_opentelemetry::layer().with_tracer(tracer);

    // 初始化 tracing 订阅者
    let subscriber = Registry::default().with(telemetry);
    tracing::subscriber::set_global_default(subscriber)?;

    Ok(())
}

fn get_resource_attr() -> Resource {
    // 获取 workspace 环境变量
    let workspace = env::var("ACS_CMS_WORKSPACE").unwrap_or_else(|_| "default".to_string());
    
    // 创建带有 workspace 的 resource
    Resource::new(vec![
        KeyValue::new("acs_cms_workspace", workspace),
    ])
}
