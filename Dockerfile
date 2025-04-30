ARG ARCH=amd64
FROM gcr.io/distroless/static-debian12
COPY ./groq-chat .
CMD ["./groq-chat"]