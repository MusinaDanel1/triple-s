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
4.Build the project:
   go build -o triple-s .
