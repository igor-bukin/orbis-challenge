## API project template

### Dependency Management

Use ``glide`` tool. Add to ``glide.yaml`` vcs links with particular versions or hash commits and run:

``make dep``

### Build

``make build``

### Run Tests

``make test``

### Run Linter checks

``make lint``

### Run migrations

``make migrate-up``

### Run it

Copy config sample

``cp config_sample.json config.json``

edit ``config.json`` if need then

``./bin/api``

or:

``go run src/main.go``