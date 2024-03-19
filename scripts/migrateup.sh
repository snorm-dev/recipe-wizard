if [ -f .env ]; then
    source .env
fi

cd sql/schema

goose sqlite3 $DATABASE_URL up
