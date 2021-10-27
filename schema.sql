
CREATE TABLE  IF NOT EXISTS Users(
   ID INT PRIMARY KEY     NOT NULL,
   NAME           TEXT    NOT NULL,
   PASSWORD       TEXT    NOT NULL,
   EMAIL          TEXT    NOT NULL,
   PHONE          TEXT,
   AGE            INT,
   ADDRESS        TEXT 
);
