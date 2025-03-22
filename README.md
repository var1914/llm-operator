# LLM Operator

A Kubernetes Operator for deploying and managing Large Language Models on Kubernetes.

## Overview

This project provides a Kubernetes Operator that allows you to deploy and manage LLM models in a Kubernetes-native way. It defines two Custom Resource Definitions (CRDs):

- **LLMModel**: Defines the model specifications
- **LLMDeployment**: Handles the actual deployment of models with configurable replicas and resource settings

## Prerequisites

- Go 1.19+
- Kubernetes cluster v1.22+
- kubectl
- Operator SDK v1.28.0+
- Docker (for building images)
- Helm v3+

## Installation

### 1. Install Required Tools

### 2. Clone the Repository

```bash
git clone https://github.com/yourusername/llm-operator.git
cd llm-operator
```

### 3. Deploy the Operator

#### Option 1: Using Pre-built Images

```bash
# Create namespace
kubectl create namespace llm-system

# Install CRDs
kubectl apply -f helm-charts/llm-operator/crds/

# Install the operator
helm install llm-operator ./helm-charts/llm-operator --namespace llm-system
```

#### Option 2: Building From Source

```bash
# Build and push the operator image (replace with your registry)
make docker-build docker-push IMG=myregistry/llm-operator:v0.1.0

# Update the values.yaml file with your image
nano helm-charts/llm-operator/values.yaml

# Install CRDs
kubectl apply -f helm-charts/llm-operator/crds/

# Install the operator
helm install llm-operator ./helm-charts/llm-operator --namespace llm-system
```

## Usage

### 1. Create an LLMModel

```yaml
apiVersion: llm.example.com/v1alpha1
kind: LLMModel
metadata:
  name: gpt-small
spec:
  modelName: "GPT-Small"
  image: "your-registry/gpt-small:v1"
  resources:
    cpu: "2"
    memory: "4Gi"
```

Save this as `model.yaml` and apply:

```bash
kubectl apply -f model.yaml
```

### 2. Create an LLMDeployment

```yaml
apiVersion: llm.example.com/v1alpha1
kind: LLMDeployment
metadata:
  name: chat-service
spec:
  modelRef: "gpt-small"
  replicas: 2
  port: 8080
```

Save this as `deployment.yaml` and apply:

```bash
kubectl apply -f deployment.yaml
```

### 3. Verify the Deployment

```bash
# Check LLM custom resources
kubectl get llmmodels
kubectl get llmdeployments

# Check the created Kubernetes resources
kubectl get deployments
kubectl get services
```

## Project Structure

```
llm-operator/
├── api/
│   └── v1alpha1/             # API definitions for CRDs
│       ├── llmmodel_types.go
│       └── llmdeployment_types.go
├── controllers/              # Controller implementations
│   ├── llmmodel_controller.go
│   └── llmdeployment_controller.go
├── config/                   # Generated configuration
│   ├── crd/                  # Generated CRDs
│   ├── rbac/                 # RBAC configurations
│   └── manager/              # Manager configurations
├── helm-charts/              # Helm charts for deployment
│   └── llm-operator/
│       ├── crds/             # CRDs to be installed
│       ├── templates/        # Helm templates
│       ├── Chart.yaml
│       └── values.yaml
└── main.go                   # Operator entrypoint
```

## Developing the Operator

### 1. Create a New Operator

```bash
# Initialize a new operator project
operator-sdk init --domain example.com --repo github.com/yourusername/llm-operator

# Create APIs
operator-sdk create api --group llm --version v1alpha1 --kind LLMModel --resource --controller
operator-sdk create api --group llm --version v1alpha1 --kind LLMDeployment --resource --controller
```

### 2. Customize the CRDs

Modify the files in `api/v1alpha1/` to define your custom resources.

### 3. Implement Controllers

Implement the reconciliation logic in the controller files in the `controllers/` directory.

### 4. Build and Deploy

```bash
# Generate manifests
make manifests

# Build and push the image
make docker-build docker-push IMG=myregistry/llm-operator:v0.1.0

# Deploy
make deploy IMG=myregistry/llm-operator:v0.1.0
```

## Custom Resource Definitions

### LLMModel

The `LLMModel` CRD defines the model to be deployed:

| Field | Description | Required |
|-------|-------------|----------|
| `spec.modelName` | Name of the LLM model | Yes |
| `spec.image` | Container image for the model | Yes |
| `spec.resources.cpu` | CPU resource requirements | No |
| `spec.resources.memory` | Memory resource requirements | No |

### LLMDeployment

The `LLMDeployment` CRD defines how to deploy a model:

| Field | Description | Required | Default |
|-------|-------------|----------|---------|
| `spec.modelRef` | Reference to an LLMModel resource | Yes | - |
| `spec.replicas` | Number of replicas to deploy | Yes | - |
| `spec.port` | Port the model service will listen on | No | 8080 |

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
