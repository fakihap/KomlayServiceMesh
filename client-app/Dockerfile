# FROM node:22-alpine AS build

# WORKDIR /app

# COPY package*.json ./

# RUN npm ci

# COPY . .

# RUN npm run build
# RUN npm prune --production


# FROM node:22-alpine

# WORKDIR /app

# COPY --from=build /app/build build/
# COPY --from=build /app/node_modules node_modules/

# COPY package.json .

# EXPOSE 3000

# ENV NODE_ENV=production

# CMD [ "node", "build" ]


FROM node:22-alpine

WORKDIR /app

COPY package*.json ./

RUN npm ci

COPY . .

EXPOSE 5173

CMD ["npm", "run", "dev"]