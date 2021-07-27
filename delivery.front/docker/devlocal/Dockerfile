FROM golang:1.15-alpine as build

LABEL maintainer="bagasl98@gmail.com"

# git is required to fetch go dependencies
RUN apk add --no-cache \
  ca-certificates \
  git 

# Create the user and group files that will be used in the running 
# container to run the process as an unprivileged user.
RUN mkdir /user && \
  echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
  echo 'nobody:x:65534:' > /user/group

# Access tokens for downloading shared module from faem repo
ARG ACCESS_TOKEN_MACHINE="gitlab.com" 
ARG ACCESS_TOKEN_USR="vlad.ja@faem.me"
ARG ACCESS_TOKEN_PWD="wso-nWAsAPpsyi3spCss"

# Create a netrc file using the credentials specified using --build-arg
RUN printf "machine ${ACCESS_TOKEN_MACHINE}\n\
  login ${ACCESS_TOKEN_USR}\n\
  password ${ACCESS_TOKEN_PWD}\n"\
  >> /root/.netrc
RUN chmod 600 /root/.netrc

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src

# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build
COPY ./go.mod ./go.sum ./
RUN go mod download

# Import the code from the context.
COPY . .

# Build the executable to `/app`. Mark the build as statically linked.
RUN CGO_ENABLED=0 GONOSUMDB=on go build \
  -o /app /src/cmd

# Final stage: the running container.
FROM scratch AS final

# Import the user and group files from the first stage.
COPY --from=build /user/group /user/passwd /etc/
# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Import the compiled executable and config file from the first stage.
COPY --from=build /app /app
COPY --from=build /src/deployment/config/deliveryfront.toml /deliveryfront.toml

# Perform any further action as an unprivileged user.
USER nobody:nobody

ARG MAIN_PORT=1400

# Specify application's main port
EXPOSE ${MAIN_PORT}

ENV TZ=Europe/Moscow 

STOPSIGNAL SIGINT

# Run the compiled binary and pass the config to the /app
ENTRYPOINT [ "/app" ]
CMD ["-config", "/deliveryfront.toml"]
