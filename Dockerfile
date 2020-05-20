# docker build -t user-service:v0.0.2 .
# docker run -it --rm --network=host user-service:v0.0.2
FROM golang:1.13

# copy app and config
COPY user-configuration.json user-configuration.json
COPY user-service app

# container port
EXPOSE 8080

# run
CMD ["./app", "start"]
