# ~/HyperZilla/ZILLA_CORE/war_room/situational_awareness.py
import asyncio
import json
from datetime import datetime
from typing import Dict, List
import websockets

class WarRoomDashboard:
    def __init__(self):
        self.connected_clients = set()
        self.real_time_data = {
            'digital_ops': [],
            'physical_tracking': [],
            'enterprise_alerts': [],
            'fusion_correlations': []
        }
        self.alert_level = "GREEN"

    async def broadcast_situational_update(self):
        """Broadcast real-time intelligence to war room"""
        while True:
            situational_data = await self._compile_situational_data()

            # Broadcast to all connected clients
            if self.connected_clients:
                message = json.dumps({
                    'type': 'SITUATIONAL_UPDATE',
                    'timestamp': datetime.now().isoformat(),
                    'data': situational_data,
                    'alert_level': self.alert_level
                })

                await asyncio.gather(*[
                    client.send(message)
                    for client in self.connected_clients
                ])

            await asyncio.sleep(2)  # Real-time updates

    async def _compile_situational_data(self) -> Dict:
        """Compile data from all intelligence streams"""
        return {
            'active_operations': len(self.real_time_data['digital_ops']),
            'tracked_targets': len(self.real_time_data['physical_tracking']),
            'enterprise_alerts': self.real_time_data['enterprise_alerts'][-10:],  # Last 10
            'recent_correlations': self.real_time_data['fusion_correlations'][-5:],
            'system_health': await self._get_system_health()
        }

    def add_digital_operation(self, op_data: Dict):
        """Add digital operation to war room display"""
        self.real_time_data['digital_ops'].append({
            **op_data,
            'war_room_timestamp': datetime.now().isoformat()
        })

        # Trim to last 50 operations
        self.real_time_data['digital_ops'] = self.real_time_data['digital_ops'][-50:]

    def update_alert_level(self, new_level: str):
        """Update global alert level"""
        levels = ["GREEN", "YELLOW", "ORANGE", "RED"]
        if new_level in levels:
            self.alert_level = new_level
            self._trigger_alert_protocols(new_level)


# WebSocket server for real-time war room
async def war_room_server(websocket, path):
    war_room = WarRoomDashboard()
    war_room.connected_clients.add(websocket)

    try:
        async for message in websocket:
            # Handle client messages
            pass
    finally:
        war_room.connected_clients.remove(websocket)
