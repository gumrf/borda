-- Set timezone
SET TIMEZONE="Europe/Moscow";

-- Create users table
CREATE TABLE IF NOT EXISTS "user" (
   id serial NOT NULL PRIMARY KEY,
   name varchar(128) NOT NULL,
   password varchar(256) NOT NULL,
   contact varchar(256) UNIQUE NULL
);

-- Create roles table
CREATE TABLE IF NOT EXISTS role (
   id serial NOT NULL PRIMARY KEY,
   name varchar(256)
);

-- Create books table
CREATE TABLE IF NOT EXISTS user_role (
   user_id integer NOT NULL PRIMARY KEY REFERENCES "user"(id),
   role_id integer NOT NULL REFERENCES role(id)
);

CREATE TABLE IF NOT EXISTS team (
   id serial NOT NULL PRIMARY KEY,
   name varchar(256) UNIQUE NOT NULL,
   token varchar(256) UNIQUE NOT NULL,
   team_leader_id integer NOT NULL REFERENCES "user"(id)
);

CREATE TABLE IF NOT EXISTS team_member (
   id serial NOT NULL PRIMARY KEY,
   team_id integer NOT NULL REFERENCES team(id),
   user_id integer NOT NULL REFERENCES "user"(id)
);

CREATE TABLE IF NOT EXISTS author (
   id serial NOT NULL PRIMARY KEY,
   name varchar(256) NOT NULL,
   contact varchar(256),
   UNIQUE(name)
);

CREATE TABLE IF NOT EXISTS task (
   id          serial PRIMARY KEY,
   title       VARCHAR(256) NOT NULL,
   description TEXT NOT NULL,
   category    VARCHAR(256) NOT NULL,
   complexity  VARCHAR(256),
   points      INTEGER NOT NULL,
   hint        TEXT,
   flag        VARCHAR(256) NOT NULL,
   is_active   bool NOT NULL,
   is_disabled bool NOT NULL,
   author_id   INTEGER REFERENCES author(id)
);

CREATE TABLE IF NOT EXISTS solved_task (
   task_id integer NOT NULL REFERENCES task(id),
   team_id integer NOT NULL REFERENCES team(id),
   timestamp timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS task_submission (
   task_id integer NOT NULL REFERENCES task(id),
   team_id integer NOT NULL REFERENCES team(id),
   user_id integer NOT NULL REFERENCES "user"(id),
   flag varchar(256) NOT NULL,
   is_correct bool NOT NULL,
   timestamp timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS settings (
   id serial NOT NULL PRIMARY KEY,
   key varchar(256) UNIQUE NOT NULL,
   value varchar(256) NOT NULL,
   timestamp timestamptz NOT NULL DEFAULT now()
);