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

### CI-CD

use gophers/ci-cd repository for multiservice architecture

and will look like .gitlab-ci.yml :

```
image: gitlab.yalantis.com:4567/gophers/ci-cd/services/go-builder:latest

stages:
  - dependencies
  - test
  - go-build
  - build
  - deploy

include:
  - project: 'GROUP_NAME/ci-cd'
    file: '/go-templates/main.yml'  
    - project: 'GROUP_NAME/ci-cd'
    file: '/go-templates/dependencies.yml'  
    - project: 'GROUP_NAME/ci-cd'
    file: '/go-templates/tests.yml'  
    - project: 'GROUP_NAME/ci-cd'
    file: '/go-templates/linter.yml'  
    - project: 'GROUP_NAME/ci-cd'
    file: '/go-templates/build.yml'  
    - project: 'GROUP_NAME/ci-cd'
 ```