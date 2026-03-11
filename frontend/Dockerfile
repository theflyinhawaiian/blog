FROM node:20-alpine AS base

ARG BACKEND_URL=http://localhost:8080
ENV BACKEND_URL=$BACKEND_URL

WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm install
COPY . .

FROM base AS dev
EXPOSE 3000
CMD ["npm", "run", "dev", "--", "--host"]

FROM base AS builder
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 3000
