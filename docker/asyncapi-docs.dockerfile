# Use node:14 docker image b/c asyncapi/generator is not yet compatible with node:15
FROM node:14

### Install dependencies
RUN npm install -g @asyncapi/generator

COPY ./generate_docs.sh /

ENTRYPOINT ["/generate_docs.sh"]
