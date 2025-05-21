// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
using System;

using cart.cartstore;
using cart.services;

using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Http;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Diagnostics.HealthChecks;
using Microsoft.Extensions.Logging;
using OpenTelemetry.Instrumentation.StackExchangeRedis;
using OpenTelemetry.Logs;
using OpenTelemetry.Metrics;
using OpenTelemetry.Resources;
using OpenTelemetry.Trace;
using OpenFeature;
using OpenFeature.Contrib.Providers.Flagd;
using OpenFeature.Contrib.Hooks.Otel;

var builder = WebApplication.CreateBuilder(args);

// OpenTelemetry 初始化
builder.Services.AddOpenTelemetry()
    .WithTracing(tracerProviderBuilder =>
        tracerProviderBuilder
            .AddSource(DiagnosticsConfig.ActivitySource.Name)
            .SetResourceBuilder(OpenTelemetry.Resources.ResourceBuilder.CreateDefault()
                .AddAttributes(new Dictionary<string, object> {
                    {"service.name", DiagnosticsConfig.ServiceName},
                    {"host.name",DiagnosticsConfig.HostName},
                    {"acs_cms_workspace","o11y-demo-cn-heyuan"}
                }))
            .AddAspNetCoreInstrumentation()
            .AddConsoleExporter() // 在控制台输出Trace数据，可选
            .AddOtlpExporter(opt =>
            {
                // 使用HTTP协议上报
                opt.Endpoint = new Uri(DiagnosticsConfig.Endpoint);
                opt.Headers = DiagnosticsConfig.LicenseKey;
                opt.Protocol = OtlpExportProtocol.HttpProtobuf;

                // 使用gRPC协议上报
                // opt.Endpoint = new Uri("<grpc_endpoint>");
                // opt.Headers = "Authentication=<token>";
                // opt.Protocol = OtlpExportProtocol.Grpc;
            })
     );

string valkeyAddress = builder.Configuration["VALKEY_ADDR"];
if (string.IsNullOrEmpty(valkeyAddress))
{
    Console.WriteLine("VALKEY_ADDR environment variable is required.");
    Environment.Exit(1);
}

builder.Logging
    .AddOpenTelemetry(options => options.AddOtlpExporter())
    .AddConsole();

builder.Services.AddSingleton<ICartStore>(x=>
{
    var store = new ValkeyCartStore(x.GetRequiredService<ILogger<ValkeyCartStore>>(), valkeyAddress);
    store.Initialize();
    return store;
});

builder.Services.AddSingleton<IFeatureClient>(x => {
    var flagdProvider = new FlagdProvider();
    Api.Instance.SetProviderAsync(flagdProvider).GetAwaiter().GetResult();
    var client = Api.Instance.GetClient();
    return client;
});

builder.Services.AddSingleton(x =>
    new CartService(
        x.GetRequiredService<ICartStore>(),
        new ValkeyCartStore(x.GetRequiredService<ILogger<ValkeyCartStore>>(), "badhost:1234"),
        x.GetRequiredService<IFeatureClient>()
));


Action<ResourceBuilder> appResourceBuilder =
    resource => resource
        .AddContainerDetector()
        .AddHostDetector();

builder.Services.AddOpenTelemetry()
    .ConfigureResource(appResourceBuilder)
    .WithTracing(tracerBuilder => tracerBuilder
        .AddSource("OpenTelemetry.Demo.Cart")
        .AddRedisInstrumentation(
            options => options.SetVerboseDatabaseStatements = true)
        .AddAspNetCoreInstrumentation()
        .AddGrpcClientInstrumentation()
        .AddHttpClientInstrumentation()
        .AddOtlpExporter())
    .WithMetrics(meterBuilder => meterBuilder
        .AddMeter("OpenTelemetry.Demo.Cart")
        .AddProcessInstrumentation()
        .AddRuntimeInstrumentation()
        .AddAspNetCoreInstrumentation()
        .SetExemplarFilter(ExemplarFilterType.TraceBased)
        .AddOtlpExporter());
OpenFeature.Api.Instance.AddHooks(new TracingHook());
builder.Services.AddGrpc();
builder.Services.AddGrpcHealthChecks()
    .AddCheck("Sample", () => HealthCheckResult.Healthy());

var app = builder.Build();

var ValkeyCartStore = (ValkeyCartStore) app.Services.GetRequiredService<ICartStore>();
app.Services.GetRequiredService<StackExchangeRedisInstrumentation>().AddConnection(ValkeyCartStore.GetConnection());

app.MapGrpcService<CartService>();
app.MapGrpcHealthChecksService();

app.MapGet("/", async context =>
{
    await context.Response.WriteAsync("Communication with gRPC endpoints must be made through a gRPC client. To learn how to create a client, visit: https://go.microsoft.com/fwlink/?linkid=2086909");
});

app.Run();

// 创建DiagnosticsConfig类
public static class DiagnosticsConfig
{
    public string ServiceName = "<your-service-name>"; // your service name
    public string HostName = "<your-host-name>"; // your host name
    public string Endpoint = "<your-host-name>"; // your host name
    public string LicenseKey = "<your-license-key>";
    ServiceName = Environment.GetEnvironmentVariable("OTEL_SERVICE_NAME")
    HostName = Environment.GetEnvironmentVariable("HostName")
    Endpoint = Environment.GetEnvironmentVariable("OTEL_EXPORTER_OTLP_ENDPOINT")
    LicenseKey = Environment.GetEnvironmentVariable("OTEL_EXPORTER_OTLP_HEADERS")
    public static ActivitySource ActivitySource = new ActivitySource(ServiceName);
}
