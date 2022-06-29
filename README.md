# lift
Lift is a service that lifts infrastructure on demand, more specifically code analysis tools.

## Supported Tools

Currently lift supports the following analysis tools:
- Sonarqube
- Jenkins

## Getting started

In order for the application to work on a machine there are a few things that need to be in place.

### Requirements

There are a set of requirements that need to be fullfilled in order for the service to work:
- Go installed
- Make installed
- Docker installed
- Docker Compose installed
- A private and public SSH key located in `lift/static/keys` called lift and lift.pub
- An AWS account (note that non free instances will be used on AWS)
- A `.env` file in the root of the project containing the following:
  - AWS_ACCESS_KEY_ID=< your-aws-access-key-id >
  - AWS_SECRET_ACCESS_KEY=< your-aws-secret-access-key >
  - DB_NAME=< database-name >
  - DB_USER=< database-username >
  - DB_PASSWORD=< database-user-password >
  - DB_ROOT_PASSWORD=< database-root-user-password >

### Starting the service

In order to run the service you can simply run `make start` which will handle the following:
- Having the GraphQL API listening on Port 8080
- Starting the Observer service in parallel
- Setting up the PostgreSQL database

### Functionalities

Lift currently supports the following:
- Creating and changing the configuration that will be used for the deployments
- Allows the execution of deployments to AWS
- Has automatic teardown of non functional instances
