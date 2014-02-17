GolangWebServer
===============

Web Server from scratch in golang

```bash
export POSTGRESQL_USER={{username}}
export POSTGRESQL_PASS={{password}}
export PGHOST={{localhost}}

psql -d {{dbname}} -U $POSTGRESQL_USER -W $POSTGRESQL_PASS
```

```SQL
 CREATE SEQUENCE users_id_sequence;
 CREATE TABLE users (
    _id integer PRIMARY KEY DEFAULT nextval('users_id_sequence'),
    username text NOT NULL UNIQUE,
    password_hash text NOT NULL,
    created_at timestamp NOT NULL
);
 ALTER SEQUENCE users_id_sequence OWNED BY users._id;

 CREATE SEQUENCE users_files_id_sequence;
 CREATE TABLE users_files (
    _id integer PRIMARY KEY DEFAULT nextval('users_files_id_sequence'),
    title text NOT NULL,
    path text NOT NULL,
    user_id integer REFERENCES users (_id),
    size bigint NOT NULL
 );
 ALTER SEQUENCE users_files_id_sequence OWNED BY users_files._id;

 CREATE SEQUENCE tags_id_sequence;
 CREATE TABLE tags (
    _id integer PRIMARY KEY DEFAULT nextval('tags_id_sequence'),
    name text NOT NULL UNIQUE
);
 ALTER SEQUENCE tags_id_sequence OWNED BY tags._id;

  CREATE SEQUENCE file_tags_id_sequence;
  CREATE TABLE file_tags (
     _id integer PRIMARY KEY DEFAULT nextval('file_tags_id_sequence'),
     file_id integer REFERENCES users_files (_id),
     tag_id integer REFERENCES tags (_id)
 );
  ALTER SEQUENCE file_tags_id_sequence OWNED BY file_tags._id;
```

```bash
git clone https://github.com/sescobb27/GolangWebServer.git

cd GolangWebServer

make

./start
```
