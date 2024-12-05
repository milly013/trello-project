FROM nginx:latest

# Kopiranje NGINX konfiguracije
COPY ./api_gateway.conf /etc/nginx/conf.d/api_gateway.conf

# Brisanje podrazumevane konfiguracije
RUN rm /etc/nginx/conf.d/default.conf

# Kopiranje sertifikata u kontejner
COPY ./certs/cert.crt /etc/nginx/certs/cert.crt
COPY ./certs/cert.key /etc/nginx/certs/cert.key
