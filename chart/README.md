# NOFX Helm Chart

This Helm Chart is used to deploy NOFX (Open Source AI Trading OS) on Kubernetes clusters.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+

## Installation

### 1. Add the NOFX Helm Repository

```bash
# Add the repository (if not already added)
helm repo add nofx https://nofxai.github.io/helm-charts
helm repo update
```

### 2. Install the Chart

#### Basic Installation

```bash
# Install with default values
helm install nofx nofx/nofx --namespace nofx --create-namespace
```

#### Custom Installation

```bash
# Create a custom values file
cat > values.yaml << EOF
backend:
  env:
    TZ: Asia/Shanghai
    AI_MAX_TOKENS: 8000
    # Add other environment variables here
  persistence:
    enabled: true
    size: 20Gi

frontend:
  resources:
    requests:
      cpu: 500m
      memory: 512Mi
    limits:
      cpu: 2
      memory: 2Gi

ingress:
  enabled: true
  hosts:
    - host: nofx.example.com
      paths:
        - path: /
          pathType: Prefix
          backend:
            service:
              name: nofx-frontend
              port:
                number: 80
EOF

# Install with custom values
helm install nofx nofx/nofx --namespace nofx --create-namespace -f values.yaml
```

## Configuration

### Key Configuration Parameters

| Parameter | Description | Default Value |
|-----------|-------------|---------------|
| `backend.image.repository` | Backend image repository | `nofxai/nofx-backend` |
| `backend.image.tag` | Backend image tag | `latest` |
| `backend.service.port` | Backend service port | `8080` |
| `backend.persistence.enabled` | Enable persistent storage for backend | `true` |
| `backend.persistence.size` | Persistent storage size | `10Gi` |
| `frontend.image.repository` | Frontend image repository | `nofxai/nofx-frontend` |
| `frontend.image.tag` | Frontend image tag | `latest` |
| `frontend.service.port` | Frontend service port | `80` |
| `ingress.enabled` | Enable ingress | `false` |
| `ingress.hosts` | Ingress host configuration | See `values.yaml` |

### Environment Variables

The backend service supports the following environment variables:

| Environment Variable | Description | Default Value |
|----------------------|-------------|---------------|
| `TZ` | Timezone | `Asia/Shanghai` |
| `AI_MAX_TOKENS` | Maximum tokens for AI models | `8000` |
| `TRANSPORT_ENCRYPTION` | Enable transport encryption | `false` |
| `NOFX_BACKEND_PORT` | Backend service port | `8080` |

## Accessing the Application

### 1. Through Ingress (Recommended for Production)

If you enabled ingress, you can access the application through the configured hostname:

```bash
# Example: Access through https://nofx.example.com
echo "Access NOFX at: https://$(kubectl get ingress nofx -n nofx -o jsonpath='{.spec.rules[0].host}')"
```

### 2. Through Port Forwarding (For Development)

```bash
# Forward frontend port
kubectl port-forward svc/nofx-frontend 3000:80 -n nofx

# Access at http://localhost:3000
```

### 3. Through LoadBalancer Service (If configured)

```bash
# Get the LoadBalancer IP
kubectl get svc nofx-frontend -n nofx -o jsonpath='{.status.loadBalancer.ingress[0].ip}'

# Access at http://<LoadBalancer-IP>
```

## Upgrading

```bash
# Upgrade the chart
helm upgrade nofx nofx/nofx --namespace nofx -f values.yaml
```

## Uninstalling

```bash
# Uninstall the chart
helm uninstall nofx --namespace nofx

# Delete the namespace (optional)
kubectl delete namespace nofx
```

## Persistence

The backend service uses a persistent volume claim to store data. By default, a 10Gi volume is created. You can adjust the size in the `values.yaml` file.

## Monitoring and Logging

### Viewing Logs

```bash
# View backend logs
kubectl logs deployment/nofx-backend -n nofx -f

# View frontend logs
kubectl logs deployment/nofx-frontend -n nofx -f
```

### Health Checks

The chart includes health checks for both backend and frontend services. You can monitor their status using:

```bash
# Check pod status
kubectl get pods -n nofx

# Check pod details
kubectl describe pod <pod-name> -n nofx
```

## Troubleshooting

### Common Issues

1. **Backend service not starting**
   - Check if the persistent volume is correctly provisioned
   - Verify environment variables are set correctly
   - Check backend logs for error messages

2. **Frontend can't connect to backend**
   - Ensure the backend service is running and accessible
   - Check if the backend service name is correctly configured in the frontend

3. **Ingress not working**
   - Verify the ingress controller is installed in your cluster
   - Check the ingress configuration and DNS settings
   - Ensure the ingress resource is correctly created

### Support

For more help, please refer to the [NOFX documentation](https://github.com/NoFxAiOS/nofx/blob/main/README.md) or join the [NOFX Developer Community](https://t.me/nofx_dev_community).

## License

This Helm Chart is licensed under the AGPL-3.0 License. See the [LICENSE](https://github.com/NoFxAiOS/nofx/blob/main/LICENSE) file for details.
