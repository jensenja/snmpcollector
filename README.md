# snmpcollector

## Run from master
If you want to build a package yourself, or contribute. Here is a guide for how to do that.

### Dependencies

- Go 1.5
- NodeJS >=6.2.1

### Get Code

```bash
go get github.com/toni-moreno/snmpcollector
```

### Building the backend


```bash
cd $GOPATH/src/github.com/toni-moreno/snmpcollector
go run build.go setup            (only needed once to install godep)
godep restore                    (will pull down all golang lib dependencies in your current GOPATH)
```

### Building frontend and backend in production mode

```bash
npm install
PATH=$(npm bin):$PATH
npm run build:pro #will build fronted and backend
```
### Creating minimal package tar.gz

```bash
npm run postbuild #will build fronted and backend
```

### Running first time
To execute without any configuration you need a minimal config.toml file on the conf directory.

```bash
cp conf/sample.config.toml conf/config.toml
./bin/snmpcollector
```

### Recompile backend on source change (only for developers)

To rebuild on source change (requires that you executed godep restore)
```bash
go get github.com/Unknwon/bra
npm start
```
will init a change autodetect webserver with angular-cli (ng serve) and also a autodetect and recompile process with bra for the backend


#### Online config

Now you wil be able to configure metrics/measuremnets and devices from the builting web server at  http://localhost:8090 or http://localhost:4390 if working in development mode (npm start)

### Offline configuration.

You will be able also insert data directly on the sqlite db that snmpcollector has been created at first execution on config/snmpcollector.db examples on example_config.sql

```
cat conf/example_config.sql |sqlite3 conf/snmpcollector.db
```
