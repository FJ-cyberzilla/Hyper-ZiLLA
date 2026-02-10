# cli/leak_detection_menu.py
class LeakDetectionMenu:
    def __init__(self, command_center):
        self.command_center = command_center
        self.forwarding_monitor = RealTimeForwardingMonitor(
            command_center.correlation_engine
        )
        
    def display_leak_detection_menu(self):
        """Display leak detection and source attribution menu"""
        print(f"""
{Color.RED}╔══════════════════════════════════════════════════════════════╗
║               🕵️  LEAK DETECTION & SOURCE TRACING            ║
╠══════════════════════════════════════════════════════════════╣
║                                                              ║
║   Detect if your conversations are being forwarded          ║
║   and trace the source of suspicious messages               ║
║                                                              ║
║   {Color.GREEN}1.{Color.RED} 🔍 Analyze Specific Message for Leaks                ║
║   {Color.GREEN}2.{Color.RED} 📡 Monitor Conversation for Forwarding              ║
║   {Color.GREEN}3.{Color.RED} 🎯 Trace Message Source & Propagation              ║
║   {Color.GREEN}4.{Color.RED} 🌐 Cross-Chat Correlation Analysis                 ║
║   {Color.GREEN}5.{Color.RED} 🚨 Active Leak Alerts & Countermeasures           ║
║   {Color.GREEN}6.{Color.RED} 📊 Leak Detection Statistics                      ║
║                                                              ║
║   {Color.YELLOW}0.{Color.RED} ← Back to Main Menu                              ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╗{Color.END}
        """)
        
        choice = input(f"{Color.YELLOW}👉 Select operation: {Color.END}")
        return choice
    
    async def analyze_message_leakage(self, message_id: int):
        """Analyze specific message for leakage"""
        print(f"{Color.CYAN}🔍 Analyzing message {message_id} for leaks...{Color.END}")
        
        # Get message from database
        message = await self._get_message_by_id(message_id)
        if not message:
            print(f"{Color.RED}❌ Message not found{Color.END}")
            return
        
        # Run leakage analysis
        leak_analysis = await self.forwarding_monitor.analyze_message_leakage(message)
        
        # Display results
        self._display_leak_analysis(leak_analysis)
        
        if leak_analysis['leak_confidence'] > 0.7:
            response = input(f"{Color.RED}🚨 High leak probability! View details? (y/N): {Color.END}")
            if response.lower() == 'y':
                await self._show_detailed_propagation(leak_analysis)
    
    def _display_leak_analysis(self, analysis: Dict):
        """Display leak analysis results"""
        confidence = analysis['leak_confidence']
        color = Color.RED if confidence > 0.7 else Color.YELLOW if confidence > 0.4 else Color.GREEN
        
        print(f"""
{Color.CYAN}╔══════════════════════════════════════════════════════════════╗
║                     LEAK ANALYSIS RESULTS                    ║
╠══════════════════════════════════════════════════════════════╣{Color.END}

{color}🔒 Leak Confidence: {confidence:.1%}{Color.END}

{Color.YELLOW}📊 Indicators:{Color.END}
  • Cross-chat appearance: {analysis['leak_indicators']['cross_chat_appearance']}
  • Temporal proximity: {analysis['leak_indicators']['temporal_anomalies']}
  • Trust violations: {analysis['leak_indicators']['trust_network_violation']}

{Color.BLUE}🎯 Suspected Leaker: {analysis.get('suspected_leaker', 'Unknown')}{Color.END}

{Color.GREEN}💡 Recommended Actions:{Color.END}
  {chr(10).join(['  • ' + action for action in analysis['recommended_actions']])}
        """)
