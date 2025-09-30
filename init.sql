CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    whatsapp INT,
    created_at TIMESTAMP,
    active BOOLEAN
)

CREATE TABLE IF NOT EXISTS workout_sets (
    set_id SERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    user_id INT,
    reps SMALLINT,
    weight SMALLINT 
)

