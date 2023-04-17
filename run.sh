cp .env.example .env

mkdir logs

docker stop servicio-sse && docker rm servicio-sse

docker image rm servicio-sse

docker build -t servicio-sse .

docker run -d \
--restart always \
--name servicio-sse \
--net=upla \
-p 8891:80 \
-v $(pwd)/logs:/etc/push \
servicio-sse