CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    whatsapp INT,
    created_at TIMESTAMP,
    active BOOLEAN
)

CREATE TABLE IF NOT EXISTS workout_sets (
    set_id SERIAL PRIMARY KEY,
    user_id INT,
    exercise VARCHAR(100),
    created_at TIMESTAMP,
    reps SMALLINT,
    weight SMALLINT 
)

