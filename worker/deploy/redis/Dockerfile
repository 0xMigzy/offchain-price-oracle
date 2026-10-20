FROM redis:7.4-alpine
COPY start.sh /usr/local/bin/start-demo-redis
RUN chmod +x /usr/local/bin/start-demo-redis
ENV REDIS_MAXMEMORY=64mb
EXPOSE 6379
HEALTHCHECK --interval=30s --timeout=3s CMD REDISCLI_AUTH="$REDIS_PASSWORD" redis-cli ping || exit 1
ENTRYPOINT ["/usr/local/bin/start-demo-redis"]
