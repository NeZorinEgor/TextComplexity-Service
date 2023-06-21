cmake_minimum_required(VERSION 2.8)

file(GLOB SRC
     "*.h"
     "*.cpp"
)


add_executable(analysis_service ${SRC})