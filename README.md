# Deployer

Deployer is a backend app with two microservices applier and builder, that lets developers deploy their applications without worrying about Dockerfiles or deployment manifests. With a single API call, Deployer builds your code and pushes the resulting image to jFrog Artifactory, streamlining the deployment process for any project.

## Features

- **Automatic Docker Handling:** If your repo contains a Dockerfile, Deployer uses it. Otherwise, it auto-detects your application's language and builds the image accordingly.
- **Tekton Pipeline Integration:** When no Dockerfile is present, Deployer runs a Tekton pipeline tailored to your project's language.
- **Easy Deployment:** Deploy new versions with a simple POST request.
- **jFrog Artifactory Support:** Built images are pushed to your configured jFrog Artifactory repository.
- **Custom Application Naming:** Specify the application name as part of your deployment request.

## How It Works

1. Send a POST request to `/build` with the following JSON body:
   ```json
   {
     "name": "myapp",
     "repo_url": "https://github.com/user/project",
     "branch": "main",
     "instance": "dev"
   }
   ```
2. Builder clones the repository and checks for a Dockerfile.
3. If a Dockerfile is found, it builds and pushes the image.
4. If not, Builder identifies the main language and runs the appropriate Tekton pipeline before pushing the image.
5. Now, if the image is built and pushed successfully, Applier ms comes into place.
6. It takes namespace name, appname, image tage and deploys the application with deployment, service and route/ingress.

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

- This application needs to be running on any k8s cluster.
- A role is specified in manifests/role.yaml with required permissions to start the pipelines and binded it with service account in manifests/tekton-sa.yaml in role-binding.yaml.
- The tekton task build-image needs jFrog secret to be mounted as volume. So jFrog secret is created as k8s secret using --from-file=config.json.
- This go application image is in dockerHub so dockerHub creds are added as imagePullSecret named regcred to the same service account and attached in deployment.yaml
- Dockerhub secret regcred attached to serviceAccount is for the deployment to pull the deployer application image from dockerHub and jFrog secret is created and mounted to tekton task to build the user 
  application images to jFrog artifactory.

  This is initial working pack of part 1 of the entire deployer app.

Phase 2 of part 1 (In development)
- The requests are tracked in DB, so that if user updates the code, if the same req body is triggered, image tag will be updated accordingly.
- GET request endpoints to track the entire status of the workflow.

Part 2 (In development)

- A dedicated go microservice will be places in all the k8s cluster could be EKS, Openshift etc. 
- This microservice will be responsible for creating the manifests like deployments, imagepullsecrets, services, routes etc., for the app.
- This microservice will poll the API server (built in phase 1) for any newly created/ updated image requests.

05/25 Phase 1 of part 2

- Miroservice 'applier' runs with an endpoint which takes namespace name, iamge tag, application name.
- It then ensures if given namespace exists, if not creates.
- It copies the jFrog artifactory secret from 'default' (to be applied manually before the 1st launch of application) to the user given namespace.
- Creates deployment with imagepullSecret attached, service of type ClusterIP and basic ingress (this can be later changed to route for openshift).
- It does all this using the service account given in applier/manifests/applier-sa.yaml. A role and role binding also declared there. 


Phase 2 of part 2 (In development)
- Monitor deployment until scaled as expected.
<To be added>


## Contributing

Contributions are welcome! Please open issues or submit pull requests for improvements and bug fixes.

## License

> _Specify license here, if applicable._

---

_Questions or feedback? Feel free to open an issue!_
