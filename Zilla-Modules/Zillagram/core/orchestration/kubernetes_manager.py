# core/orchestration/kubernetes_manager.py
import asyncio
import yaml
from typing import Dict, List, Optional
from kubernetes import client, config
from kubernetes.client.rest import ApiException
import subprocess
import logging

class KubernetesOrchestrator:
    """
    Enterprise Kubernetes orchestration for container deployment and management
    """
    
    def __init__(self, k8s_config: Dict):
        self.k8s_config = k8s_config
        self.load_kube_config()
        self.apps_v1 = client.AppsV1Api()
        self.core_v1 = client.CoreV1Api()
        self.networking_v1 = client.NetworkingV1Api()
        self.monitoring = KubernetesMonitor(self)
        
    def load_kube_config(self):
        """Load Kubernetes configuration based on environment"""
        try:
            if self.k8s_config['in_cluster']:
                config.load_incluster_config()
            else:
                config.load_kube_config(context=self.k8s_config.get('context'))
        except Exception as e:
            logging.error(f"Failed to load kube config: {e}")
            raise
    
    async def deploy_enterprise_stack(self) -> Dict:
        """Deploy complete enterprise stack to Kubernetes"""
        deployments = {
            "api_gateway": await self.deploy_api_gateway(),
            "worker_nodes": await self.deploy_worker_pool(),
            "database_cluster": await self.deploy_database(),
            "monitoring_stack": await self.deploy_monitoring(),
            "security_services": await self.deploy_security_services(),
            "data_lake": await self.deploy_data_lake()
        }
        
        # Wait for all deployments to be ready
        await self.wait_for_deployments_ready(list(deployments.keys()))
        
        logging.info("✅ Enterprise stack deployed successfully")
        return deployments
    
    async def deploy_api_gateway(self) -> Dict:
        """Deploy API Gateway with security and monitoring"""
        api_gateway_manifest = {
            "apiVersion": "apps/v1",
            "kind": "Deployment",
            "metadata": {"name": "osint-api-gateway"},
            "spec": {
                "replicas": 3,
                "selector": {"matchLabels": {"app": "osint-api-gateway"}},
                "template": {
                    "metadata": {
                        "labels": {"app": "osint-api-gateway"},
                        "annotations": {
                            "prometheus.io/scrape": "true",
                            "prometheus.io/port": "9090"
                        }
                    },
                    "spec": {
                        "containers": [{
                            "name": "api-gateway",
                            "image": f"{self.k8s_config['registry']}/osint-api-gateway:latest",
                            "ports": [{"containerPort": 8080}],
                            "env": self._get_api_gateway_env(),
                            "resources": {
                                "requests": {"cpu": "200m", "memory": "512Mi"},
                                "limits": {"cpu": "1000m", "memory": "2Gi"}
                            },
                            "livenessProbe": {
                                "httpGet": {"path": "/health", "port": 8080},
                                "initialDelaySeconds": 30,
                                "periodSeconds": 10
                            },
                            "readinessProbe": {
                                "httpGet": {"path": "/ready", "port": 8080},
                                "initialDelaySeconds": 5,
                                "periodSeconds": 5
                            }
                        }],
                        "securityContext": {
                            "runAsNonRoot": True,
                            "runAsUser": 1000
                        }
                    }
                }
            }
        }
        
        # Create deployment
        api_response = self.apps_v1.create_namespaced_deployment(
            namespace=self.k8s_config['namespace'],
            body=api_gateway_manifest
        )
        
        # Create service
        service_manifest = {
            "apiVersion": "v1",
            "kind": "Service",
            "metadata": {"name": "osint-api-gateway"},
            "spec": {
                "selector": {"app": "osint-api-gateway"},
                "ports": [{"port": 80, "targetPort": 8080}],
                "type": "LoadBalancer"
            }
        }
        
        service_response = self.core_v1.create_namespaced_service(
            namespace=self.k8s_config['namespace'],
            body=service_manifest
        )
        
        return {
            "deployment": api_response.metadata.name,
            "service": service_response.metadata.name,
            "status": "deployed"
        }
    
    async def deploy_worker_pool(self) -> Dict:
        """Deploy scalable worker pool for processing"""
        worker_manifest = {
            "apiVersion": "apps/v1",
            "kind": "Deployment",
            "metadata": {"name": "osint-workers"},
            "spec": {
                "replicas": 5,
                "selector": {"matchLabels": {"app": "osint-worker"}},
                "template": {
                    "metadata": {"labels": {"app": "osint-worker"}},
                    "spec": {
                        "containers": [{
                            "name": "worker",
                            "image": f"{self.k8s_config['registry']}/osint-worker:latest",
                            "env": self._get_worker_env(),
                            "resources": {
                                "requests": {"cpu": "500m", "memory": "1Gi"},
                                "limits": {"cpu": "2000m", "memory": "4Gi"}
                            },
                            "volumeMounts": [{
                                "name": "worker-secrets",
                                "mountPath": "/etc/secrets",
                                "readOnly": True
                            }]
                        }],
                        "volumes": [{
                            "name": "worker-secrets",
                            "secret": {"secretName": "worker-credentials"}
                        }],
                        "affinity": {
                            "podAntiAffinity": {
                                "preferredDuringSchedulingIgnoredDuringExecution": [{
                                    "weight": 100,
                                    "podAffinityTerm": {
                                        "labelSelector": {
                                            "matchExpressions": [{
                                                "key": "app",
                                                "operator": "In",
                                                "values": ["osint-worker"]
                                            }]
                                        },
                                        "topologyKey": "kubernetes.io/hostname"
                                    }
                                }]
                            }
                        }
                    }
                }
            }
        }
        
        # Create Horizontal Pod Autoscaler
        hpa_manifest = {
            "apiVersion": "autoscaling/v2",
            "kind": "HorizontalPodAutoscaler",
            "metadata": {"name": "osint-workers-hpa"},
            "spec": {
                "scaleTargetRef": {
                    "apiVersion": "apps/v1",
                    "kind": "Deployment",
                    "name": "osint-workers"
                },
                "minReplicas": 3,
                "maxReplicas": 20,
                "metrics": [{
                    "type": "Resource",
                    "resource": {
                        "name": "cpu",
                        "target": {"type": "Utilization", "averageUtilization": 70}
                    }
                }]
            }
        }
        
        deployment = self.apps_v1.create_namespaced_deployment(
            namespace=self.k8s_config['namespace'],
            body=worker_manifest
        )
        
        hpa = self.autoscaling_v2.create_namespaced_horizontal_pod_autoscaler(
            namespace=self.k8s_config['namespace'],
            body=hpa_manifest
        )
        
        return {
            "deployment": deployment.metadata.name,
            "hpa": hpa.metadata.name,
            "replicas": "3-20",
            "status": "deployed"
        }
    
    async def deploy_database(self) -> Dict:
        """Deploy highly available database cluster"""
        # StatefulSet for database with persistent storage
        db_manifest = {
            "apiVersion": "apps/v1",
            "kind": "StatefulSet",
            "metadata": {"name": "osint-database"},
            "spec": {
                "serviceName": "osint-database",
                "replicas": 3,
                "selector": {"matchLabels": {"app": "osint-database"}},
                "template": {
                    "metadata": {"labels": {"app": "osint-database"}},
                    "spec": {
                        "containers": [{
                            "name": "postgres",
                            "image": "postgres:14",
                            "ports": [{"containerPort": 5432}],
                            "env": [
                                {"name": "POSTGRES_DB", "value": "osint_platform"},
                                {"name": "POSTGRES_USER", "valueFrom": {"secretKeyRef": {"name": "db-secrets", "key": "username"}}},
                                {"name": "POSTGRES_PASSWORD", "valueFrom": {"secretKeyRef": {"name": "db-secrets", "key": "password"}}}
                            ],
                            "volumeMounts": [{"name": "db-data", "mountPath": "/var/lib/postgresql/data"}],
                            "resources": {
                                "requests": {"cpu": "1000m", "memory": "2Gi"},
                                "limits": {"cpu": "4000m", "memory": "8Gi"}
                            }
                        }],
                        "securityContext": {
                            "fsGroup": 999,
                            "runAsUser": 999
                        }
                    }
                },
                "volumeClaimTemplates": [{
                    "metadata": {"name": "db-data"},
                    "spec": {
                        "accessModes": ["ReadWriteOnce"],
                        "storageClassName": "fast-ssd",
                        "resources": {"requests": {"storage": "100Gi"}}
                    }
                }]
            }
        }
        
        statefulset = self.apps_v1.create_namespaced_stateful_set(
            namespace=self.k8s_config['namespace'],
            body=db_manifest
        )
        
        return {
            "statefulset": statefulset.metadata.name,
            "replicas": 3,
            "storage": "100Gi per replica",
            "status": "deployed"
        }
