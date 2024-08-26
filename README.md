# PowerPoint Parser Microservice

This Go-based microservice parses PowerPoint (`.pptx`) files and extracts slide content, including text and images. The images are returned as Base64-encoded strings within a JSON response. This microservice is designed to be run as a Docker container for easy deployment.

## Features

- Parses PowerPoint presentations (`.pptx`)
- Extracts slide content including titles, paragraphs, and images
- Encodes images as Base64 for inclusion in JSON responses
- Runs as a lightweight, portable Docker container

## Prerequisites

Before you begin, ensure you have the following installed on your machine:

- [Docker](https://www.docker.com/get-started)
- [Git](https://git-scm.com/)

## Getting Started

### 1. Clone the Repository

First, clone this repository to your local machine:

```bash
git clone https://github.com/yourusername/pptx-parser-microservice.git
cd pptx-parser-microservice
