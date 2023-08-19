# fold
fold assignment solution

## Backend Architecture Solution
![image](https://github.com/gagan-gaurav/fold/assets/51356017/c72e4f5c-7b65-48bc-8fc4-f41237943917)

So it has 6 Components:
1. SQS FIFO Queue
2. Postgres Database
3. Backend Service (written in GO)
4. Lmanda Sync Service (written in Python)
5. Elasticsearch Service (Setup on EC2)
6. Search Service (written in GO)


# Step 1: SQS FIFO Queue
Create a SQS FIFO Queue from AWS console. Keep the queue_url

# Step 2: Setting up Postgres database (AWS RDS)
Setup a AWS RDS database with postgres engine. Keep the database_url

# Step 3: Building Backend Service
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

# Step 4: Setting Lambda to intercept and process Queue data.
Create a lambda function with **python** runtime and necessary permissions and set your queue (created in step 1) as trigger.
link to lambdahandler code: 



