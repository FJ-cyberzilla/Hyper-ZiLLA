#!/usr/bin/env python3
"""
Advanced AI Code Analysis Service
Enterprise-grade code analysis with pattern matching, security detection,
and intelligent recommendations.
"""

import sys
import json
import re
import hashlib
from typing import Dict, List, Any, Tuple
from dataclasses import dataclass, asdict
from enum import Enum
import datetime


class Severity(Enum):
    """Analysis severity levels"""
    INFO = 1
    LOW = 2
    MEDIUM = 3
    HIGH = 4
    CRITICAL = 5


@dataclass
class VulnerabilityPattern:
    """Pattern definition for vulnerability detection"""
    name: str
    pattern: str
    severity: Severity
    description: str
    recommendation: str
    cwe_id: str = ""


@dataclass
class AnalysisResult:
    """Structured analysis result"""
    success: bool
    analysis: str
    confidence: float
    severity: int
    recommendations: List[str]
    error: str = ""
    detected_issues: List[Dict[str, Any]] = None
    
    def to_dict(self) -> Dict[str, Any]:
        result = asdict(self)
        if self.detected_issues is None:
            result['detected_issues'] = []
        return result


class CodeAnalyzer:
    """Advanced code analyzer with pattern matching and AI simulation"""
    
    def __init__(self, model_type: str = "advanced"):
        self.model_type = model_type
        self.vulnerability_patterns = self._load_vulnerability_patterns()
        self.quality_patterns = self._load_quality_patterns()
        
    def _load_vulnerability_patterns(self) -> Dict[str, List[VulnerabilityPattern]]:
        """Load vulnerability detection patterns by language"""
        return {
            "cpp": [
                VulnerabilityPattern(
                    name="Buffer Overflow",
                    pattern=r'\b(strcpy|strcat|gets|sprintf)\s*\(',
                    severity=Severity.CRITICAL,
                    description="Unsafe string function can cause buffer overflow",
                    recommendation="Use safe alternatives: strncpy, strncat, fgets, snprintf",
                    cwe_id="CWE-120"
                ),
                VulnerabilityPattern(
                    name="Memory Leak",
                    pattern=r'\bnew\b(?!.*\bdelete\b)',
                    severity=Severity.HIGH,
                    description="Potential memory leak detected",
                    recommendation="Ensure all dynamically allocated memory is properly freed or use smart pointers",
                    cwe_id="CWE-401"
                ),
                VulnerabilityPattern(
                    name="Null Pointer Dereference",
                    pattern=r'(?:->|\*)\w+(?!\s*(?:if|&&|\|\|))',
                    severity=Severity.HIGH,
                    description="Potential null pointer dereference",
                    recommendation="Add null pointer checks before dereferencing",
                    cwe_id="CWE-476"
                ),
                VulnerabilityPattern(
                    name="Integer Overflow",
                    pattern=r'\b(?:int|long|short)\s+\w+\s*=\s*\w+\s*[\+\-\*]\s*\w+',
                    severity=Severity.MEDIUM,
                    description="Potential integer overflow in arithmetic operation",
                    recommendation="Use checked arithmetic or larger integer types",
                    cwe_id="CWE-190"
                ),
                VulnerabilityPattern(
                    name="Use After Free",
                    pattern=r'\bdelete\b.*?(?:\n.*?)*?\b\1\b',
                    severity=Severity.CRITICAL,
                    description="Potential use-after-free vulnerability",
                    recommendation="Set pointers to nullptr after deletion or use smart pointers",
                    cwe_id="CWE-416"
                ),
            ],
            "python": [
                VulnerabilityPattern(
                    name="SQL Injection",
                    pattern=r'execute\s*\(\s*["\'].*?%s.*?["\'].*?%',
                    severity=Severity.CRITICAL,
                    description="SQL injection vulnerability detected",
                    recommendation="Use parameterized queries or prepared statements",
                    cwe_id="CWE-89"
                ),
                VulnerabilityPattern(
                    name="Command Injection",
                    pattern=r'(?:os\.system|subprocess\.call|eval|exec)\s*\(\s*.*?\+',
                    severity=Severity.CRITICAL,
                    description="Command injection vulnerability",
                    recommendation="Sanitize inputs and use subprocess with array arguments",
                    cwe_id="CWE-78"
                ),
                VulnerabilityPattern(
                    name="Hardcoded Credentials",
                    pattern=r'(?:password|api_key|secret)\s*=\s*["\'][^"\']+["\']',
                    severity=Severity.HIGH,
                    description="Hardcoded credentials detected",
                    recommendation="Use environment variables or secure credential stores",
                    cwe_id="CWE-798"
                ),
                VulnerabilityPattern(
                    name="Pickle Deserialization",
                    pattern=r'\bpickle\.loads?\s*\(',
                    severity=Severity.HIGH,
                    description="Unsafe deserialization with pickle",
                    recommendation="Use json or validate/sanitize pickled data",
                    cwe_id="CWE-502"
                ),
            ],
            "javascript": [
                VulnerabilityPattern(
                    name="XSS Vulnerability",
                    pattern=r'innerHTML\s*=|document\.write\s*\(',
                    severity=Severity.HIGH,
                    description="Potential Cross-Site Scripting (XSS) vulnerability",
                    recommendation="Use textContent or sanitize HTML input",
                    cwe_id="CWE-79"
                ),
                VulnerabilityPattern(
                    name="Eval Usage",
                    pattern=r'\beval\s*\(',
                    severity=Severity.CRITICAL,
                    description="Use of eval() detected",
                    recommendation="Avoid eval() - use JSON.parse() or Function constructor",
                    cwe_id="CWE-95"
                ),
                VulnerabilityPattern(
                    name="Weak Random",
                    pattern=r'Math\.random\s*\(',
                    severity=Severity.MEDIUM,
                    description="Math.random() is not cryptographically secure",
                    recommendation="Use crypto.getRandomValues() for security-sensitive operations",
                    cwe_id="CWE-330"
                ),
            ],
            "java": [
                VulnerabilityPattern(
                    name="SQL Injection",
                    pattern=r'Statement\s+\w+\s*=.*?executeQuery\s*\(\s*".*?\+',
                    severity=Severity.CRITICAL,
                    description="SQL injection vulnerability",
                    recommendation="Use PreparedStatement with parameterized queries",
                    cwe_id="CWE-89"
                ),
                VulnerabilityPattern(
                    name="Path Traversal",
                    pattern=r'new\s+File\s*\(\s*.*?\+',
                    severity=Severity.HIGH,
                    description="Potential path traversal vulnerability",
                    recommendation="Validate and sanitize file paths",
                    cwe_id="CWE-22"
                ),
            ],
            "go": [
                VulnerabilityPattern(
                    name="SQL Injection",
                    pattern=r'db\.Exec\s*\(\s*".*?\+',
                    severity=Severity.CRITICAL,
                    description="SQL injection vulnerability",
                    recommendation="Use parameterized queries with $1, $2, etc.",
                    cwe_id="CWE-89"
                ),
                VulnerabilityPattern(
                    name="Unhandled Error",
                    pattern=r'\w+\s*:=.*?\n(?!.*if\s+\w+\s*!=\s*nil)',
                    severity=Severity.MEDIUM,
                    description="Error not handled",
                    recommendation="Always check and handle errors",
                    cwe_id="CWE-252"
                ),
            ]
        }
    
    def _load_quality_patterns(self) -> Dict[str, List[VulnerabilityPattern]]:
        """Load code quality patterns"""
        return {
            "cpp": [
                VulnerabilityPattern(
                    name="Missing Const",
                    pattern=r'(?:std::)?string\s+\w+\s*\(',
                    severity=Severity.INFO,
                    description="Consider using const reference for string parameters",
                    recommendation="Use const std::string& for read-only parameters"
                ),
                VulnerabilityPattern(
                    name="Raw Pointer Usage",
                    pattern=r'\w+\s*\*\s*\w+\s*=\s*new\b',
                    severity=Severity.LOW,
                    description="Consider using smart pointers",
                    recommendation="Use std::unique_ptr or std::shared_ptr instead of raw pointers"
                ),
            ],
            "python": [
                VulnerabilityPattern(
                    name="Mutable Default Argument",
                    pattern=r'def\s+\w+\s*\([^)]*=\s*\[\]',
                    severity=Severity.LOW,
                    description="Mutable default argument detected",
                    recommendation="Use None as default and initialize inside function"
                ),
                VulnerabilityPattern(
                    name="Broad Exception Catching",
                    pattern=r'except\s*:',
                    severity=Severity.LOW,
                    description="Catching all exceptions is too broad",
                    recommendation="Catch specific exception types"
                ),
            ]
        }
    
    def analyze(self, code: str, language: str, file_path: str = "") -> AnalysisResult:
        """Perform comprehensive code analysis"""
        try:
            # Normalize language name
            language = language.lower()
            
            # Collect all issues
            vulnerabilities = []
            quality_issues = []
            
            # Check for vulnerabilities
            if language in self.vulnerability_patterns:
                for pattern in self.vulnerability_patterns[language]:
                    matches = list(re.finditer(pattern.pattern, code, re.MULTILINE))
                    for match in matches:
                        line_num = code[:match.start()].count('\n') + 1
                        vulnerabilities.append({
                            'type': 'vulnerability',
                            'name': pattern.name,
                            'severity': pattern.severity.name,
                            'severity_level': pattern.severity.value,
                            'description': pattern.description,
                            'recommendation': pattern.recommendation,
                            'cwe_id': pattern.cwe_id,
                            'line': line_num,
                            'snippet': match.group(0)[:50]
                        })
            
            # Check for quality issues
            if language in self.quality_patterns:
                for pattern in self.quality_patterns[language]:
                    matches = list(re.finditer(pattern.pattern, code, re.MULTILINE))
                    for match in matches:
                        line_num = code[:match.start()].count('\n') + 1
                        quality_issues.append({
                            'type': 'quality',
                            'name': pattern.name,
                            'severity': pattern.severity.name,
                            'severity_level': pattern.severity.value,
                            'description': pattern.description,
                            'recommendation': pattern.recommendation,
                            'line': line_num
                        })
            
            # Calculate metrics
            lines_of_code = len(code.split('\n'))
            complexity_score = self._calculate_complexity(code, language)
            
            # Determine overall severity
            max_severity = 0
            if vulnerabilities:
                max_severity = max(v['severity_level'] for v in vulnerabilities)
            elif quality_issues:
                max_severity = max(q['severity_level'] for q in quality_issues)
            
            # Calculate confidence based on pattern matches and code size
            confidence = min(0.95, 0.7 + (len(vulnerabilities) * 0.05))
            
            # Generate analysis text
            analysis_parts = []
            
            if vulnerabilities:
                analysis_parts.append(f"Found {len(vulnerabilities)} security issue(s):")
                for vuln in vulnerabilities[:5]:  # Limit to top 5
                    analysis_parts.append(
                        f"  • {vuln['name']} (Line {vuln['line']}): {vuln['description']}"
                    )
            
            if quality_issues:
                analysis_parts.append(f"\nFound {len(quality_issues)} code quality issue(s):")
                for issue in quality_issues[:3]:  # Limit to top 3
                    analysis_parts.append(
                        f"  • {issue['name']} (Line {issue['line']}): {issue['description']}"
                    )
            
            analysis_parts.append(f"\nCode Metrics:")
            analysis_parts.append(f"  • Lines of Code: {lines_of_code}")
            analysis_parts.append(f"  • Estimated Complexity: {complexity_score:.1f}")
            
            if not vulnerabilities and not quality_issues:
                analysis_parts.append("\n✓ No major issues detected. Code appears clean.")
            
            # Collect recommendations
            recommendations = []
            for vuln in vulnerabilities:
                recommendations.append(vuln['recommendation'])
            for issue in quality_issues[:3]:
                recommendations.append(issue['recommendation'])
            
            # Add general recommendations
            if language == "cpp":
                recommendations.extend([
                    "Consider using static analysis tools like cppcheck or clang-tidy",
                    "Enable compiler warnings (-Wall -Wextra -Werror)",
                    "Use address sanitizer during development"
                ])
            elif language == "python":
                recommendations.extend([
                    "Use pylint or flake8 for additional checks",
                    "Consider type hints for better code clarity",
                    "Use virtual environments for dependency management"
                ])
            
            return AnalysisResult(
                success=True,
                analysis="\n".join(analysis_parts),
                confidence=confidence,
                severity=max_severity,
                recommendations=list(dict.fromkeys(recommendations)),  # Remove duplicates
                detected_issues=vulnerabilities + quality_issues
            )
            
        except Exception as e:
            return AnalysisResult(
                success=False,
                analysis="",
                confidence=0.0,
                severity=0,
                recommendations=[],
                error=f"Analysis failed: {str(e)}"
            )
    
    def _calculate_complexity(self, code: str, language: str) -> float:
        """Calculate estimated cyclomatic complexity"""
        complexity = 1.0
        
        # Control flow keywords that increase complexity
        control_keywords = {
            'cpp': ['if', 'else', 'for', 'while', 'switch', 'case', 'catch', '&&', '||'],
            'python': ['if', 'elif', 'else', 'for', 'while', 'except', 'and', 'or'],
            'javascript': ['if', 'else', 'for', 'while', 'switch', 'case', 'catch', '&&', '||'],
            'java': ['if', 'else', 'for', 'while', 'switch', 'case', 'catch', '&&', '||'],
            'go': ['if', 'else', 'for', 'switch', 'case', '&&', '||']
        }
        
        if language in control_keywords:
            for keyword in control_keywords[language]:
                complexity += code.count(keyword)
        
        # Normalize by lines of code
        lines = len(code.split('\n'))
        if lines > 0:
            complexity = complexity / lines * 10
        
        return min(complexity, 10.0)  # Cap at 10


