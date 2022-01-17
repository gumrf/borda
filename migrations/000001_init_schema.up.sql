CREATE TABLE IF NOT EXISTS "user" (
   id serial NOT NULL PRIMARY KEY,
   name varchar(128) NOT NULL,
   password varchar(256) NOT NULL,
   contact varchar(256) UNIQUE NULL
);

CREATE TABLE IF NOT EXISTS role (
   id serial NOT NULL PRIMARY KEY,
   name varchar(256)
);

CREATE TABLE IF NOT EXISTS user_role (
   user_id integer NOT NULL REFERENCES "user"(id),
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
   contact varchar(256)
);

CREATE TABLE IF NOT EXISTS task (
   id serial NOT NULL PRIMARY KEY,
   title varchar(256) NOT NULL,
   description text NOT NULL,
   category varchar(256) NOT NULL,
   complexity varchar(256) NOT NULL,
   points integer NOT NULL,
   hint text,
   flag varchar(256) NOT NULL,
   is_active bool NOT NULL,
   is_disabled bool NOT NULL,
   author_id integer NULL REFERENCES author(id)
);

CREATE TABLE IF NOT EXISTS solved_task (
   task_id integer NOT NULL REFERENCES task(id),
   team_id integer NOT NULL REFERENCES team(id),
   timestamp timestamptz NOT NULL
);

CREATE TABLE IF NOT EXISTS task_submission (
   task_id integer NOT NULL REFERENCES task(id),
   team_id integer NOT NULL REFERENCES team(id),
   user_id integer NOT NULL REFERENCES "user"(id),
   submission varchar(256) NOT NULL,
   is_correct bool NOT NULL,
   timestamp timestamptz NOT NULL
);

CREATE TABLE IF NOT EXISTS manage_settings (
   key varchar(256) NOT NULL,
   value varchar(256) NOT NULL
);