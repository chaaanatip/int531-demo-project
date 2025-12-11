DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'sreuser') THEN
    CREATE ROLE sreuser WITH LOGIN PASSWORD 'srepass';
  END IF;
END
$$;

CREATE DATABASE IF NOT EXISTS sredb OWNER sreuser;