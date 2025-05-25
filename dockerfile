FROM alpine:latest
WORKDIR /app
COPY my-go-binary .
RUN chmod +x my-go-binary
CMD ["./my-go-binary"]