import json
import http.client
import base64
import ssl
import os

def lambda_handler(event, context):
    response_body = "sync failed."
    url = "/projects/_doc"
    username = os.environ['ELASTICSEARCH_USERNAME']
    password = os.environ['ELASTICSEARCH_PASSWORD']
    
    try:
        for record in event['Records']:
            print(record)
            # Get the message body from the SQS record
            message_body = json.loads(record['body']) # Extract the body field
            doc_id = message_body['doc']['id'] # Extract project id
            method = message_body['method'] # Extract method
            url += ("/" + str(doc_id))
            
            # Create an HTTP connection
            conn = http.client.HTTPSConnection("3.108.40.246", 9200, context=ssl._create_unverified_context())
            
            # Headers with basic authentication
            auth = f"{username}:{password}"
            auth_bytes = base64.b64encode(auth.encode())
            headers = {
                "Authorization": f"Basic {auth_bytes.decode()}",
                "Content-Type": "application/json"
            }
            
            # Send the request to index the data in Elasticsearch
            if(method == "POST"):
                conn.request("POST", url, body= json.dumps(message_body['doc']), headers=headers)  # Convert message_body to JSON
            else:
                conn.request("DELETE", url, headers=headers)
            
            response = conn.getresponse()
            
            if response.status == 200 or response.status == 201:
                if(method == "POST") :
                    print("Data indexed successfully!")
                else : 
                    print ("Data deleted successfully")
                response_body = "successfully synced the data on elasticsearch."
            else:
                print(f"Failed to index data. Response: {response.read().decode()}")
                
            
            conn.close()  # Close the connection
    
    except Exception as e:
        print(f"An error occurred: {str(e)}")
    
    return {
        'statusCode': 200,
        'body': json.dumps(response_body)
    }
