# gcp-magic-sql-proxy

This tool is a very simple wrapper for gce-proxy (gcr.io/cloudsql-docker/gce-proxy:1.16)

Based on the following env vars it gets a DB_SOCKET environment. With that CloudSQL connection it starts a cloud_sql_proxy allowing to connect to localhost 3306.

`CR_SERVICE_NAME`
`GCP_PROJECT`

`REGION`
- CR_DB_ENV_NAME, defaulting to `DB_SOCKET`

We use this in the pipeline like this:

```yaml
deploy:test:db:
  stage: deploy
  image:
    name: ${BASE_REPO}/api:${CI_PIPELINE_ID}
    entrypoint: ["bash"]
  services:
    - name: thijsdev/gcp-magic-sql-proxy:latest
      alias: sql-proxy
      entrypoint:
        - /gcp-magic-sql-proxy
  variables:
    DB_HOST: sql-proxy
    DB_USERNAME: cloudsqlproxy
    DB_PASSWORD: ""
    DB_DATABASE: template-app
    REGION: europe-west1
    CR_SERVICE_NAME: template-app-api
    GCP_PROJECT: phoenix-central-test
  script:
    - dockerize -wait tcp://sql-proxy:3306 -timeout 1m
    - /opt/application/artisan migrate --force
```

It allows us to create different databases while still use the same provision script. 