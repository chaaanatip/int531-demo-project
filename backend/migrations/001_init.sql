CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- safe seed: only insert if id not present
INSERT INTO users (id, name) VALUES
(65130500098, 'Chai'),
(65130500069, 'leng'),
(65130500095, 'Oat'),
(65130500116, 'Save')
ON CONFLICT (id) DO NOTHING;
