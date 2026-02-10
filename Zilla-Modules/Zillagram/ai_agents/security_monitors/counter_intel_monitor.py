# ai_agents/security_monitors/counter_intel_monitor.py
class RealTimeCounterIntelMonitor:
    def __init__(self, telegram_client):
        self.client = telegram_client
        self.spy_detector = GeolocationSpyDetector()
        self.trap_engine = CounterIntelligenceTrap(telegram_client)
        self.alert_system = MissionAccomplishedAlert()
        self.active_operations = {}
        
    async def monitor_for_geolocation_threats(self):
        """Continuous monitoring for geolocation threats"""
        while True:
            try:
                # Get recent messages from monitored chats
                recent_messages = await self._get_recent_messages()
                
                for message in recent_messages:
                    threat_analysis = await self.spy_detector.detect_geolocation_harvesting(
                        message, await self._get_user_data(message.sender_id)
                    )
                    
                    if threat_analysis["threat_level"] > 0.7:
                        await self._activate_counter_intel_protocol(message, threat_analysis)
                
                await asyncio.sleep(30)  # Check every 30 seconds
                
            except Exception as e:
                logger.error(f"Counter-intel monitoring error: {e}")
                await asyncio.sleep(60)
    
    async def _activate_counter_intel_protocol(self, message: Dict, threat_analysis: Dict):
        """Activate counter-intelligence measures against detected spy"""
        spy_user_id = message.sender_id
        
        # Don't activate multiple traps for same user
        if spy_user_id in self.active_operations:
            return
        
        logger.warning(f"🚨 Activating counter-intel against suspected spy: {spy_user_id}")
        
        # Setup location deception trap
        trap_config = await self.trap_engine.setup_geolocation_trap(spy_user_id, message)
        
        # Monitor for trap activation
        self.active_operations[spy_user_id] = {
            "trap_config": trap_config,
            "activation_time": datetime.now(),
            "monitoring_task": asyncio.create_task(
                self._monitor_trap_activation(spy_user_id, trap_config["trap_id"])
            )
        }
    
    async def _monitor_trap_activation(self, spy_user_id: int, trap_id: str):
        """Monitor for when spy takes the bait"""
        trap_activation_indicators = [
            "location_based_followup",
            "timezone_reference",
            "cultural_reference",
            "meetup_attempt",
            "increased_questioning"
        ]
        
        start_time = datetime.now()
        max_monitoring_time = timedelta(hours=24)
        
        while datetime.now() - start_time < max_monitoring_time:
            try:
                # Check for trap activation indicators
                activation_detected = await self._check_trap_activation(spy_user_id, trap_activation_indicators)
                
                if activation_detected:
                    # Mission accomplished!
                    deception_details = await self._gather_deception_intel(spy_user_id, trap_id)
                    await self.alert_system.trigger_mission_accomplished(trap_id, spy_user_id, deception_details)
                    
                    # Clean up
                    del self.active_operations[spy_user_id]
                    return
                
                await asyncio.sleep(300)  # Check every 5 minutes
                
            except Exception as e:
                logger.error(f"Trap monitoring error: {e}")
                await asyncio.sleep(60)
        
        # Trap expired
        logger.info(f"Trap {trap_id} expired without activation")
        del self.active_operations[spy_user_id]
