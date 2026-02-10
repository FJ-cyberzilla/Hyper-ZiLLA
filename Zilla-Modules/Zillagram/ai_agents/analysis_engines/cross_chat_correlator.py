# ai_agents/analysis_engines/cross_chat_correlator.py
import sqlite3
from datetime import datetime, timedelta
from collections import defaultdict
import networkx as nx

class CrossChatCorrelationEngine:
    def __init__(self, db_path: str):
        self.db_path = db_path
        self.message_graph = nx.MultiGraph()
        
    async def track_message_propagation(self, target_message: Dict) -> Dict:
        """Track how a message propagates across different chats"""
        target_fingerprint = self._create_message_fingerprint(target_message)
        
        # Search for similar messages across all monitored chats
        similar_messages = await self._find_similar_messages(target_fingerprint)
        
        propagation_analysis = {
            "original_source": await self._identify_original_source(similar_messages),
            "propagation_path": self._reconstruct_propagation_path(similar_messages),
            "forwarding_network": self._build_forwarding_network(similar_messages),
            "leakage_points": self._identify_leakage_points(similar_messages),
            "trust_metrics": self._calculate_trust_metrics(similar_messages)
        }
        
        return propagation_analysis
    
    async def _find_similar_messages(self, target_fingerprint: MessageFingerprint) -> List[Dict]:
        """Find messages with similar fingerprints across all chats"""
        similar_messages = []
        
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        
        # Search by content hash (exact matches)
        cursor.execute('''
            SELECT * FROM messages 
            WHERE content_hash = ? OR semantic_hash = ?
        ''', (target_fingerprint.content_hash, target_fingerprint.semantic_hash))
        
        exact_matches = cursor.fetchall()
        similar_messages.extend(self._format_message_results(exact_matches))
        
        # Search for semantic similarities
        cursor.execute('''
            SELECT * FROM messages 
            WHERE date > datetime('now', '-7 days')
            AND chat_id != ?
        ''', (target_message['chat_id'],))
        
        all_recent_messages = cursor.fetchall()
        
        for msg in all_recent_messages:
            if self._calculate_semantic_similarity(target_fingerprint, msg) > 0.8:
                similar_messages.append(self._format_message_result(msg))
        
        conn.close()
        return similar_messages
    
    def _reconstruct_propagation_path(self, similar_messages: List[Dict]) -> List[Dict]:
        """Reconstruct the path of message propagation"""
        if not similar_messages:
            return []
        
        # Sort by timestamp to find chronological order
        sorted_messages = sorted(similar_messages, key=lambda x: x['date'])
        
        propagation_path = []
        current_chain = [sorted_messages[0]]
        
        for i in range(1, len(sorted_messages)):
            current_msg = sorted_messages[i]
            prev_msg = current_chain[-1]
            
            # Check if this continues the chain
            if self._is_likely_forward(prev_msg, current_msg):
                current_chain.append(current_msg)
            else:
                # Start new chain
                if len(current_chain) > 1:
                    propagation_path.append(current_chain)
                current_chain = [current_msg]
        
        if len(current_chain) > 1:
            propagation_path.append(current_chain)
        
        return propagation_path
    
    def _build_forwarding_network(self, similar_messages: List[Dict]) -> Dict:
        """Build network graph of message forwarding"""
        G = nx.DiGraph()
        
        for message in similar_messages:
            user_id = message['user_id']
            chat_id = message['chat_id']
            timestamp = message['date']
            
            # Add node for user
            G.add_node(user_id, type='user', first_seen=timestamp)
            
            # Add node for chat
            G.add_node(chat_id, type='chat')
            
            # Add edge from user to chat
            G.add_edge(user_id, chat_id, message_id=message['id'], timestamp=timestamp)
        
        # Analyze network properties
        return {
            "graph_data": self._serialize_graph(G),
            "centrality_metrics": self._calculate_centrality(G),
            "community_structure": self._detect_communities(G),
            "key_influencers": self._identify_influencers(G),
            "propagation_speed": self._calculate_propagation_speed(G, similar_messages)
        }
    
    def _identify_leakage_points(self, similar_messages: List[Dict]) -> List[Dict]:
        """Identify where information leaked between chats"""
        leakage_points = []
        
        # Group by time windows to find coordinated sharing
        time_windows = self._create_time_windows(similar_messages)
        
        for window in time_windows:
            if len(window['messages']) > 1:
                # Multiple similar messages in short time frame = potential leakage
                leakage_point = {
                    "timestamp": window['window_start'],
                    "messages": window['messages'],
                    "leakage_confidence": self._calculate_leakage_confidence(window),
                    "suspected_leakers": self._identify_suspected_leakers(window['messages'])
                }
                leakage_points.append(leakage_point)
        
        return leakage_points
