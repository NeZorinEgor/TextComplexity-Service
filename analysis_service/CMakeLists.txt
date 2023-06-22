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

add_executable(analysis_service main.cpp router.cpp ${PROTOBUF_MODELS_SOURCES})

include_directories(protos/generate/cpp)

target_link_libraries(analysis_service grpc++)