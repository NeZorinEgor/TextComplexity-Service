#include "router.hpp"

::grpc::Status getResult(::grpc::ServerContext* context,
        const ::Text_analys::SettingsTextPB* request,
        ::Text_analys::ResultParsingPB* response) {
    if(request->text() == "test")
        response->set_mood(ResultParsingPB_Mood_happy);
        //context->
    return ::grpc::Status::OK;
}

RouterService::~RouterService() {}