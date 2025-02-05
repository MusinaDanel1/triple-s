# Triple-S: Simple Storage Service

**Triple-S** is a simplified implementation of a cloud storage service, similar to Amazon S3. It allows users to create buckets, upload, retrieve, and delete files (objects), and manage object metadata. This project demonstrates fundamental concepts of RESTful APIs, HTTP, basic networking, and object storage.

## Table of Contents
- [Project Overview](#project-overview)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
  - [Bucket Management](#bucket-management)
  - [Object Operations](#object-operations)
- [Directory Structure](#directory-structure)
- [Error Handling](#error-handling)
- [Metadata Storage](#metadata-storage)
- [Examples](#examples)
- [License](#license)

## Project Overview

The **Triple-S** project is a simplified version of Amazon S3, designed to help you understand the core principles behind cloud storage solutions. The system allows users to manage storage containers (buckets) and objects (files), offering the following functionalities:
- **Bucket Management**: Create, list, and delete storage buckets.
- **Object Operations**: Upload, retrieve, and delete objects within buckets.
- **RESTful API**: Interact with the system via HTTP-based API endpoints.
- **Metadata Storage**: Metadata for buckets and objects is stored in CSV files.
  
The system responds with XML format in compliance with Amazon S3's specifications.

## Features

- Create and manage storage buckets.
- Upload, retrieve, and delete files within those buckets.
- Simple REST API with XML responses.
- Metadata storage in CSV format for buckets and objects.

## Installation

### Prerequisites:
- Go (version 1.16 or higher)

### Steps to Set Up:
1. Clone the repository:
   git clone https://github.com/your-username/triple-s.git
2. Navigate to the project directory:
   cd triple-s
3. Initialize the Go module:
   go mod tidy
4. Build the project:
   go build -o triple-s .
   
##Usage
To run the server, use the following command:

./triple-s -port <port-number> -dir <storage-directory>

##Where:
-port <port-number> specifies the port the server will listen on (default: 8080).
-dir <storage-directory> specifies the path to the directory where the buckets and objects will be stored.

##Example:
To run the server on port 8080 with the storage directory at /path/to/storage:
./triple-s -port 8080 -dir /path/to/storage

##Show help:
./triple-s --help
This will display the available options for configuring the server.

#API Endpoints
Bucket Management

1. Create a Bucket:
HTTP Method: PUT
Endpoint: /:{BucketName}
Request Body: Empty
Response: 200 OK on success or error message.
List All Buckets:

2. HTTP Method: GET
Endpoint: /
Response: XML list of all buckets.

3. Delete a Bucket:
HTTP Method: DELETE
Endpoint: /:{BucketName}
Response: 204 No Content if successful, error message otherwise.

#Object Operations
1. Upload a New Object:
HTTP Method: PUT
Endpoint: /:{BucketName}/{ObjectKey}
Request Body: Binary content (file).
Response: 200 OK on success.

2. Retrieve an Object:
HTTP Method: GET
Endpoint: /:{BucketName}/{ObjectKey}
Response: Binary content of the object, appropriate MIME type.

3. Delete an Object:
HTTP Method: DELETE
Endpoint: /:{BucketName}/{ObjectKey}
Response: 204 No Content on success.

#Directory Structure
The project stores data in a data/ directory. The structure is as follows:
/data
  /{bucket-name}
    /objects.csv         # Metadata of objects in the bucket
    /{object-key}        # Stored object (file)
  /buckets.csv           # Metadata of all buckets

The objects.csv file stores metadata for objects, including their keys, sizes, and content types.
The buckets.csv file stores metadata for buckets, including names, creation times, and modification times.

#Error Handling
The server handles errors gracefully and returns appropriate HTTP status codes:
400 Bad Request: Invalid bucket or object name.
404 Not Found: Bucket or object does not exist.
409 Conflict: Bucket already exists or bucket is not empty when trying to delete.
500 Internal Server Error: Server errors (e.g., permission issues, file system errors).

#Metadata Storage
Bucket Metadata (buckets.csv)
Each line represents a bucket:
BucketName,CreationTime,LastModifiedTime,Status

Object Metadata (objects.csv)
Each line represents an object within a bucket:
ObjectKey,Size,ContentType,LastModified

#Examples

Example 1: Create a Bucket
Request:
PUT /photos
Response:
<Bucket>
  <Name>photos</Name>
  <CreationDate>2025-02-05T12:00:00Z</CreationDate>
</Bucket>

Example 2: Upload an Object
Request:
PUT /photos/sunset.png
Response:
<Object>
  <Key>sunset.png</Key>
  <Size>1024</Size>
  <ContentType>image/png</ContentType>
  <LastModified>2025-02-05T12:00:00Z</LastModified>
</Object>

Example 3: Delete an Object
Request:
DELETE /photos/sunset.png
Response:
<Status>204 No Content</Status>
