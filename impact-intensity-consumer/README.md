# Impact Intensity Consumer

Listens to an SQS queue for intensity messages and saves them to the impact DB.

## Message Format

## Deployment

### AWS Resources

There is an AWS Cloudformation template in the cloudformation directory to create the SQS queue and a user to access it.  Use the ouput of this template for the configuration in impact-intensity-consumer.json.  Subscribe a suitable endpoint to the alarm topic for the stack.

### Properties 

Either or both of: 
1. Copy an appropriately edited version of `impact-intensity-consumer.json` to `/etc/sysconfig/impact-intensity-consumer.json`  This should include write access credentials for accessing the impact database.
2. Refer to docker-run.sh for overriding from env var.

### Scaling

The number of concurrent SQS listeners can be can be controlled with the config parameter `NumberOfListeners`.  

### Database Fault Tolerance.

The app is tolerant of (most) database faults:

1. Won't fully start until the DB can be contacted.
2. Blocks processing if a message can't be stored in the DB.

Case 2. can be caused by the DB being uncontactable for a while in which case the app will block message processing until the database becomes available again and then recover.  If it is caused by schema or permission changes then app will block message processing until the problem is fixed.

### AWS SQS Fault Tolerance

The app is tolerant of networking faults for SQS.  This does mean that it will not exit for SQS config errors and will loop logging ERROR messages.