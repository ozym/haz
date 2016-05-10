# This docker file should be invoked from docker.sh
FROM scratch
COPY docker-build-tmp/quakesearch-go /
COPY docker-build-tmp/common /
EXPOSE 8080
USER nobody
CMD ["/quakesearch-go"]
