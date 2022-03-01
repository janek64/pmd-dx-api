#!/bin/bash

if [ "$1" == "--help" ];
then
  echo -e " setup <database-url> <port> <database-name> <user> <password> <data-path>\n Creates the tables for the 'pmd-dx-api' database and inserts data from .csv files in <data-path>.";
  exit 0;
fi

if [ $# -lt 6 ];
then
  echo "Missing arguments. See '--help'";
	exit 1;
fi

export PGHOST=$1
export PGPORT=$2
export PGDATABASE=$3
export PGUSER=$4
export PGPASSWORD=$5
export DATAPATH=$6

echo "Starting the database setup...";

echo "Creating all tables...";
psql -f create-tables.sql;
echo "Done.";

echo "Importing data from .csv files...";
psql -c "\copy camp FROM '%DATAPATH%\camp.csv' CSV HEADER";
psql -c "\copy pokemon_type FROM '%DATAPATH%\pokemon_type.csv' CSV HEADER";
psql -c "\copy ability FROM '%DATAPATH%\ability.csv' CSV HEADER";
psql -c "\copy attack_move FROM '%DATAPATH%\attack_move.csv' CSV HEADER";
psql -c "\copy dungeon FROM '%DATAPATH%\dungeon.csv' CSV HEADER";
psql -c "\copy pokemon FROM '%DATAPATH%\pokemon.csv' CSV HEADER";
psql -c "\copy effectiveness FROM '%DATAPATH%\effectiveness.csv' CSV HEADER";
psql -c "\copy encountered_in FROM '%DATAPATH%\encountered_in.csv' CSV HEADER";
psql -c "\copy learns FROM '%DATAPATH%\learns.csv' CSV HEADER";
psql -c "\copy pokemon_has_ability FROM '%DATAPATH%\pokemon_has_ability.csv' CSV HEADER";
psql -c "\copy pokemon_has_type FROM '%DATAPATH%\pokemon_has_type.csv' CSV HEADER";
echo "Done.";