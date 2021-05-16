# Contacts

Usage:
```sh
# Run
go run .

# Create
curl -d "name=Jane&phone=88005553535" -X POST http://localhost:8080/contacts
curl -d "name=Alice&phone=89992221122" -X POST http://localhost:8080/contacts

# Get all
curl http://localhost:8080/contacts

# Get contact by ID
curl http://localhost:8080/contact/1

# Edit
curl -X PATCH -d "name=Mary" http://localhost:8080/contact/1

# Delete
curl -X DELETE http://localhost:8080/contact/1
```