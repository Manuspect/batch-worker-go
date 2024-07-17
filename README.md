# Data flow

1. Get data-message from queue (NATS)

    - userID

    - recordId

    - timestamp (milis)

    - filePath

2. Download data (batch) from S3 using filePath

3. Unpack batch

4. Convert video to images

5. Upload images to S3

6. Send requests to processing service (each requestr consint of image and meta data)

    - get mata data from `.csv` file

7. Send acknowledge to queue if OK

8. Send to another queue processed files