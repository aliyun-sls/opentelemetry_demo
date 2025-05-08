#!/usr/bin/python

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0
from opentelemetry.instrumentation.grpc import GrpcInstrumentorClient, GrpcInstrumentorServer
GrpcInstrumentorClient().instrument()
GrpcInstrumentorServer().instrument()
# Python
import os
import random
from concurrent import futures

# Pip
import grpc
from langchain.llms.openai import OpenAI
from langchain.chat_models import ChatOpenAI
from langchain.schema import HumanMessage, SystemMessage
from opentelemetry import trace, metrics
from opentelemetry._logs import set_logger_provider
from opentelemetry.exporter.otlp.proto.grpc._log_exporter import (
    OTLPLogExporter,
)
from opentelemetry.sdk._logs import LoggerProvider, LoggingHandler
from opentelemetry.sdk._logs.export import BatchLogRecordProcessor
from opentelemetry.sdk.resources import Resource

from openfeature import api
from openfeature.contrib.provider.flagd import FlagdProvider

from openfeature.contrib.hook.opentelemetry import TracingHook

# Local
import logging
import demo_pb2
import demo_pb2_grpc
from grpc_health.v1 import health_pb2
from grpc_health.v1 import health_pb2_grpc

from metrics import (
    init_metrics
)

cached_ids = []
first_run = True

class RecommendationService(demo_pb2_grpc.RecommendationServiceServicer):
    def ListRecommendations(self, request, context):
        prod_list = get_product_list(request.product_ids)
        span = trace.get_current_span()
        span.set_attribute("app.products_recommended.count", len(prod_list))
        logger.info(f"Receive ListRecommendations for product ids:{prod_list}")

        # build and return response
        response = demo_pb2.ListRecommendationsResponse()
        response.product_ids.extend(prod_list)

        # Collect metrics for this service
        rec_svc_metrics["app_recommendations_counter"].add(len(prod_list), {'recommendation.type': 'catalog'})

        return response

    def Check(self, request, context):
        return health_pb2.HealthCheckResponse(
            status=health_pb2.HealthCheckResponse.SERVING)

    def Watch(self, request, context):
        return health_pb2.HealthCheckResponse(
            status=health_pb2.HealthCheckResponse.UNIMPLEMENTED)


def get_model_config():
    """从flagd获取模型配置"""
    client = api.get_client()
    try:
        model_config = client.get_object_value("recommendationModelConfig", {
            "model": "qwen-plus",
            "temperature": 0.7,
            "max_tokens": 100,
            "ai_probability": 0.3
        })
        return model_config
    except Exception as e:
        logger.warning(f"无法从flagd获取模型配置，使用默认配置: {str(e)}")
        return {
            "model": "qwen-plus",
            "temperature": 0.7,
            "max_tokens": 100,
            "ai_probability": 0.3
        }

