#include "router.hpp"
#include "analys_alg.h"

::grpc::Status RouterService::getResult(::grpc::ServerContext* context,
        const ::Text_analys::SettingsTextPB* request,
        ::Text_analys::ResultParsingPB* response) {
    auto result = aal::analys(request->text());
    
    response->set_hard_reading(result.hard);
    response->set_mood(static_cast<ResultParsingPB_Mood>(result.mood));
    response->set_water_value(result.water);

    return ::grpc::Status::OK;
}

RouterService::~RouterService() {}