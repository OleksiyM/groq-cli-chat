ARG ARCH=amd64
FROM gcr.io/distroless/static-debian12
COPY ./bin/release/linux-${ARCH}/groq-chat .
CMD ["./groq-chat"]