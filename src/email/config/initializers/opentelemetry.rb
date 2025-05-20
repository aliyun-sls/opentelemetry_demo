# config/initializers/opentelemetry.rb
require 'opentelemetry/sdk'
require 'opentelemetry/exporter/otlp'
require 'opentelemetry/instrumentation/all'

endpoint = ENV['OTEL_EXPORTER_OTLP_TRACES_ENDPOINT']
serviceName = ENV['OTEL_SERVICE_NAME']
hostName = ENV['POD_HOSTNAME']
OpenTelemetry::SDK.configure do |c|
  c.add_span_processor(
    OpenTelemetry::SDK::Trace::Export::BatchSpanProcessor.new(
      OpenTelemetry::Exporter::OTLP::Exporter.new(
        endpoint: endpoint # HTTP方式接入
      )
    )
  )
  c.resource = OpenTelemetry::SDK::Resources::Resource.create({
    OpenTelemetry::SemanticConventions::Resource::HOST_NAME => hostName, # 主机名
  })
  c.service_name = serviceName    # 服务名
  c.use_all()    # 自动观测opentelemetry支持的所有库，
end