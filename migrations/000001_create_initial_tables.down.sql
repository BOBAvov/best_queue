ALTER TABLE queue_participants DROP CONSTRAINT IF EXISTS "Queue_participants_user_fk";
ALTER TABLE queue_participants DROP CONSTRAINT IF EXISTS "Queue_participants_queue_fk";
ALTER TABLE users DROP CONSTRAINT IF EXISTS "Users_in_Groups_fk";

DROP TABLE IF EXISTS queue_participants;
DROP TABLE IF EXISTS queues;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS groups;