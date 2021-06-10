FROM scratch

# Copy the existing binary, must be named kvs, from the host.
COPY kvs .

EXPOSE 8080

CMD ["/kvs"]