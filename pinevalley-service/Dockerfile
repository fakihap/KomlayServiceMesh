FROM node:slim

WORKDIR /app

COPY package.json .

RUN npm install

COPY . .

CMD [ "npm", "run", "pine-valley" ]

EXPOSE 9090