FROM golang:1.12 as build-stage

RUN set -x

RUN mkdir /build

WORKDIR /build

COPY . /build

# RUN git clone https://github.com/kuritka/webhook-abbsa-interview . && \
RUN  go mod vendor  && \
     go list -e $(go list -f . -m all) && \
#viz https://drailing.net/2018/02/building-go-binaries-for-small-docker-images/
     CGO_ENABLED=0 go build -a -o main . && \
     groupadd -g 1001 lb && \
     useradd -r -u 1001 -g lb lb

#------------------------------------------------------------  << 20MB
FROM scratch as release-stage

WORKDIR /app

#multistage containers - copying from build stage /build to /app
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-stage /build/main /app/main
COPY --from=build-stage /build/server/public /app/server/public
COPY --from=build-stage /build/server/templates /app/server/templates


#scratch is missing bash, so cannot call useradd command. That's we created user at build-stage, now we copy him to scratch
COPY --from=build-stage /etc/passwd /etc/passwd

USER lb

ENTRYPOINT ["./main"]


#delete all <none> images
#sudo docker rmi $(sudo docker images | grep "^<none>" | awk '{ print $3 }')
#docker container prune
#docker image prune
#docker network prune
#docker volume prune
# https://www.projectatomic.io/blog/2015/07/what-are-docker-none-none-images/
# docker rmi $(docker images -f "dangling=true" -q)



