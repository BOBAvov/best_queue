CREATE TABLE IF NOT EXISTS "Faculties" (
    "id" serial PRIMARY KEY,
    "code" varchar(16) UNIQUE NOT NULL,
    "name" varchar(255) NOT NULL,
    "comments" varchar(255)
);

CREATE TABLE IF NOT EXISTS "Departments" (
    "id" serial PRIMARY KEY,
    "code" varchar(32) UNIQUE NOT NULL,
    "name" varchar(255) NOT NULL,
    "faculty_id" integer NOT NULL,
    "comment" varchar(255)
);

CREATE TABLE IF NOT EXISTS "Groups" (
    "id" serial PRIMARY KEY,
    "code" varchar(32) UNIQUE NOT NULL,
    "name" varchar(255) NOT NULL,
    "department_id" integer NOT NULL,
    "comment" varchar(255)
);

CREATE TABLE IF NOT EXISTS "Users" (
    "id" serial PRIMARY KEY,
    "username" varchar(255) NOT NULL,
    "tg_nick" varchar(255) NOT NULL UNIQUE,
    "group_id" integer NOT NULL,
    "password_hash" varchar(255) NOT NULL,
    "is_admin" boolean NOT NULL
);

CREATE TABLE IF NOT EXISTS "Available" (
	"id" serial NOT NULL UNIQUE,
	"type" varchar(255) NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "Queues" (
	"id" serial NOT NULL UNIQUE,
	"title" varchar(255),
	"group_id" bigint NOT NULL,
	"available_id" bigint NOT NULL,
	"time_start" time without time zone NOT NULL,
	"time_end" time without time zone NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "Groups_in_queue" (
	"id" bigint NOT NULL UNIQUE,
	"queue_id" bigint NOT NULL,
	"groupe_id" bigint NOT NULL,
	PRIMARY KEY ("queue_id", "groupe_id")
);


ALTER TABLE "Users" ADD CONSTRAINT "Users_fk4" FOREIGN KEY ("group_id") REFERENCES "Groups"("id");
ALTER TABLE "Departments" ADD CONSTRAINT "Departments_fk3" FOREIGN KEY ("faculty_id") REFERENCES "Faculties"("id");
ALTER TABLE "Groups" ADD CONSTRAINT "Groups_fk3" FOREIGN KEY ("department_id") REFERENCES "Departments"("id");
ALTER TABLE "Queues" ADD CONSTRAINT "Queues_fk2" FOREIGN KEY ("group_id") REFERENCES "Groups"("id");
ALTER TABLE "Queues" ADD CONSTRAINT "Queues_fk3" FOREIGN KEY ("available_id") REFERENCES "Available"("id");
ALTER TABLE "Groups_in_queue" ADD CONSTRAINT "Groups_in_queue_fk1" FOREIGN KEY ("queue_id") REFERENCES "Queues"("id");
ALTER TABLE "Groups_in_queue" ADD CONSTRAINT "Groups_in_queue_fk2" FOREIGN KEY ("groupe_id") REFERENCES "Groups"("id")