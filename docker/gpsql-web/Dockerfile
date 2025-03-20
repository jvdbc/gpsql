# Étape 1 : Build
FROM public.ecr.aws/docker/library/golang:1.23.8 AS builder

# Définir le répertoire de travail
WORKDIR /app

# Copier les fichiers du projet dans l'image
COPY ./cmd/gpsql-web .
COPY ./go.* .

# Compiler l'application en un binaire statique
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gpsql-web .

# Étape 2 : Image minimale
FROM scratch

# Copier le binaire depuis l'étape de build
COPY --from=builder /app/gpsql-web /gpsql-web

# Note: L'instruction EXPOSE est uniquement à des fins de documentation et ne publie pas réellement le port.
EXPOSE 8000

# Commande par défaut pour exécuter l'application
CMD ["./gpsql-web"]