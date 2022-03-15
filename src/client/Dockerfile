FROM node:16.14.0-bullseye

WORKDIR /studier/client

COPY . .

RUN npm ci && npm run build

CMD ["node", "build"]
