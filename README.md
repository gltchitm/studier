# Studier
Open-source online studying platform

## Description
Studier is a platform that allows you to create, manage, and share flashcard decks. It includes various studying tools such as write practice and traditional flashcards. A friend system it built in, allowing you to create and manage relationships with others.

## Running
You will need Docker to run Studier. Before starting, create `config/config.json` and fill it out based on the template in `config/config.json.example`. Once you have created the configuration, place the server's TLS certificate in `config/server.crt` and place the key in `config/server.key`. After doing this, you can build and start Studier through Docker Compose:
```sh
docker-compose -f docker-compose.yml up --build
```
Studier will start exclusively over HTTPS on port 443.

## License
[MIT](LICENSE)
