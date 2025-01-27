CREATE TABLE IF NOT EXISTS channels (
    id SERIAL PRIMARY KEY,
    full_name TEXT NOT NULL,  
    email TEXT UNIQUE NOT NULL,
    public_key TEXT NOT NULL,
    public_url TEXT NOT NULL,
    public_qr TEXT NOT NULL,
    verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);

