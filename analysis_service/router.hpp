#pragma once
#include <text_analys.grpc.pb.h>

using namespace Text_analys;

class RouterService final : public TextAnalysService::Service {
public:
    virtual ::grpc::Status getResult(::grpc::ServerContext* context,
        const ::Text_analys::SettingsTextPB* request,
        ::Text_analys::ResultParsingPB* response) override;
        
    virtual ~RouterService();
};
