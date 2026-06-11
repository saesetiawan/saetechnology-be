CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   full_name VARCHAR(255) NOT NULL,
   email VARCHAR(255) UNIQUE,
   phone VARCHAR(30) UNIQUE,
   password TEXT NOT NULL,
   avatar_url TEXT,
   status VARCHAR(50) DEFAULT 'active',
   role VARCHAR(50) DEFAULT 'customer',
   email_verified_at TIMESTAMPTZ,
   phone_verified_at TIMESTAMPTZ,
   last_login_at TIMESTAMPTZ,
   created_at TIMESTAMPTZ DEFAULT NOW(),
   updated_at TIMESTAMPTZ DEFAULT NOW(),
   deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
