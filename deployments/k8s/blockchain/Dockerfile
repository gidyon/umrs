FROM alpine
LABEL maintainer="gideonhacer@gmail.com"
RUN apk update && \
   apk add ca-certificates && \
   update-ca-certificates && \
   rm -rf /var/cache/apk/* && \
   apk add libc6-compat

WORKDIR /app
COPY ledger .
EXPOSE 80 443
ENTRYPOINT [ "/app/ledger" ]