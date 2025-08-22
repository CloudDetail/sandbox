FROM node:18-alpine

WORKDIR /app

RUN apk add --no-cache bash iproute2

COPY package*.json ./
RUN npm ci --only=production

COPY . .

USER root

EXPOSE 3500

RUN npm install --save @opentelemetry/api && npm install --save @opentelemetry/auto-instrumentations-node

CMD ["npm", "start"]
