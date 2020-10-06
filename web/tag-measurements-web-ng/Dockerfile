# этап сборки (build stage)
FROM node:lts-alpine as build-stage
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm install -g @angular/cli
RUN ng build --prod

# этап production (production-stage)
FROM nginx:stable-alpine as production-stage
COPY --from=build-stage /app/nginx_config/nginx.conf /etc/nginx/nginx.conf
RUN mkdir -p /var/log/nginx
RUN rm -rf /usr/share/nginx/html/*
COPY --from=build-stage /app/dist/tag-measurements-web-ng /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
