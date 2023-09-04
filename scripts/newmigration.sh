if [ $# -ne 1 ]; then
    echo "You must pass in a single argument for the migration name"
    exit 1
fi

cd sql/schema
goose create $1 sql