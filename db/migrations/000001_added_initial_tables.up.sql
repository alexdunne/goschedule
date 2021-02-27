CREATE TABLE IF NOT EXISTS users (
   "id" VARCHAR(255) NOT NULL, 
   "name" VARCHAR(255) NOT NULL,
   "email" VARCHAR(255) UNIQUE NOT NULL, 
   "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL, 
   "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
   PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS user_external_logins (
   "id" VARCHAR(255) NOT NULL, 
   "user_id" VARCHAR(255) NOT NULL,
   "source" VARCHAR(255) NOT NULL, 
   "source_id" TEXT NOT NULL, 
   "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL, 
   "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
   PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS organisations (
   "id" VARCHAR(255) NOT NULL, 
   "name" VARCHAR(255) NOT NULL,
   "owner_id" VARCHAR(255) NOT NULL, 
   "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL, 
   "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
   PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS schedules (
   "id" VARCHAR(255) NOT NULL, 
   "name" VARCHAR(255) NOT NULL,
   "owner_id" VARCHAR(255) NOT NULL, 
   "organisation_id" VARCHAR(255) NOT NULL, 
   "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL, 
   "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
   PRIMARY KEY(id)
);

CREATE INDEX IF NOT EXISTS IX_user_external_logins_user_id ON user_external_logins (user_id);
CREATE INDEX IF NOT EXISTS IX_organisations_owner_id ON organisations (owner_id);
CREATE INDEX IF NOT EXISTS IX_schedules_owner_id ON schedules (owner_id);
CREATE INDEX IF NOT EXISTS IX_schedules_organisation_id ON schedules (organisation_id);


ALTER TABLE user_external_logins ADD CONSTRAINT FK_users_user_external_logins FOREIGN KEY (user_id) REFERENCES users (id) NOT DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE organisations ADD CONSTRAINT FK_users_organisations FOREIGN KEY (owner_id) REFERENCES users (id) NOT DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE schedules ADD CONSTRAINT FK_users_schedules FOREIGN KEY (owner_id) REFERENCES users (id) NOT DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE schedules ADD CONSTRAINT FK_organisations_schedules FOREIGN KEY (organisation_id) REFERENCES organisations (id) NOT DEFERRABLE INITIALLY IMMEDIATE;