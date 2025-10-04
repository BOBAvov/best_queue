CREATE TABLE IF NOT EXISTS faculties (
    "id" serial PRIMARY KEY,
    "code" varchar(16) UNIQUE NOT NULL,
    "name" varchar(255) NOT NULL,
    "comments" varchar(255)
);

CREATE TABLE IF NOT EXISTS departments (
    "id" serial PRIMARY KEY,
    "code" varchar(32) UNIQUE NOT NULL,
    "name" varchar(255) NOT NULL,
    "faculty_id" integer NOT NULL,
    "comment" varchar(255)
);

CREATE TABLE IF NOT EXISTS groups (
    "id" serial PRIMARY KEY,
    "code" varchar(32) UNIQUE NOT NULL,
    "name" varchar(255) NOT NULL,
    "department_id" integer NOT NULL,
    "comment" varchar(255)
);

CREATE TABLE IF NOT EXISTS users (
    "id" serial PRIMARY KEY,
    "username" varchar(255) NOT NULL,
    "tg_nick" varchar(255) NOT NULL UNIQUE,
    "group_id" integer NOT NULL,
    "password_hash" varchar(255) NOT NULL,
    "is_admin" boolean NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS available (
	"id" serial NOT NULL UNIQUE,
	"type" varchar(255) NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS queues (
	"id" serial NOT NULL UNIQUE,
	"title" varchar(255),
	"group_id" bigint NOT NULL,
	"available_id" bigint NOT NULL,
	"time_start" time without time zone NOT NULL,
	"time_end" time without time zone NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS groups_in_queue (
	"id" bigint NOT NULL UNIQUE,
	"queue_id" bigint NOT NULL,
	"group_id" bigint NOT NULL,
	PRIMARY KEY ("queue_id", "group_id")
);


ALTER TABLE users ADD CONSTRAINT "Users_in_Groups_fk" FOREIGN KEY ("group_id") REFERENCES groups("id");
ALTER TABLE departments ADD CONSTRAINT "Departments_in_Faculties_fk" FOREIGN KEY ("faculty_id") REFERENCES faculties("id");
ALTER TABLE groups ADD CONSTRAINT "Groups_in_Department_fk" FOREIGN KEY ("department_id") REFERENCES departments("id");
ALTER TABLE queues ADD CONSTRAINT "Queues_in_Groups_fk" FOREIGN KEY ("group_id") REFERENCES groups("id");
ALTER TABLE queues ADD CONSTRAINT "Queues_in_Available_fk" FOREIGN KEY ("available_id") REFERENCES available("id");
ALTER TABLE groups_in_queue ADD CONSTRAINT "Groups_queue_in_Queues_fk" FOREIGN KEY ("queue_id") REFERENCES queues("id");
ALTER TABLE groups_in_queue ADD CONSTRAINT "Groups_in_queue_Groups_fk" FOREIGN KEY ("group_id") REFERENCES groups("id")