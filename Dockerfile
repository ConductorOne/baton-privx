FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-privx"]
COPY baton-privx /