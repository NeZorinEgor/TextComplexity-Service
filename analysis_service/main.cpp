#include "router.hpp"
#include <grpc/grpc.h>
#include <grpc++/server.h>
#include <grpc++/server_builder.h>
#include <memory>

//188.168.25.28 21112
int main(int argc, char* argv[]) {
    if(argc!=3)
        return -1;
    RouterService router;
    grpc::ServerBuilder builder;
    builder.AddListeningPort(std::string(argv[1]) + ":" + std::string(argv[2]),grpc::InsecureServerCredentials());
    builder.RegisterService(&router); 
    std::unique_ptr<grpc::Server> server(builder.BuildAndStart());
    server->Wait();
    return 0;
}