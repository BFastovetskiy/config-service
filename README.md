# Configuration & discovery service

Service for managing configuration files and discovering services. The service is oriented for use in systems with microservice architecture.

## 1. Launch

## 1.1. Launch parameters

Command line options

- **i** - Interactive mode the service setup. Allows you to configure the service in question-answer mode and save it to a configuration file.
- **listener-port** - Listener HTTP port
- **listener-ssl-port** - Listener SSL port
- **listener-cluster-port** - Listener cluster port
- **database-name** - Database name
- **only-ssl** - Listener only SSL
- **certificate-pem** - PEM file certificate
- **certificate-key** - Key file certificate
- **secret** - Secret for connection to service
- **frequency** - Frequency timeout of checking the availability of services. Seconds
- **join-to** - Address to join on cluster
- **debug** - Use http listener for profilers and debug

Configuration priority over **app.config** file. Example file

```yaml
application:
  srvPort: 3000
  srvSslPort: 3001
  clusterProt: 4000
  db: database.db
  pem: public.crt
  key: private.key
  onlySsl: false
  discovery:
    secret: EIaUPvdI1ONo6IQowmo6HsSRZBxUv4Hb
    frequency: 60
  cluster:
    join: :4000
```
