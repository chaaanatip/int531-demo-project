# frontend/Dockerfile (static)
FROM node:20 AS builder
WORKDIR /app
COPY package.json package-lock.json* ./

# install deps (lockfile recommended)
RUN if [ -f package-lock.json ]; then npm ci; else npm install; fi

COPY . .

# build will produce 'out' because next.config.js sets output: 'export'
ARG NEXT_PUBLIC_API_BASE
ENV NEXT_PUBLIC_API_BASE=${NEXT_PUBLIC_API_BASE:-""}
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/.next/static /usr/share/nginx/html/_next/static
COPY --from=builder /app/out /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
