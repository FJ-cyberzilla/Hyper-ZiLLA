# cli/uninstall_menu.py
class UninstallMenu:
    def __init__(self, command_center):
        self.command_center = command_center
        self.uninstall_manager = UninstallManager(command_center.base_path)
        self.install_tracker = InstallationTracker(command_center.base_path)
        
    def display_uninstall_menu(self):
        """Display uninstall confirmation menu"""
        print(f"""
{Color.RED}╔══════════════════════════════════════════════════════════════╗
║                    🚨 UNINSTALLATION MENU 🚨                 ║
╠══════════════════════════════════════════════════════════════╣
║                                                              ║
║   This will completely remove Telegram OSINT Platform       ║
║   from your system. All data and configurations will        ║
║   be deleted.                                               ║
║                                                              ║
║   {Color.YELLOW}Options:{Color.RED}                                                  ║
║                                                              ║
║   {Color.GREEN}1.{Color.RED} 📦 Create backup and uninstall                        ║
║   {Color.GREEN}2.{Color.RED} 🗑️  Uninstall without backup                          ║
║   {Color.GREEN}3.{Color.RED} 📊 Show usage statistics                              ║
║   {Color.GREEN}4.{Color.RED} ❌ Cancel and return to main menu                     ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╗{Color.END}
        """)
        
        choice = input(f"{Color.YELLOW}👉 Select option: {Color.END}")
        return choice
    
    async def execute_uninstall(self, backup: bool = True):
        """Execute the uninstallation process"""
        # Generate final report
        report = self.uninstall_manager.generate_uninstall_report()
        
        # Display final statistics
        self._display_final_stats(report)
        
        # Confirm one more time
        confirm = input(f"{Color.RED}🚨 Type 'CONFIRM UNINSTALL' to proceed: {Color.END}")
        if confirm != "CONFIRM UNINSTALL":
            print(f"{Color.YELLOW}Uninstallation cancelled.{Color.END}")
            return False
        
        # Execute uninstall
        print(f"{Color.RED}🚀 Starting uninstallation...{Color.END}")
        success = self.uninstall_manager.execute_uninstall(backup=backup)
        
        if success:
            print(f"{Color.GREEN}✅ Uninstallation completed successfully!{Color.END}")
            return True
        else:
            print(f"{Color.RED}❌ Uninstallation failed!{Color.END}")
            return False
