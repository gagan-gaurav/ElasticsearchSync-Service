# Fold assignment solution documentation

## Backend Architecture Solution
![image](https://github.com/gagan-gaurav/fold/assets/51356017/c72e4f5c-7b65-48bc-8fc4-f41237943917)

This architecture has 6 Components:
1. SQS FIFO Queue
2. Postgres Database
3. Backend Service (written in GO)
4. Lmanda Sync Service (written in Python)
5. Elasticsearch Service (Setup on EC2)
6. Search Service (written in GO)


## Step 1: SQS FIFO Queue
Create a SQS FIFO Queue from AWS console. Keep the queue_url

## Step 2: Setting up Postgres database (AWS RDS)
Setup a AWS RDS database with postgres engine. Keep the database_url

## Step 3: Building Backend Service
Clone this fold repo
```bash
git clone https://github.com/gagan-gaurav/fold.git
```

Build docker image (make sure docker daemon is running)
```bash
docker build -t backendservice:v1 .
```

Create an .env file in current folder with following environment variables:
```
POSTGRES_URL=your-postgres-database-url-from-step-1
AWS_ACCESS_KEY_ID=your-aws-access-key-id
AWS_SECRET_ACCESS_KEY=your-aws-secret-access-key
SQS_QUEUE_URL=your-sqs-queue-url-from-step-1 
```

Running docker container 
```bash
docker run -p 8080:8080 --env-file .env backendservice:v1
```

### API Documentation for Bacend Service

The following routes are available for interacting with the API:

| Route                          | Method | Description               |
|-------------------------------|--------|---------------------------|
| `/users`                       | POST   | Create user               |
| `/users/{id}`                  | GET    | Get user                  |
| `/users/update/{id}`           | POST   | Update user               |
| `/users/delete/{id}`           | DELETE | Delete user               |
| `/hashtags`                    | POST   | Create hashtag            |
| `/hashtags/{id}`               | GET    | Get hashtag               |
| `/hashtags/update/{id}`        | POST   | Update hashtag            |
| `/hashtags/delete/{id}`        | DELETE | Delete hashtag            |
| `/projects`                    | POST   | Create project            |
| `/projects/{id}`               | GET    | Get project               |
| `/projects/update/{id}`        | POST   | Update project            |
| `/projects/delete/{id}`        | DELETE | Delete project            |

Each route is associated with a specific HTTP method and provides functionality related to creating, retrieving, updating, or deleting users, hashtags, and projects.

Make sure to use the appropriate HTTP method and route to perform the desired action on the API.

## Step 4: Setting Elasticsearch on EC2

Create an EC2 instance. SSH into it. (Make sure to your security group has ingress for port:22 and port:9200)
Install Elasticsearch RPM package -> Follow this [instructions](https://www.elastic.co/guide/en/elasticsearch/reference/current/rpm.html) for installation.
Grab your username(default_username="elastic") and password.

## Step 5: Setting Lambda to intercept and process Queue data.

Create a lambda function with **python** runtime and necessary permissions and set your queue (created in step 1) as trigger.
link to lambdahandler code: https://github.com/gagan-gaurav/fold/blob/main/internal/services/lambdahandler.py

Setup environment variables for your lambda through AWS CLI
```
ELASTICSEARCH_PASSWORD = your-elasticsearch-password
ELASTICSEARCH_USERNAME = your-elasticsearch-username
```

## Step 6: Building Search Service
Clone searchService repo.
```bash
git clone https://github.com/gagan-gaurav/searchService.git
```
Build docker image (make sure docker daemon is running)
```bash
docker build -t searchservice:v1 .
```

Create an .env file in current folder with following environment variables:
```
ELASTICSEARCH_USERNAME=your-elasticsearch-username
ELASTICSEARCH_PASSWORD=your-elasticsearch-password
ELASTICSEARCH_URL=https:your-elasticsearch-url
```

Running docker container 
```bash
docker run -p 8080:8080 --env-file .env searchservice:v1
```
### API Documentation for Search Service

The following routes are defined for serarch service:

| Route               | Method | Description                              |
|-----------------------|--------|----------------------------------------|
| `/users`              | GET    | Users search with query parameter      |
| `/hashtags`           | GET    | Hashtags search with query parameter   |
| `/fuzzy`              | GET    | Fuzzy search with query parameter      |

Each route is associated with a specific HTTP method and provides functionality related to searching users, hashtags, and performing fuzzy searches.

Make sure to use the appropriate HTTP method and route along with the required query parameter to perform the desired search operation.







