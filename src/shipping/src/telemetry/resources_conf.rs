// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

use opentelemetry_sdk::{
    resource::{
        EnvResourceDetector, OsResourceDetector, ProcessResourceDetector, ResourceDetector,
        SdkProvidedResourceDetector, TelemetryResourceDetector,
    },
    Resource,
};
use std::time::Duration;
use std::env;

pub fn get_resource_attr() -> Resource {
    let os_resource = OsResourceDetector.detect(Duration::from_secs(0));
    let process_resource = ProcessResourceDetector.detect(Duration::from_secs(0));
    let sdk_resource = SdkProvidedResourceDetector.detect(Duration::from_secs(0));
    let env_resource = EnvResourceDetector::new().detect(Duration::from_secs(0));
    let telemetry_resource = TelemetryResourceDetector.detect(Duration::from_secs(0));

    // 获取 workspace 环境变量
    let workspace = env::var("ACS_CMS_WORKSPACE").unwrap_or_else(|_| "default".to_string());
    
    // 创建带有 workspace 的 resource
    let workspace_resource = Resource::new(vec![
        ("acs_cms_workspace".into(), workspace.into()),
    ]);

    os_resource
        .merge(&process_resource)
        .merge(&sdk_resource)
        .merge(&env_resource)
        .merge(&telemetry_resource)
        .merge(&workspace_resource)
}
