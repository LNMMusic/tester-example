# Image Base of MySQL
FROM mysql:8.0

# Copy
COPY ../db/init.sql /docker-entrypoint-initdb.d/init.sql