FROM node:24-alpine AS build
WORKDIR /src
COPY package.json package-lock.json tsconfig.base.json ./
COPY packages/domain/package.json packages/domain/package.json
COPY apps/web/package.json apps/web/package.json
RUN npm ci
COPY packages/domain packages/domain
COPY apps/web apps/web
RUN npm run build -w @dcraft-fusion/web

FROM nginx:1.29-alpine
COPY apps/web/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /src/apps/web/dist /usr/share/nginx/html
EXPOSE 80