def main():
    """Main entry point for AI service"""
    if len(sys.argv) != 3:
        print(json.dumps({
            'success': False,
            'error': 'Usage: ai_service.py <input_json> <output_json>'
        }))
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_file = sys.argv[2]
    
    try:
        # Read input
        with open(input_file, 'r') as f:
            request = json.load(f)
        
        code = request.get('code', '')
        language = request.get('language', '')
        file_path = request.get('file_path', '')
        model_type = request.get('model_type', 'advanced')
        
        # Validate input
        if not code:
            raise ValueError("Code is required")
        if not language:
            raise ValueError("Language is required")
        
        # Perform analysis
        analyzer = CodeAnalyzer(model_type=model_type)
        result = analyzer.analyze(code, language, file_path)
        
        # Write output
        with open(output_file, 'w') as f:
            json.dump(result.to_dict(), f, indent=2)
        
        sys.exit(0 if result.success else 1)
        
    except Exception as e:
        error_result = {
            'success': False,
            'analysis': '',
            'confidence': 0.0,
            'severity': 0,
            'recommendations': [],
            'error': str(e),
            'detected_issues': []
        }
        
        try:
            with open(output_file, 'w') as f:
                json.dump(error_result, f, indent=2)
        except:
            pass
        
        print(json.dumps(error_result))
        sys.exit(1)


if __name__ == '__main__':
    main()
