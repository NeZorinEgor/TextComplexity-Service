cmake_minimum_required(VERSION 2.8)

project(Analisys_service)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

set(ABSL_ENABLE_INSTALL ON)

add_subdirectory(lib/SQLiteCpp)
add_subdirectory(lib/grpc)

file(GLOB SRC
     "*.h"
     "*.cpp"
)

list(APPEND CMAKE_MODULE_PATH "${CMAKE_CURRENT_LIST_DIR}/lib/grpc/cmake")

file(GLOB PROTOBUF_MODELS_SOURCES "protos/generate/cpp/*.pb.cc")

add_executable(analysis_service ${SRC} ${PROTOBUF_MODELS_SOURCES})

include_directories(protos/generate/cpp lib/SQLiteCpp/include)

target_link_libraries(analysis_service grpc++ SQLiteCpp sqlite3)