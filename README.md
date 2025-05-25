# Deployer

Deployer is a backend app that lets developers deploy their applications without worrying about Dockerfiles or deployment manifests. With a single API call, Deployer builds your code and pushes the resulting image to jFrog Artifactory, streamlining the deployment process for any project.

## Features

- **Automatic Docker Handling:** If your repo contains a Dockerfile, Deployer uses it. Otherwise, it auto-detects your application's language and builds the image accordingly.
- **Tekton Pipeline Integration:** When no Dockerfile is present, Deployer runs a Tekton pipeline tailored to your project's language.
- **Easy Deployment:** Deploy new versions with a simple POST request.
- **jFrog Artifactory Support:** Built images are pushed to your configured jFrog Artifactory repository.
- **Custom Application Naming:** Specify the application name as part of your deployment request.

## How It Works

1. Send a POST request to `/deploy` with the following JSON body:
   ```json
   {
     "name": "myapp",
     "repo_url": "https://github.com/user/project",
     "branch": "main",
     "instance": "dev"
   }
   ```
2. Deployer clones the repository and checks for a Dockerfile.
3. If a Dockerfile is found, it builds and pushes the image.
4. If not, Deployer identifies the main language and runs the appropriate Tekton pipeline before pushing the image.

## API

### `POST /deploy`

**Request Body:**
- `name` (string): Name of the application to deploy.
- `repo_url` (string): URL of the repository to deploy.
- `branch` (string): Branch to deploy from.
- `instance` (string): Target deployment environment (e.g., dev, prod).

**Example:**
```json
{
  "name": "payment-service",
  "repo_url": "https://github.com/user/myapp",
  "branch": "main",
  "instance": "production"
}
```

**Response:**
- Returns the build and deployment status, along with image information.

## Requirements

- Access to a jFrog Artifactory instance for image storage.
- Sufficient permissions for the target repositories.

## Getting Started

> _Instructions for installing/running locally, if applicable. (Please fill in as needed!)_

## Contributing

Contributions are welcome! Please open issues or submit pull requests for improvements and bug fixes.

## License

> _Specify license here, if applicable._

---

_Questions or feedback? Feel free to open an issue!_
