#include "router.hpp"

::grpc::Status RouterService::getResult(::grpc::ServerContext* context,
        const ::Text_analys::SettingsTextPB* request,
        ::Text_analys::ResultParsingPB* response) {
    if(request->text() == "test")
        response->set_mood(ResultParsingPB_Mood_happy);
        
    response->set_mood(ResultParsingPB_Mood_boring);
    response->set_water_value(7);
    response->set_hard_reading(2);

    return ::grpc::Status::OK;
}

RouterService::~RouterService() {}