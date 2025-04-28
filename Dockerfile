FROM gcr.io/distroless/static-debian12
COPY ./bin/release/linux-amd64/groq-chat .
CMD ["./groq-chat"]