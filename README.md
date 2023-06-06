# tanzu-go-helloworld

## Bindings

This application requires a binding of type `postgresql`.

```
tanzu service class list
  NAME                  DESCRIPTION
  gcp-cloudsql-psql     Google Cloud Platform PostgreSQL
  mysql-unmanaged       MySQL by Bitnami
  postgresql-unmanaged  PostgreSQL by Bitnami
  rabbitmq-unmanaged    RabbitMQ by Bitnami
  redis-unmanaged       Redis by Bitnami
```

```
tanzu service class-claim create pgdb-1 --class gcp-cloudsql-psql -p parameters.storageGB=20 -n dev
```
