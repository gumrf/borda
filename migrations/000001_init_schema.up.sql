CREATE TABLE IF NOT EXISTS users (
   id serial NOT NULL PRIMARY KEY,
   name varchar(128) NOT NULL,
   password varchar(256) NOT NULL,
   contact varchar(256) UNIQUE NULL
);

CREATE TABLE IF NOT EXISTS roles (
   id serial NOT NULL PRIMARY KEY,
   name varchar(256)
);

CREATE TABLE IF NOT EXISTS user_roles (
   user_id integer NOT NULL REFERENCES users(id),
   role_id integer NOT NULL REFERENCES roles(id)
);

CREATE TABLE IF NOT EXISTS teams (
   id serial NOT NULL PRIMARY KEY,
   name varchar(256) NOT NULL,
   team_leader_id integer NOT NULL REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS team_members (
   team_id integer NOT NULL PRIMARY KEY REFERENCES teams(id),
   user_id integer NOT NULL REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS authors (
   id serial NOT NULL PRIMARY KEY,
   name varchar(256) NOT NULL,
   contact varchar(256)
);

CREATE TABLE IF NOT EXISTS tasks (
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
   author_id integer NULL REFERENCES authors(id)
);

CREATE TABLE IF NOT EXISTS solved_tasks (
   task_id integer NOT NULL REFERENCES tasks(id),
   team_id integer NOT NULL REFERENCES teams(id),
   timestamp timestamptz NOT NULL
);

CREATE TABLE IF NOT EXISTS task_submissions (
   task_id integer NOT NULL REFERENCES tasks(id),
   team_id integer NOT NULL REFERENCES teams(id),
   user_id integer NOT NULL REFERENCES users(id),
   submission varchar(256) NOT NULL,
   is_correct bool NOT NULL,
   timestamp timestamptz NOT NULL
);