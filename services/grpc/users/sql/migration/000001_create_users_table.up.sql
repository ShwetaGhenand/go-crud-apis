BEGIN;

CREATE TABLE  IF NOT EXISTS users(
   id             INT   PRIMARY KEY     NOT NULL,
   name           TEXT                  NOT NULL,
   password       TEXT                  NOT NULL,
   email          TEXT                  NOT NULL,
   phone          TEXT,
   age            INT,
   address        TEXT,
   createTime     TIMESTAMP                NOT NULL,
   updateTime     TIMESTAMP                NOT NULL 
);

COMMIT;
