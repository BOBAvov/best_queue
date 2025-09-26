CREATE TABLE IF NOT EXISTS "Faculties" (
	"id" serial NOT NULL UNIQUE,
	"Title" varchar(255) UNIQUE,
	"Comments" varchar(255) NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "Departments" (
	"id" serial NOT NULL UNIQUE,
	"Title" varchar(255) UNIQUE,
	"Comment" varchar(255) NOT NULL,
	"Faculcy_id" bigint NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "Groups" (
	"id" serial NOT NULL UNIQUE,
	"Title" varchar(255) UNIQUE,
	"Comments" varchar(255) NOT NULL,
	"Departament_id" bigint NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "Users" (
	"id" serial NOT NULL UNIQUE,
	"username" varchar(255) NOT NULL,
	"tg_nick" varchar(255) NOT NULL UNIQUE,
	"email" varchar(255) NOT NULL UNIQUE,
	"Group_id" bigint NOT NULL,
	"Password_hash" varchar(255) NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "Available" (
	"id" serial NOT NULL UNIQUE,
	"Type" varchar(255) NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "Queues" (
	"id" serial NOT NULL UNIQUE,
	"Title" varchar(255),
	"Group_id" bigint NOT NULL,
	"Available_id" bigint NOT NULL,
	"Time_start" time without time zone NOT NULL,
	"Time_end" time without time zone NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "Groups_in_queue" (
	"id" bigint NOT NULL UNIQUE,
	"queue_id" bigint NOT NULL,
	"groupe_id" bigint NOT NULL,
	PRIMARY KEY ("queue_id", "groupe_id")
);

ALTER TABLE "Users" ADD CONSTRAINT "Users_fk4" FOREIGN KEY ("Group_id") REFERENCES "Groups"("id");

ALTER TABLE "Departments" ADD CONSTRAINT "Departments_fk3" FOREIGN KEY ("Faculcy_id") REFERENCES "Faculties"("id");

ALTER TABLE "Groups" ADD CONSTRAINT "Groups_fk3" FOREIGN KEY ("Departament_id") REFERENCES "Departments"("id");

ALTER TABLE "Queues" ADD CONSTRAINT "Queues_fk2" FOREIGN KEY ("Group_id") REFERENCES "Groups"("id");

ALTER TABLE "Queues" ADD CONSTRAINT "Queues_fk3" FOREIGN KEY ("Available_id") REFERENCES "Available"("id");

ALTER TABLE "Groups_in_queue" ADD CONSTRAINT "Groups_in_queue_fk1" FOREIGN KEY ("queue_id") REFERENCES "Queues"("id");

ALTER TABLE "Groups_in_queue" ADD CONSTRAINT "Groups_in_queue_fk2" FOREIGN KEY ("groupe_id") REFERENCES "Groups"("id");