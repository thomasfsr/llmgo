CREATE TABLE IF NOT EXISTS users ( 
  id PRIMARY KEY AUTOINCREMENT, 
  first_name CHAR(20), 
  last_name CHAR(50), 
  whatsapp INTEGER, 
  active BOOLEAN, 
  created_at TIMESTAMP,
);

CREATE TABLE IF NOT EXISTS messages (
  id INTEGER PRIMARY KEY AUTOINCREMENT, 
  user_id INTEGER,
  message CHAR(100),
  created_at TIMESTAMP,
);

CREATE TABLE IF NOT EXISTS exercise_sets (
  id INTEGER PRIMARY KEY AUTOINCREMENT, 
  exercise CHAR(100),
  weight INTEGER,
  reps INTEGER,
  created_at TIMESTAMP,
); 
