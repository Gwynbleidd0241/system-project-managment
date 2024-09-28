
CREATE TABLE users (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   email VARCHAR(255) NOT NULL UNIQUE ,
   password TEXT NOT NULL
);

-- DROP TABLE IF EXISTS users;
