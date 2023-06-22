#include "router.hpp"
#include <grpc/grpc.h>
#include <grpc++/server.h>
#include <grpc++/server_builder.h>
#include <memory>

//188.168.25.28 21112
int main() {
    
    RouterService router;
    grpc::ServerBuilder builder;
    builder.AddListeningPort("192.168.0.200:21112",grpc::InsecureServerCredentials());
    builder.RegisterService(&router); 
    std::unique_ptr<grpc::Server> server(builder.BuildAndStart());
    server->Wait();
    return 0;
}