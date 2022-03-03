-- Declare all custom enum types
DROP TYPE IF EXISTS evolve_condition CASCADE;
CREATE TYPE evolve_condition AS ENUM('level', 'crystal', 'no_evolve');

DROP TYPE IF EXISTS move_category CASCADE;
CREATE TYPE move_category AS ENUM('Physical', 'Special', 'Status');

DROP TYPE IF EXISTS move_range CASCADE;
CREATE TYPE move_range AS ENUM('Entire room', '4 tiles away', '2 tiles away', '1 tile around user', 'Front', 'User');

DROP TYPE IF EXISTS move_target CASCADE;
CREATE TYPE move_target AS ENUM('Enemy', 'Party', 'User', 'All except user');

DROP TYPE IF EXISTS camp_unlock_type CASCADE;
CREATE TYPE camp_unlock_type AS ENUM('obtain', 'buy');

DROP TYPE IF EXISTS move_learn_type CASCADE;
CREATE TYPE move_learn_type AS ENUM('level', 'tutor', 'tm');

DROP TYPE IF EXISTS type_interaction CASCADE;
CREATE TYPE type_interaction AS ENUM('super effective', 'not very effective', 'not effective');

-- Declare all tables
DROP TABLE IF EXISTS camp CASCADE;
CREATE TABLE camp (
  camp_ID smallserial PRIMARY KEY,
  camp_name varchar(50) NOT NULL,
  unlock_type camp_unlock_type NOT NULL,
  cost smallint,
  description varchar(300) NOT NULL
);

DROP TABLE IF EXISTS pokemon CASCADE;
CREATE TABLE pokemon (
  dex_number smallint NOT NULL PRIMARY KEY,
  pokemon_name varchar(50) NOT NULL,
  evolution_stage smallint,
  evolve_condition evolve_condition NOT NULL,
  evolve_level smallint,
  evolve_crystals smallint,
  classification varchar(50) NOT NULL
);

DROP TABLE IF EXISTS pokemon_type CASCADE;
CREATE TABLE pokemon_type (
  type_ID smallserial PRIMARY KEY,
  type_name varchar(50) NOT NULL
);

DROP TABLE IF EXISTS attack_move CASCADE;
CREATE TABLE attack_move (
  move_ID smallserial PRIMARY KEY,
  move_name varchar(50) NOT NULL,
  category move_category NOT NULL,
  move_range move_range NOT NULL,
  target move_target NOT NULL,
  initial_pp smallint NOT NULL,
  initial_power smallint NOT NULL,
  accuracy smallint NOT NULL,
  description varchar(300) NOT NULL
);

DROP TABLE IF EXISTS ability CASCADE;
CREATE TABLE ability (
  ability_ID smallserial PRIMARY KEY,
  ability_name varchar(50) NOT NULL,
  description varchar(300) NOT NULL
);

DROP TABLE IF EXISTS dungeon CASCADE;
CREATE TABLE dungeon (
  dungeon_ID smallserial PRIMARY KEY,
  dungeon_name varchar(50) NOT NULL,
  levels smallint NOT NULL,
  start_level smallint,
  team_size smallint NOT NULL,
  items_allowed boolean NOT NULL,
  pokemon_joining boolean NOT NULL,
  map_visible boolean NOT NULL  
);

DROP TABLE IF EXISTS encountered_in;
CREATE TABLE encountered_in (
  dex_number smallint NOT NULL REFERENCES pokemon (dex_number),
  dungeon_ID smallint NOT NULL REFERENCES dungeon (dungeon_ID),
  super_enemy boolean NOT NULL,
  PRIMARY KEY(dex_number, dungeon_ID)
);

DROP TABLE IF EXISTS learns;
CREATE TABLE learns (
  learns_ID smallserial,
  dex_number smallint NOT NULL REFERENCES pokemon (dex_number),
  move_ID smallint NOT NULL REFERENCES attack_move (move_ID),
  learn_type move_learn_type NOT NULL,
  cost smallint,
  level smallint,
  PRIMARY KEY(learns_ID, dex_number, move_ID)
);

DROP TABLE IF EXISTS pokemon_has_ability;
CREATE TABLE pokemon_has_ability (
  dex_number smallint NOT NULL REFERENCES pokemon (dex_number),
  ability_ID smallint NOT NULL REFERENCES ability (ability_ID),
  PRIMARY KEY(dex_number, ability_ID)
);

DROP TABLE IF EXISTS pokemon_has_type;
CREATE TABLE pokemon_has_type (
  dex_number smallint NOT NULL REFERENCES pokemon (dex_number),
  type_ID smallint NOT NULL REFERENCES pokemon_type (type_ID),
  PRIMARY KEY(dex_number, type_ID)
);

DROP TABLE IF EXISTS effectiveness;
CREATE TABLE effectiveness (
  attacker smallint REFERENCES pokemon_type (type_ID),
  defender smallint REFERENCES pokemon_type (type_ID),
  interaction type_interaction NOT NULL,
  PRIMARY KEY(attacker, defender)
);

-- Add missing foreign keys
ALTER TABLE pokemon ADD COLUMN camp_id smallint NOT NULL REFERENCES camp (camp_ID);

ALTER TABLE attack_move ADD COLUMN type_ID smallint NOT NULL REFERENCES pokemon_type (type_ID);
