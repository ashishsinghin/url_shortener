# URL shortener

To Run the application 

`go mod tidy`

`go run main.go`


Connect to PostgreSQL and create a database for your URL shortener:

`CREATE DATABASE url_shortener;`

Switch to the new database:

`url_shortener`

Create a table to store URLs:

`CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    short_key VARCHAR(10) UNIQUE NOT NULL,
    original_url TEXT NOT NULL
);`

edit your db `username` and `password` in `.env` file

## Create a Short URL:

### Use curl or any HTTP client to POST a URL:

`curl -X POST http://localhost:8080/shorten -d '{"url":"https://urvashisingh.in/"}' -H "Content-Type: application/json"`

## Access the Short URL:

### Open the short URL in your browser or use curl to test the redirection:

`curl -i http://localhost:8080/{short_key}`

Replace {short_key} with the short key you received in the response.

## Deployment

It's hosted on AWS EC2 instance with public Ip `13.200.237.89`

You can test application using below link

`curl -X POST http://13.200.237.89:8080/shorten -d '{"url":"https://urvashisingh.in/"}' -H "Content-Type: application/json"`

`curl -i http://13.200.237.89:8080/{short_key}`