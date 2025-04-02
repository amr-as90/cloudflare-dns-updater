FROM alpine:latest

WORKDIR /app

COPY cloudflare-dns-updater .

ENV ZONE_ID=""
ENV RECORD_NAMES=""
ENV EMAIL=""
ENV ACCESS_TOKEN=""
ENV UPDATE_INTERVAL="60"
ENV PUSHOVER_APP_TOKEN=""
ENV PUSHOVER_USER_KEY=""

RUN echo "ZONE_ID=\"$ZONE_ID\"" > .config && \
    echo "RECORD_NAMES=\"$RECORD_NAMES\"" >> .config && \
    echo "EMAIL=\"$EMAIL\"" >> .config && \
    echo "ACCESS_TOKEN=\"$ACCESS_TOKEN\"" >> .config && \
    echo "UPDATE_INTERVAL=\"$UPDATE_INTERVAL\"" >> .config && \
    echo "PUSHOVER_APP_TOKEN=\"$PUSHOVER_APP_TOKEN\"" >> .config && \
    echo "PUSHOVER_USER_KEY=\"$PUSHOVER_USER_KEY\"" >> .config

CMD ["./cloudflare-dns-updater"]