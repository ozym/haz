FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/
ADD zoneinfo /usr/share/zoneinfo
ADD geonet-rest /
ADD tmpl /tmpl
ADD geonet-rest.json /
EXPOSE 8080
CMD ["/geonet-rest"]
