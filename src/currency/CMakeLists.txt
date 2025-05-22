cmake_minimum_required(VERSION 3.1)
project(currency)

# 设置 CURL 的查找路径
set(CURL_ROOT_DIR "/usr" CACHE PATH "CURL root directory")
set(CURL_INCLUDE_DIR "/usr/include" CACHE PATH "CURL include directory")
set(CURL_LIBRARY "/usr/lib/libcurl.so" CACHE FILEPATH "CURL library")

find_package(Protobuf REQUIRED)
find_package(gRPC CONFIG REQUIRED)
find_package(opentelemetry-cpp CONFIG REQUIRED)
find_package(CURL REQUIRED)
find_package(nlohmann_json)

set(GENERATED_PROTOBUF_PATH "${CMAKE_BINARY_DIR}/generated/proto")

add_subdirectory(genproto)

include_directories("${GENERATED_PROTOBUF_PATH}" "${OPENTELEMETRY_CPP_INCLUDE_DIRS}")

add_executable(currency src/server.cpp)
add_dependencies(currency demo-proto)

# 设置 OpenTelemetry 的库
set(OTEL_LIBS
    opentelemetry_trace
    opentelemetry_common
    opentelemetry_exporter_otlp_grpc
    opentelemetry_exporter_otlp_grpc_client
    opentelemetry_proto
    opentelemetry_otlp_recordable
    opentelemetry_resources
    opentelemetry_metrics
    opentelemetry_exporter_otlp_grpc_metrics
    opentelemetry_http_client_curl
)

target_link_libraries(
    currency 
    demo-proto 
    protobuf::libprotobuf
    ${OPENTELEMETRY_CPP_LIBRARIES}
    ${OTEL_LIBS}
    gRPC::grpc++
    CURL::libcurl
)

install(TARGETS currency DESTINATION bin)
