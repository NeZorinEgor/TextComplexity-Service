#include "router.hpp"
#include <grpc/grpc.h>
#include <grpc++/server.h>
#include <grpc++/server_builder.h>
#include <memory>

//188.168.25.28 21112
int main() {
    
    RouterService router;
    grpc::ServerBuilder builder;
    builder.AddListeningPort("127.0.0.1:1111",grpc::InsecureServerCredentials());
    builder.RegisterService(&router); 
    std::unique_ptr<grpc::Server> server(builder.BuildAndStart());
    server->Wait();
    return 0;
}