# core/monitoring/alert_manager.py
class AlertManager:
    """Intelligent alert management with escalation and correlation"""
    
    def __init__(self, config: Dict):
        self.config = config
        self.alert_rules = self._load_alert_rules()
        self.escalation_policies = self._load_escalation_policies()
        self.alert_correlator = AlertCorrelator(config)
        
    async def trigger_alert(self, alert_type: str, data: Dict, severity: str = "medium"):
        """Trigger intelligent alert with correlation and escalation"""
        alert = {
            'alert_id': self._generate_alert_id(),
            'type': alert_type,
            'severity': severity,
            'timestamp': datetime.utcnow(),
            'source': data.get('source', 'monitoring_system'),
            'data': data,
            'status': 'active',
            'acknowledged': False,
            'escalation_level': 0
        }
        
        # Correlate with existing alerts
        correlated_alerts = await self.alert_correlator.correlate_alert(alert)
        if correlated_alerts:
            alert['correlated_with'] = correlated_alerts
        
        # Store alert
        await self._store_alert(alert)
        
        # Execute alert actions based on severity
        await self._execute_alert_actions(alert)
        
        # Start escalation timer if not acknowledged
        asyncio.create_task(self._start_escalation_timer(alert))
        
        logger.warning(f"🚨 Alert triggered: {alert_type} - Severity: {severity}")
        return alert
    
    async def _execute_alert_actions(self, alert: Dict):
        """Execute appropriate actions based on alert severity"""
        severity_actions = {
            'low': [self._log_alert, self._update_dashboard],
            'medium': [self._log_alert, self._update_dashboard, self._notify_team],
            'high': [self._log_alert, self._update_dashboard, self._notify_team, self._page_on_call],
            'critical': [self._log_alert, self._update_dashboard, self._notify_team, 
                        self._page_on_call, self._escalate_management]
        }
        
        actions = severity_actions.get(alert['severity'], [])
        for action in actions:
            try:
                await action(alert)
            except Exception as e:
                logger.error(f"Alert action failed: {e}")
    
    async def _start_escalation_timer(self, alert: Dict):
        """Start escalation timer for unacknowledged alerts"""
        escalation_delays = {
            'low': timedelta(hours=4),
            'medium': timedelta(hours=2),
            'high': timedelta(minutes=30),
            'critical': timedelta(minutes=5)
        }
        
        delay = escalation_delays.get(alert['severity'], timedelta(hours=1))
        await asyncio.sleep(delay.total_seconds())
        
        # Check if alert still active and unacknowledged
        current_alert = await self._get_alert(alert['alert_id'])
        if current_alert and not current_alert['acknowledged']:
            await self._escalate_alert(current_alert)
