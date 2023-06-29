#include "db.h"
#include <iostream>

static const std::string kDBPath = "./analysis_local_db.sqlite";

AppDB &AppDB::Instanse() {
    static AppDB inst = AppDB();
    return inst;
}

std::optional<Result> AppDB::getResult(int32_t hash) {
try
{    
    SQLite::Statement query(m_db, "SELECT * FROM analysis_hash WHERE hash = ?");
    query.bind(1,hash);
    while(query.executeStep()) {
        Result out;
        out.water = query.getColumn(1);
        out.mood = static_cast<Result::TypeMood>(int(query.getColumn(2)));
        out.hard = query.getColumn(3);
        return out;
    }
    return std::nullopt;
} catch (std::exception& e) {
    std::cout << "DB ERROR: " << e.what() << std::endl;
}   
    return std::nullopt;
}

void AppDB::addResult(int32_t hash, const Result &result) {
try {
    SQLite::Transaction transaction(m_db);
    m_db.exec(std::string("INSERT INTO analysis_hash VALUES (") 
             + std::to_string(hash) + std::string(",")
             + std::to_string(result.water) + std::string(",")
             + std::to_string(static_cast<int>(result.mood)) + std::string(",")
             + std::to_string(result.hard) 
             + ")");
    transaction.commit();
} catch(std::exception& e) {
    std::cout << "DB ERROR: " << e.what() <<std::endl;
}
}

AppDB::AppDB()
: m_db(kDBPath, SQLite::OPEN_READWRITE|SQLite::OPEN_CREATE) {
try {
    SQLite::Transaction transaction(m_db);
    m_db.exec("CREATE TABLE IF NOT EXISTS analysis_hash (hash INTEGER PRIMARY KEY, water INTEGER, mood INTEGER, harder INTEGER)");
    transaction.commit();
} catch(std::exception& e) {
    std::cout << "DB ERROR: " << e.what() <<std::endl;
}
}