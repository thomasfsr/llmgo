CREATE TABLE IF NOT EXISTS users ( 
  id SERIAL PRIMARY KEY, 
  first_name CHAR(20), 
  last_name CHAR(50), 
  whatsapp BIGINT, 
  active BOOLEAN, 
  created_at TIMESTAMP,
)

CREATE TABLE IF NOT EXISTS messages (
  id SERIAL PRIMARY KEY,
  user_id INT,
  message CHAR(100),
  created_at TIMESTAMP,
)

CREATE TABLE IF NOT EXISTS exercise_sets (
  id SERIAL PRIMARY KEY,
  exercise CHAR(100),
  weight INT,
  reps INT,
  created_at TIMESTAMP,
) 
