FROM postgres:9.4

RUN apt-get update \
      && apt-get install -y \
           postgresql-$PG_MAJOR-postgis-2.1 \
      && rm -rf /var/lib/apt/lists/*

ADD ddl /docker-entrypoint-initdb.d
