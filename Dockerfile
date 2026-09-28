# One image serving the blog at /, the admin at /admin/ and the API at /api/v1.
# Build from this directory: docker compose up -d --build

FROM node:24-alpine AS web
ARG NPM_REGISTRY=https://registry.npmmirror.com
WORKDIR /src
COPY web/package.json web/package-lock.json web/
COPY web-admin/package.json web-admin/package-lock.json web-admin/
RUN cd web && npm ci --registry "$NPM_REGISTRY" --no-fund --no-audit \
 && cd ../web-admin && npm ci --registry "$NPM_REGISTRY" --no-fund --no-audit
COPY web web
COPY web-admin web-admin
RUN cd web && npm run build && cd ../web-admin && npm run build

FROM golang:1.26.8-alpine AS server
ARG GOPROXY=https://goproxy.cn,direct
WORKDIR /src
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/blog-server ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=server /out/blog-server /blog-server
COPY --from=web /src/web/build /srv/web
COPY --from=web /src/web-admin/build /srv/admin
ENV APP_ENV=production \
    HTTP_ADDR=:8080 \
    WEB_DIR=/srv/web \
    ADMIN_DIR=/srv/admin
EXPOSE 8080
ENTRYPOINT ["/blog-server"]
