# Dockerfile za project-service

# Koristi golang bazni imidž
FROM golang:1.20-alpine

# Podešava radni direktorijum u kontejneru
WORKDIR /app

# Kopira go.mod i go.sum fajlove
COPY go.mod go.sum ./

# Instalira zavisnosti
RUN go mod download

# Kopira sve fajlove projekta u radni direktorijum u kontejneru
COPY . .

# Kompajlira aplikaciju
RUN go build -o project-service

# Definiše komandu za pokretanje aplikacije
CMD ["./project-service"]
