# utils/debug_engine.py
class DebugEngine:
    def __init__(self):
        self.log_level = "INFO"
        self.performance_monitor = PerformanceMonitor()
        self.error_tracker = ErrorTracker()
        
    def comprehensive_debug_menu(self):
        """Advanced debugging interface"""
        debug_options = {
            "1": "System Health Check",
            "2": "Agent Performance Analysis", 
            "3": "Database Integrity Scan",
            "4": "Network Connection Test",
            "5": "ML Model Validation",
            "6": "Security Audit Log",
            "7": "Real-time Monitoring",
            "8": "Generate Diagnostic Report"
        }
        
        while True:
            self._display_debug_menu(debug_options)
            choice = input("Select debug option: ")
            
            if choice == "1":
                self._run_health_check()
            elif choice == "2":
                self._analyze_agent_performance()
            # ... other options
    
    def _run_health_check(self):
        """Comprehensive system healt            "ML Models Loaded": self._check_models(),
h check"""
        checks = {
            "Database Connection": self._check_database(),
            "AI Agents Status": self._check_agents(),
            "Security Systems": self._check_security(),
            "Network Connectivity": self._check_network()
        }
        
        for check_name, check_func in checks.items():
            status = check_func()
            color = Color.GREEN if status else Color.RED
            print(f"{color}✓ {check_name}: {'PASS' if status else 'FAIL'}{Color.END}")
