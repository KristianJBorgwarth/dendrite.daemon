CREATE TABLE IF NOT EXISTS index_build_flag (
    lock CHAR(1) PRIMARY KEY DEFAULT 'X' CHECK (lock = 'X'),
    rebuilding BOOLEAN NOT NULL DEFAULT 0,
);
