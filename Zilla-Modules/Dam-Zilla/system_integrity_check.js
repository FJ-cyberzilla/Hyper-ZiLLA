import { readdirSync, existsSync } from 'fs';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';

class SystemIntegrityCheck {
    constructor() {
        this.criticalFiles = [
            'main.js',
            'security_lock.js', 
            'startup.sh',
            'package.json',
            'core/quantum_ml/conscious_decision_engine.jl',
            'core/recon_master/recon_master.js',
            'autonomous_agents/self_aware_system/conscious_decision_framework.js',
            'vintage_war_room/retro_cli_interface/vintage_green_terminal.js',
            'config/main_config.yaml'
        ];
    }

    async verifyCompleteStructure() {
        console.log('🔍 VERIFYING ZILLA-DAM COMPLETE STRUCTURE...\n');

        const results = {
            passed: 0,
            failed: 0,
            missing: [],
            details: {}
        };

        for (const file of this.criticalFiles) {
            if (existsSync(file)) {
                console.log(`✅ ${file}`);
                results.passed++;
                results.details[file] = 'PRESENT';
            } else {
                console.log(`❌ ${file} - MISSING`);
                results.failed++;
                results.missing.push(file);
            }
        }

        // Check directory structure
        const requiredDirs = [
            'core', 'advanced_sensors', 'autonomous_agents', 
            'vintage_war_room', 'modular_system', 'deployment_suite'
        ];

        console.log('\n📁 VERIFYING DIRECTORY STRUCTURE...');
        for (const dir of requiredDirs) {
            if (existsSync(dir)) {
                console.log(`✅ ${dir}/`);
            } else {
                console.log(`❌ ${dir}/ - MISSING`);
                results.missing.push(dir);
            }
        }

        return results;
    }

    async checkDependencies() {
        console.log('\n📦 VERIFYING DEPENDENCIES...');
        
        try {
            const packageJson = JSON.parse(await readFileSync('package.json', 'utf8'));
            const requiredDeps = [
                'axios', 'chalk', 'ora', 'figlet', 
                'crypto-js', 'machine-id', 'node-os-utils'
            ];

            const missingDeps = requiredDeps.filter(dep => !packageJson.dependencies[dep]);
            
            if (missingDeps.length === 0) {
                console.log('✅ All critical dependencies configured');
            } else {
                console.log('❌ Missing dependencies:', missingDeps);
            }

            return missingDeps;
        } catch (error) {
            console.log('❌ Cannot read package.json');
            return ['package.json'];
        }
    }

    async runCompleteCheck() {
        console.log('🐉 ZILLA-DAM COMPLETE SYSTEM INTEGRITY CHECK\n');
        console.log('═'.repeat(60));

        const structureResults = await this.verifyCompleteStructure();
        const dependencyResults = await this.checkDependencies();

        console.log('\n📊 SUMMARY:');
        console.log(`✅ Files Present: ${structureResults.passed}`);
        console.log(`❌ Files Missing: ${structureResults.failed}`);
        console.log(`📦 Dependencies Missing: ${dependencyResults.length}`);

        if (structureResults.failed === 0 && dependencyResults.length === 0) {
            console.log('\n🎉 ZILLA-DAM SYSTEM INTEGRITY: PERFECT');
            console.log('🚀 READY FOR DEPLOYMENT!');
        } else {
            console.log('\n⚠️  ZILLA-DAM SYSTEM INTEGRITY: ISSUES DETECTED');
            console.log('🔧 Please fix the missing components above');
        }

        return {
            structure: structureResults,
            dependencies: dependencyResults,
            overall: structureResults.failed === 0 && dependencyResults.length === 0
        };
    }
}

// Run integrity check
const checker = new SystemIntegrityCheck();
checker.runCompleteCheck();
