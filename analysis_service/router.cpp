#include "router.hpp"
#include "analys_alg.h"
#include "db.h"
#include <iostream>

static int32_t crc32(size_t len, const void* data, const unsigned int POLY = 0x04C11DB7) {
    const unsigned char* buffer = (const unsigned char*)data;
    unsigned int crc = -1;

    while (len--)
    {
        crc = crc ^ (*buffer++ << 24);
        for (int bit = 0; bit < 8; bit++)
        {
            if (crc & (1L << 31)) crc = (crc << 1) ^ POLY;
            else                  crc = (crc << 1);
        }
    }
    uint32_t crc_u = ~crc;
    int32_t crc_out = *reinterpret_cast<int*>(&crc);

    return std::abs(crc_out);
}

::grpc::Status RouterService::getResult(::grpc::ServerContext *context,
                                        const ::Text_analys::SettingsTextPB *request,
                                        ::Text_analys::ResultParsingPB *response)
{
    Result result;
//work with hashes
    int32_t hash = crc32(request->text().size(), request->text().c_str());
    auto opt = AppDB::Instanse().getResult(hash); //TODO CRC32 
    if(opt.has_value()) {
        result = opt.value();
std::cout << "Db use hashes\n";
    }
    else {
        result = aal::analys(request->text()); 
        AppDB::Instanse().addResult(hash,result); //TODO CRC32
    }
//set settings
    response->set_hard_reading(result.hard);
    response->set_mood(static_cast<ResultParsingPB_Mood>(result.mood));
    response->set_water_value(result.water);

    return ::grpc::Status::OK;
}

RouterService::~RouterService() {}