def get_product_list(request_product_ids):
    global first_run
    global cached_ids
    with tracer.start_as_current_span("get_product_list") as span:
        max_responses = 5

        # Formulate the list of characters to list of strings
        request_product_ids_str = ''.join(request_product_ids)
        request_product_ids = request_product_ids_str.split(',')

        # 获取所有商品信息
        cat_response = product_catalog_stub.ListProducts(demo_pb2.Empty())
        all_products = cat_response.products
        
        # 获取请求的商品详情
        requested_products = []
        for product_id in request_product_ids:
            product = next((p for p in all_products if p.id == product_id), None)
            if product:
                requested_products.append(product)
            else:
                logger.warning(f"未找到商品ID: {product_id}")

        # 获取模型配置
        model_config = get_model_config()
        span.set_attribute("app.recommendation.model", model_config["model"])
        span.set_attribute("app.recommendation.ai_probability", model_config.get("ai_probability", 0.3))

        # 根据概率决定是否使用AI推荐
        if random.random() > model_config.get("ai_probability", 0.3):
            logger.info("使用随机推荐")
            filtered_products = list(set([p.id for p in all_products]) - set(request_product_ids))
            num_products = len(filtered_products)
            num_return = min(max_responses, num_products)
            indices = random.sample(range(num_products), num_return)
            return [filtered_products[i] for i in indices]

        # 构建提示词
        if requested_products:
            prompt = f"""基于以下用户正在查看的商品，推荐5个相关的商品：

用户正在查看的商品：
{', '.join([f"ID:{p.id}, 名称:{p.name}, 描述:{p.description}" for p in requested_products])}

可选的商品列表：
{', '.join([f"ID:{p.id}, 名称:{p.name}, 描述:{p.description}" for p in all_products if p.id not in request_product_ids])}

请只返回商品ID列表，用逗号分隔。"""
        else:
            prompt = f"""请从以下商品列表中随机推荐5个商品：

可选的商品列表：
{', '.join([f"ID:{p.id}, 名称:{p.name}, 描述:{p.description}" for p in all_products])}

请只返回商品ID列表，用逗号分隔。"""

        try:
            # 调用AI API
            client = ChatOpenAI(
                api_key=os.getenv('OPENAI_API_KEY'),
                base_url=os.getenv('OPENAI_BASE_URL', 'https://dashscope.aliyuncs.com/api/v1'),
                model=model_config["model"],
                temperature=model_config.get("temperature", 0.7),
                max_tokens=model_config.get("max_tokens", 100)
            )
            
            messages = [
                SystemMessage(content="你是一个专业的商品推荐助手。请根据用户正在查看的商品，推荐最相关的商品。只返回商品ID列表，用逗号分隔。"),
                HumanMessage(content=prompt)
            ]
            
            response = client.invoke(messages)
            logger.info(f"AI API调用结果: {response}")
            
            # 解析推荐结果
            recommended_ids = response.content.strip().split(',')
            recommended_ids = [id.strip() for id in recommended_ids if id.strip()]
            
            # 确保返回的商品数量不超过max_responses
            prod_list = recommended_ids[:max_responses]
            
        except Exception as e:
            logger.error(f"AI API调用失败: {str(e)}")
            # 如果API调用失败，回退到随机推荐
            filtered_products = list(set([p.id for p in all_products]) - set(request_product_ids))
            num_products = len(filtered_products)
            num_return = min(max_responses, num_products)
            indices = random.sample(range(num_products), num_return)
            prod_list = [filtered_products[i] for i in indices]

        span.set_attribute("app.filtered_products.list", prod_list)
        return prod_list


def must_map_env(key: str):
    value = os.environ.get(key)
    if value is None:
        raise Exception(f'{key} environment variable must be set')
    return value


def check_feature_flag(flag_name: str):
    # Initialize OpenFeature
    client = api.get_client()
    return client.get_boolean_value("recommendationCacheFailure", False)


if __name__ == "__main__":
    service_name = must_map_env('OTEL_SERVICE_NAME')
    api.set_provider(FlagdProvider(host=os.environ.get('FLAGD_HOST', 'flagd'), port=os.environ.get('FLAGD_PORT', 8013)))
    api.add_hooks([TracingHook()])

    # Initialize Traces and Metrics
    tracer = trace.get_tracer_provider().get_tracer(service_name)
    meter = metrics.get_meter_provider().get_meter(service_name)
    rec_svc_metrics = init_metrics(meter)

    # Initialize Logs
    logger_provider = LoggerProvider(
        resource=Resource.create(
            {
                'service.name': service_name,
            }
        ),
    )
    set_logger_provider(logger_provider)
    log_exporter = OTLPLogExporter(insecure=True)
    logger_provider.add_log_record_processor(BatchLogRecordProcessor(log_exporter))
    handler = LoggingHandler(level=logging.NOTSET, logger_provider=logger_provider)

    # Attach OTLP handler to logger
    logger = logging.getLogger('main')
    logger.addHandler(handler)

    catalog_addr = must_map_env('PRODUCT_CATALOG_ADDR')
    pc_channel = grpc.insecure_channel(catalog_addr)
    product_catalog_stub = demo_pb2_grpc.ProductCatalogServiceStub(pc_channel)

    # Create gRPC server
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    # Add class to gRPC server
    service = RecommendationService()
    demo_pb2_grpc.add_RecommendationServiceServicer_to_server(service, server)
    health_pb2_grpc.add_HealthServicer_to_server(service, server)

    # Start server
    port = must_map_env('RECOMMENDATION_PORT')
    server.add_insecure_port(f'[::]:{port}')
    server.start()
    logger.info(f'Recommendation service started, listening on port {port}')
    server.wait_for_termination()
