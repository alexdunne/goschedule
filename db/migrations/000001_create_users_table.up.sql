CREATE TABLE IF NOT EXISTS users (
   "id" TEXT PRIMARY KEY, 
   "name" VARCHAR(255) NOT NULL,
   "email" VARCHAR(255) UNIQUE NOT NULL, 
   "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP, 
   "updated_at" TIMESTAMP(3) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_external_logins (
   "id" TEXT PRIMARY KEY, 
   "user_id" TEXT REFERENCES users (id) ON DELETE CASCADE,
   "source" VARCHAR(255) NOT NULL, 
   "source_id" TEXT NOT NULL, 
   "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP, 
   "updated_at" TIMESTAMP(3) NOT NULL
);
