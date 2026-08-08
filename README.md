# Job Assistant

AI job analysis assistant built with Go, go-zero, GORM, MySQL, and an OpenAI-compatible large model API.

## Features

- Analyze a job description with an AI model.
- Return structured API responses through go-zero.
- Persist analysis records to MySQL with GORM.
- Keep local API keys out of source control through environment variables.

## Tech Stack

- Go
- go-zero / goctl
- GORM
- MySQL
- OpenAI-compatible Chat Completions API

## Getting Started

### 1. Create the database

```sql
CREATE DATABASE job_analyzer DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_unicode_ci;

CREATE TABLE analyze_records (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    job_description LONGTEXT NOT NULL,
    question TEXT NOT NULL,
    answer LONGTEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 2. Create local config

```powershell
Copy-Item .\etc\job-api.example.yaml .\etc\job-api.yaml
```

Then update `etc/job-api.yaml` with your local MySQL DSN and model name.

### 3. Set API key

```powershell
$env:LLM_API_KEY="your-api-key"
```

### 4. Run the service

```powershell
go run . -f etc/job-api.yaml
```

### 5. Test the API

```powershell
$body = @{
    jobDescription = "Go backend developer, requires MySQL and RESTful API"
    question = "What skills does this job require?"
} | ConvertTo-Json

Invoke-RestMethod `
    -Uri "http://127.0.0.1:8888/api/analyze" `
    -Method Post `
    -ContentType "application/json; charset=utf-8" `
    -Body ([System.Text.Encoding]::UTF8.GetBytes($body))
```

## API

### `POST /api/analyze`

Request:

```json
{
  "jobDescription": "Go backend developer, requires MySQL and RESTful API",
  "question": "What skills does this job require?"
}
```

Response:

```json
{
  "code": 200,
  "message": "success",
  "success": true,
  "data": {
    "id": 1,
    "answer": "..."
  }
}
```
