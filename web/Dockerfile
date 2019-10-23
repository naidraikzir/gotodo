FROM node:10.16.3-alpine
RUN npm i -g @vue/cli @vue/cli-service-global
WORKDIR /app
COPY . .
EXPOSE 8080
CMD vue serve
