# cli/advanced_menu.py
class AdvancedIntelligenceMenu:
    def __init__(self):
        self.bot_hunter = MaliciousBotHunter()
        self.payload_analyzer = AdvancedPayloadAnalyzer()
        self.shared_detector = SharedAccountsDetector()
        
    def display_advanced_intel_menu(self):
        """Advanced intelligence operations menu"""
        menu_text = f"""
{Color.BOLD}🔍 ADVANCED INTELLIGENCE OPERATIONS{Color.END}
{Color.CYAN}═{'═' * 58}{Color.END}

{Color.GREEN}1. {Color.BOLD}🕵️  User Deep Forensic Analysis{Color.END}
{Color.GREEN}2. {Color.BOLD}🤖 Malicious Bot Detection & Neutralization{Color.END}  
{Color.GREEN}3. {Color.BOLD}📡 Payload Analysis & Source Attribution{Color.END}
{Color.GREEN}4. {Color.BOLD}👥 Shared Accounts Detection{Color.END}
{Color.GREEN}5. {Color.BOLD}🎭 Deception & Manipulation Detection{Color.END}
{Color.GREEN}6. {Color.BOLD}🔋 Hardware & Behavioral Fingerprinting{Color.END}
{Color.GREEN}7. {Color.BOLD}🛡️  Active Countermeasures{Color.END}
{Color.GREEN}8. {Color.BOLD}📊 Intelligence Dashboard{Color.END}

{Color.RED}0. {Color.BOLD}← Back to Main Menu{Color.END}

{Color.CYAN}═{'═' * 58}{Color.END}
        """
        print(menu_text)
        
        choice = input(f"{Color.YELLOW}👉 Select operation: {Color.END}")
        return choice
    
    async def execute_bot_hunt(self, target_user: str):
        """Execute comprehensive bot hunting operation"""
        print(f"{Color.RED}🚀 INITIATING BOT HUNT: {target_user}{Color.END}")
        
        # Animated loading
        await self._show_loading_animation("Scanning for malicious bots")
        
        results = await self.bot_hunter.detect_malicious_bots(target_user)
        
        # Display results with advanced visualization
        self._display_bot_hunt_results(results)
        
        if results["recommended_actions"]:
            response = input(f"{Color.RED}🚨 Execute countermeasures? (y/N): {Color.END}")
            if response.lower() == 'y':
                await self._execute_countermeasures(results)
