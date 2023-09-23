FROM --platform=linux/amd64 debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates

RUN ls

ADD recipe-wizard /usr/bin/recipe-wizard

CMD ["recipe-wizard"]
