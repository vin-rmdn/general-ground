# General Ground

General Ground is a short project from vin-rmdn for creating a web server in Go.
Creation of this repository is for an upcoming walk-through on author's Medium
for practicing Software Engineering paradigms and principles.

## Sample curls
`curl --http3 -vv --max-time 1 --cacert ~/cert.pem -X POST -d '{"to": "other", "message": "hi!"}' -H "User-ID: me" https://localhost:8080/chat`

`curl --http3 -vv --max-time 1 --cacert ~/cert.pem -X GET -H "User-ID: me" https://localhost:8080/chat\?with\=other`