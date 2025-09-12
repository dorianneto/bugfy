FROM node:20-alpine

RUN apk update
RUN apk upgrade

WORKDIR /app

COPY package*.json .
RUN npm ci
COPY . .

CMD ["npm", "run", "dev"]
