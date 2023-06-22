cmake_minimum_required(VERSION 2.8)

project(Analisys_service)

set(ABSL_ENABLE_INSTALL ON)

file(GLOB SRC
     "*.h"
     "*.cpp"
)

list(APPEND CMAKE_MODULE_PATH "${CMAKE_CURRENT_LIST_DIR}/lib/grpc/cmake")
add_subdirectory(lib/grpc)

file(GLOB PROTOBUF_MODELS_SOURCES "protos/generate/cpp/*.pb.cc")
file(GLOB PROTOBUF_MODELS_INCLUDES "protos/generate/cpp/*.hpp")

include_directories(${PROTOBUF_MODELS_INCLUDES})

add_executable(analysis_service ${SRC} ${PROTOBUF_MODELS_SOURCES})
target_include_directories(analysis_service PUBLIC ${PROTOBUF_MODELS_INCLUDES})

target_link_libraries(analysis_service grpc++)