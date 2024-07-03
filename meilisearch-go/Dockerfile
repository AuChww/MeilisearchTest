# Stage 1: Build wal2json from source
FROM postgres:latest AS wal2json-builder

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
       git \
       build-essential \
       postgresql-server-dev-all \
       ca-certificates \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN git clone https://github.com/eulerto/wal2json.git \
    && cd wal2json \
    && make && make install

# Verify the installation paths
RUN find /usr -name wal2json -type f

# Stage 2: Final image
FROM postgres:latest

# Copy the wal2json.so file to the PostgreSQL library directory
COPY --from=wal2json-builder /usr/lib/postgresql /usr/lib/postgresql

# Optionally copy configuration files
# COPY ./postgresql.conf /etc/postgresql/postgresql.conf
# COPY ./pg_hba.conf /etc/postgresql/pg_hba.conf

# Expose PostgreSQL port
EXPOSE 5432

CMD ["postgres"]