# go-mono-repo-aws-template

This is a golang mono-repo project template for AWS lambda based services.

# **Project Structure**

- **`infra`**  
  This directory contains deployment scripts for the project. These scripts are typically written using tools like **Terraform** or **AWS CloudFormation** to automate the deployment of project infrastructure.

- **`pkg`**  
  This directory is used to store shared code that is utilized across multiple services. Common examples include utility functions, shared libraries, or modules.

- **`services`**  
  This directory contains all the microservices, organized by domain. Each microservice is deployed as an AWS Lambda service and provides AWS Lambda handlers.  
  Each service is structured as follows:

    - **`application`**  
      This layer serves as the entry point for the service and defines its main use cases.
        - **`usecase`**: This subdirectory implements use cases, with each use case typically corresponding to one API endpoint or specific functionality.

    - **`domain`**  
      This is the domain layer where all the business logic of the service is implemented. It encapsulates the core functionality and ensures separation from infrastructure-related concerns.

    - **`infrastructure`**  
      This directory contains implementations related to infrastructure, such as the repository layer.
        - For example, the repository implementation can be tailored to use **MongoDB**, **MySQL**, or any other specific database.
        - It allows for flexibility in selecting and integrating infrastructure components based on the service requirements.



# Building

To build project:

```
make clean build
```

# Deploying

To deploy project

```
make deploy
```