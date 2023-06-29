#pragma once

#include <sqlite3.h>
#include <SQLiteCpp/SQLiteCpp.h>
#include "constants.hpp"
#include <optional>

class AppDB {
public:
    static AppDB & Instanse();
    std::optional<Result> getResult(int32_t hash);
    void addResult(int32_t hash, const Result &);

private:
    AppDB();
    SQLite::Database m_db;
};