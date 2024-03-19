if [ -f .env ]; then
    source .env
fi

./scripts/libsql-goose down
