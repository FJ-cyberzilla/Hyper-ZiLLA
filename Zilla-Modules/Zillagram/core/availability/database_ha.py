# core/availability/database_ha.py
class DatabaseHighAvailability:
    """High availability configuration for databases with automatic failover"""
    
    def __init__(self, config: Dict):
        self.config = config
        self.postgres_ha = PostgresHAManager(config)
        self.redis_ha = RedisHAManager(config)
        
    async def setup_database_ha(self) -> Dict:
        """Setup high availability for all databases"""
        ha_configuration = {
            "postgresql_cluster": await self.postgres_ha.setup_ha_cluster(),
            "redis_sentinel": await self.redis_ha.setup_sentinel_cluster(),
            "connection_pooling": await self._configure_connection_pooling(),
            "automatic_failover": await self._configure_automatic_failover(),
            "backup_synchronization": await self._configure_backup_sync()
        }
        
        return ha_configuration

class PostgresHAManager:
    """PostgreSQL high availability with Patroni or Stolon"""
    
    async def setup_ha_cluster(self) -> Dict:
        """Setup PostgreSQL HA cluster"""
        patroni_config = {
            "apiVersion": "apps/v1",
            "kind": "StatefulSet",
            "metadata": {"name": "postgres-ha"},
            "spec": {
                "serviceName": "postgres-ha",
                "replicas": 3,
                "selector": {"matchLabels": {"app": "postgres-ha"}},
                "template": {
                    "metadata": {"labels": {"app": "postgres-ha"}},
                    "spec": {
                        "containers": [{
                            "name": "patroni",
                            "image": "patroni:latest",
                            "env": [
                                {"name": "PATRONI_NAME", "valueFrom": {"fieldRef": {"fieldPath": "metadata.name"}}},
                                {"name": "PATRONI_NAMESPACE", "valueFrom": {"fieldRef": {"fieldPath": "metadata.namespace"}}},
                                {"name": "PATRONI_SCOPE", "value": "postgres-ha"},
                                {"name": "PATRONI_POSTGRESQL_CONNECT_ADDRESS", "value": "$(PATRONI_NAME).postgres-ha:5432"},
                                {"name": "PATRONI_RESTAPI_CONNECT_ADDRESS", "value": "$(PATRONI_NAME).postgres-ha:8008"}
                            ],
                            "ports": [
                                {"name": "postgres", "containerPort": 5432},
                                {"name": "api", "containerPort": 8008}
                            ]
                        }],
                        "affinity": {
                            "podAntiAffinity": {
                                "requiredDuringSchedulingIgnoredDuringExecution": [{
                                    "labelSelector": {
                                        "matchExpressions": [{
                                            "key": "app",
                                            "operator": "In", 
                                            "values": ["postgres-ha"]
                                        }]
                                    },
                                    "topologyKey": "kubernetes.io/hostname"
                                }]
                            }
                        }
                    }
                }
            }
        }
        
        # Create Patroni cluster
        statefulset = await self._create_patroni_cluster(patroni_config)
        
        return {
            "cluster_type": "patroni",
            "replicas": 3,
            "failover_automation": True,
            "status": "active"
        }
