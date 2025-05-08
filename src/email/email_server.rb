# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

require "ostruct"
require "pony"
require "sinatra"
require "logger"

require "opentelemetry/sdk"
require "opentelemetry/exporter/otlp"
require "opentelemetry/instrumentation/sinatra"

set :port, ENV["EMAIL_PORT"]

# 配置日志记录器
logger = Logger.new(STDOUT)
logger.level = Logger::INFO
logger.formatter = proc do |severity, datetime, progname, msg|
  "#{datetime.strftime('%Y-%m-%d %H:%M:%S')} [#{severity}] #{msg}\n"
end

OpenTelemetry::SDK.configure do |c|
  c.use "OpenTelemetry::Instrumentation::Sinatra"
end

post "/send_order_confirmation" do
  data = JSON.parse(request.body.read, object_class: OpenStruct)

  logger.info("收到订单确认邮件请求: 订单ID=#{data.order.order_id}, 收件人=#{data.email}")

  # get the current auto-instrumented span
  current_span = OpenTelemetry::Trace.current_span
  current_span.add_attributes({
    "app.order.id" => data.order.order_id,
  })

  begin
    send_email(data)
    logger.info("订单确认邮件发送成功: 订单ID=#{data.order.order_id}, 收件人=#{data.email}")
  rescue => e
    logger.error("订单确认邮件发送失败: 订单ID=#{data.order.order_id}, 收件人=#{data.email}, 错误=#{e.message}")
    raise e
  end
end

error do
  error = env['sinatra.error']
  logger.error("邮件服务发生错误: #{error.message}")
  OpenTelemetry::Trace.current_span.record_exception(error)
end

def send_email(data)
  # create and start a manual span
  tracer = OpenTelemetry.tracer_provider.tracer('email')
  tracer.in_span("send_email") do |span|
    logger.info("开始发送订单确认邮件: 订单ID=#{data.order.order_id}, 收件人=#{data.email}")
    
    Pony.mail(
      to:       data.email,
      from:     "noreply@example.com",
      subject:  "Your confirmation email",
      body:     erb(:confirmation, locals: { order: data.order }),
      via:      :test
    )
    
    span.set_attribute("app.email.recipient", data.email)
    logger.info("订单确认邮件发送完成: 订单ID=#{data.order.order_id}, 收件人=#{data.email}")
  end
  # manually created spans need to be ended
  # in Ruby, the method `in_span` ends it automatically
  # check out the OpenTelemetry Ruby docs at: 
  # https://opentelemetry.io/docs/instrumentation/ruby/manual/#creating-new-spans 
end
