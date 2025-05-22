FROM node:20.19-alpine AS builder

WORKDIR /src

# this variable need at build stage
ARG NEXT_PUBLIC_BACKEND_URL
ENV NEXT_PUBLIC_BACKEND_URL=$NEXT_PUBLIC_BACKEND_URL

# download and install dependencies
COPY package.json package-lock.json ./
RUN npm ci

# build app
COPY . .
RUN sed -i "/const nextConfig/a \ output: 'standalone'," next.config.ts
RUN npm run build

# 2nd stage - run
FROM node:20.19-alpine

ENV NODE_ENV=production

WORKDIR /src

# copy app files
# COPY --from=builder /src/.env .
COPY --from=builder /src/public ./public
COPY --from=builder /src/.next/static ./.next/static
COPY --from=builder /src/.next/standalone ./

EXPOSE 3000

# run app
CMD ["node", "server.js"]