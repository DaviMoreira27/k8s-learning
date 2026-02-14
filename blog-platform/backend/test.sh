curl -X GET http://localhost:8080/posts \
        -H "Content-Type: application/json"

curl -X POST http://localhost:8080/posts \
        -H "Content-Type: application/json" \
        -d '{"title":"First Post","text":"Hello Kubernetes"}'
