import sqlite3 from 'sqlite3';
import { Pool } from 'pg';
import { Sequelize } from 'sequelize';

export class DatabaseManager {
    constructor() {
        this.sqlite = null;
        this.postgres = null;
        this.sequelize = null;
        this.init();
    }

    async init() {
        await this.initializeSQLite();
        await this.initializePostgreSQL();
        await this.initializeORM();
    }

    async initializeSQLite() {
        this.sqlite = new sqlite3.Database(':memory:', (err) => {
            if (err) {
                console.error('SQLite initialization error:', err);
                return;
            }
            console.log('✅ SQLite (in-memory) initialized');
        });

        // Create optimized tables
        await this.createSQLiteSchema();
    }

    async initializePostgreSQL() {
        if (process.env.POSTGRES_URL) {
            this.postgres = new Pool({
                connectionString: process.env.POSTGRES_URL,
                max: 20,
                idleTimeoutMillis: 30000,
                connectionTimeoutMillis: 2000,
            });

            try {
                await this.postgres.query('SELECT 1');
                console.log('✅ PostgreSQL connected');
                await this.createPostgresSchema();
            } catch (error) {
                console.log('❌ PostgreSQL unavailable, using SQLite only');
                this.postgres = null;
            }
        }
    }

    async initializeORM() {
        this.sequelize = new Sequelize({
            dialect: 'sqlite',
            storage: ':memory:',
            logging: false,
            pool: {
                max: 5,
                min: 0,
                acquire: 30000,
                idle: 10000
            }
        });

        await this.defineModels();
    }

    async createSQLiteSchema() {
        const schemas = [
            `CREATE TABLE IF NOT EXISTS operations (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                phone_number TEXT UNIQUE,
                raw_data BLOB,
                analyzed_data BLOB,
                risk_score REAL,
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
            )`,
            
            `CREATE TABLE IF NOT EXISTS social_profiles (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                operation_id INTEGER,
                platform TEXT,
                username TEXT,
                user_id TEXT,
                profile_data BLOB,
                confidence REAL,
                FOREIGN KEY (operation_id) REFERENCES operations (id)
            )`,
            
            `CREATE TABLE IF NOT EXISTS vpn_detections (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                operation_id INTEGER,
                ip_address TEXT,
                location TEXT,
                vpn_provider TEXT,
                confidence REAL,
                detected_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                FOREIGN KEY (operation_id) REFERENCES operations (id)
            )`,
            
            `CREATE INDEX IF NOT EXISTS idx_phone_number ON operations(phone_number)`,
            `CREATE INDEX IF NOT EXISTS idx_risk_score ON operations(risk_score)`,
            `CREATE INDEX IF NOT EXISTS idx_platform ON social_profiles(platform)`
        ];

        for (const schema of schemas) {
            await this.runSQLite(schema);
        }
    }
}
