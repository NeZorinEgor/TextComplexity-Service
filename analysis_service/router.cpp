#include "router.hpp"
#include "analys_alg.h"
#include "db.h"

::grpc::Status RouterService::getResult(::grpc::ServerContext *context,
                                        const ::Text_analys::SettingsTextPB *request,
                                        ::Text_analys::ResultParsingPB *response)
{
    Result result;
//work with hashes
    auto opt = AppDB::Instanse().getResult(0); //TODO CRC32 
    if(opt.has_value())
        result = opt.value();
    else {
        result = aal::analys(request->text()); 
        AppDB::Instanse().addResult(0,result); //TODO CRC32
    }
//set settings
    response->set_hard_reading(result.hard);
    response->set_mood(static_cast<ResultParsingPB_Mood>(result.mood));
    response->set_water_value(result.water);

    return ::grpc::Status::OK;
}

RouterService::~RouterService() {}