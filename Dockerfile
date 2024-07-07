# Utiliser une image de base légère pour Go
FROM golang:1.20-alpine

# Définir le répertoire de travail à l'intérieur du container
WORKDIR /app

# Copier le code source dans le répertoire de travail
COPY . .

# Construire l'application
RUN go build -o gateway-server

# Exposer le port
EXPOSE 8080

# Commande pour lancer le serveur
CMD ["./gateway-server"]
