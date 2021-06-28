#!/bin/sh

sqlite3 perfectward.db <<'END_SQL'
.timeout 2000
create table averages
(
uuid           INTEGER PRIMARY KEY,
questions      INTEGER NOT NULL,
positives      INTEGER NOT NULL
);
END_SQL