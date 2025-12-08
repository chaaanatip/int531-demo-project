CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- safe seed: only insert if id not present
INSERT INTO users (id, name) VALUES
(65130500029, 'Guy'),
(65130500039, 'Book'),
(65130500105, 'James')
ON CONFLICT (id) DO NOTHING;
