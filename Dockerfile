LABEL org.opencontainers.image.source=https://github.com/ConductorOne/baton-minecraft-luckperms
FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-minecraft-luckperms"]
COPY baton-minecraft-luckperms /