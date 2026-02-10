# core/network/dns_manager.py
import aiodns
import asyncio
from typing import List, Dict, Any
import time
import logging
import random

logger = logging.getLogger(__name__)

class EnterpriseDNSManager:
    def __init__(self):
        self.resolver = aiodns.DNSResolver()
        self.dns_cache = {}
        self.dns_servers = [
            '8.8.8.8',  # Google
            '1.1.1.1',  # Cloudflare
            '9.9.9.9',  # Quad9
            '208.67.222.222'  # OpenDNS
        ]
        self.performance_stats = {server: {"reliability": 1.0, "latency": 0.0} for server in self.dns_servers}
        
    async def _resolve_with_server(self, domain: str, query_type: str, dns_server: str) -> Dict:
        """Simulate DNS resolution with a specific server."""
        logger.debug(f"Attempting to resolve {domain} ({query_type}) via {dns_server}")
        start_time = time.time()
        # Simulate network delay and potential failure
        await asyncio.sleep(random.uniform(0.05, 0.5))
        
        if random.random() < 0.1: # 10% chance of failure
            raise aiodns.error.DNSError("Simulated DNS resolution failure")

        # Simulate A record response
        if query_type == 'A':
            results = [{'host': f'192.0.2.{random.randint(1, 254)}', 'ttl': random.randint(60, 3600)}]
        else:
            results = [{'host': 'stub_result', 'ttl': random.randint(60, 3600)}]

        latency = (time.time() - start_time) * 1000 # milliseconds
        return {
            'domain': domain,
            'query_type': query_type,
            'results': results,
            'response_time': latency,
            'source': dns_server,
            'ttl': results[0]['ttl'] # Assuming at least one result
        }

    def _update_dns_performance(self, dns_server: str, success: bool):
        if dns_server not in self.performance_stats:
            self.performance_stats[dns_server] = {"reliability": 1.0, "latency": 0.0}
        
        if success:
            self.performance_stats[dns_server]["reliability"] = min(1.0, self.performance_stats[dns_server]["reliability"] + 0.05)
        else:
            self.performance_stats[dns_server]["reliability"] = max(0.0, self.performance_stats[dns_server]["reliability"] - 0.1)
        # Latency update would require passing actual latency from _resolve_with_server

    def _calculate_geo_relevance(self, result: Dict, domain: str) -> float:
        # Placeholder for geographic relevance calculation
        return 0.7 # Assume moderately relevant
        
    async def smart_dns_resolution(self, domain: str, query_type: str = 'A') -> Dict:
        """AI-driven DNS resolution with fallback and optimization"""
        # Check cache first
        cache_key = f"{domain}_{query_type}"
        if cache_key in self.dns_cache:
            cached = self.dns_cache[cache_key]
            if time.time() - cached['timestamp'] < 300:  # 5 minute cache
                logger.info(f"Using cached DNS result for {domain} from {cached['source']}.")
                return cached['result']
        
        # Try different DNS servers intelligently
        results = []
        for dns_server in self.dns_servers:
            try:
                result = await self._resolve_with_server(domain, query_type, dns_server)
                results.append(result)
                
                # Update performance stats
                self._update_dns_performance(dns_server, True)
                
            except Exception as e:
                logger.warning(f"DNS resolution for {domain} via {dns_server} failed: {e}")
                self._update_dns_performance(dns_server, False)
                continue
        
        if not results:
            raise Exception(f"DNS resolution failed for {domain}")
        
        # Select best result based on multiple factors
        best_result = self._select_best_dns_result(results, domain)
        
        # Cache the result
        self.dns_cache[cache_key] = {
            'result': best_result,
            'timestamp': time.time(),
            'source': best_result['source']
        }
        logger.info(f"Successfully resolved {domain} via {best_result['source']}.")
        
        return best_result
    
    def _select_best_dns_result(self, results: List[Dict], domain: str) -> Dict:
        """AI-driven selection of best DNS result"""
        scored_results = []
        
        for result in results:
            score = self._calculate_dns_result_score(result, domain)
            scored_results.append((result, score))
        
        return max(scored_results, key=lambda x: x[1])[0]
    
    def _calculate_dns_result_score(self, result: Dict, domain: str) -> float:
        """Calculate score for DNS result quality"""
        # Ensure result['ttl'] is an int for scoring
        try:
            ttl_val = int(result.get('ttl', 0))
        except (ValueError, TypeError):
            ttl_val = 0 # Default to 0 if conversion fails

        scores = {
            "response_time": max(0, 1 - (result.get('response_time', 1000) / 2000.0)),  # Faster = better (normalize to 2000ms max)
            "ttl_optimal": 1.0 if 300 <= ttl_val <= 3600 else 0.5,  # Optimal TTL
            "server_reliability": self.performance_stats.get(result['source'], {}).get('reliability', 0.5),
            "geographic_relevance": self._calculate_geo_relevance(result, domain)
        }
        
        return sum(scores.values()) / len(scores)
