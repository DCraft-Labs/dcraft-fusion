# syntax=docker/dockerfile:1.7
#
# Fusion SPA (apps/web) image.
#
# IMPORTANT: the production bundle is ALWAYS built inside this image from
# source via `npm run build -w @dcraft-fusion/web`. A pre-built `dist/` must
# NEVER be copied in from the host — `.dockerignore` excludes `**/dist/` and
# the build stage explicitly removes any stale `dist/` before rebuilding so
# a cached/committed bundle can never leak into production.
#
# If you change this Dockerfile, keep the `rm -rf apps/web/dist` line below
# and keep `dist/` listed in `.dockerignore` and `.gitignore`.
FROM node:24-alpine AS build
WORKDIR /src
COPY package.json package-lock.json tsconfig.base.json ./
COPY packages/domain/package.json packages/domain/package.json
COPY apps/web/package.json apps/web/package.json
RUN npm ci
COPY packages/domain packages/domain
COPY apps/web apps/web
# Defensive: nuke any stale bundle that may have been carried in by COPY.
RUN rm -rf apps/web/dist
RUN npm run build -w @dcraft-fusion/web
# Fail the build if the bundle is missing or older than the newest source file.
RUN test -d apps/web/dist && \
    NEWEST_SRC="$(find apps/web/src packages/domain/src -type f -printf '%T@ %p\n' | sort -nr | head -1 | cut -d' ' -f2)" && \
    NEWEST_DIST="$(find apps/web/dist -type f -printf '%T@ %p\n' | sort -nr | head -1 | cut -d' ' -f2)" && \
    if [ -n "${NEWEST_SRC}" ] && [ -n "${NEWEST_DIST}" ] && \
       [ "$(stat -c %Y "${NEWEST_DIST}")" -lt "$(stat -c %Y "${NEWEST_SRC}")" ]; then \
      echo "::error:: apps/web/dist is older than newest src/ — bundle not rebuilt" >&2; exit 1; \
    fi

FROM nginx:1.29-alpine
COPY apps/web/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /src/apps/web/dist /usr/share/nginx/html
EXPOSE 80
