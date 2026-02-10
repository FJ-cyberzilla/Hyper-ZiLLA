import { createLogger, transports, format } from 'winston';
import { promises as fs } from 'fs';

export class ConsolidatedErrorHandler {
    constructor() {
        this.logger = this.setupLogger();
        this.errorCounts = new Map();
        this.setupGlobalHandlers();
    }

    setupLogger() {
        return createLogger({
            level: 'info',
            format: format.combine(
                format.timestamp(),
                format.errors({ stack: true }),
                format.json()
            ),
            defaultMeta: { service: 'zilla-dam' },
            transports: [
                new transports.File({ 
                    filename: 'logs/error.log', 
                    level: 'error',
                    maxsize: 5242880, // 5MB
                    maxFiles: 5
                }),
                new transports.File({ 
                    filename: 'logs/combined.log',
                    maxsize: 5242880,
                    maxFiles: 5
                }),
                new transports.Console({
                    format: format.combine(
                        format.colorize(),
                        format.simple()
                    )
                })
            ]
        });
    }

    setupGlobalHandlers() {
        process.on('uncaughtException', (error) => {
            this.handleFatalError('uncaughtException', error);
        });

        process.on('unhandledRejection', (reason, promise) => {
            this.handleFatalError('unhandledRejection', reason);
        });

        process.on('SIGTERM', () => {
            this.handleGracefulShutdown('SIGTERM');
        });

        process.on('SIGINT', () => {
            this.handleGracefulShutdown('SIGINT');
        });
    }

    async handleError(context, error, severity = 'error') {
        const errorId = this.generateErrorId(context, error);
        const count = this.errorCounts.get(errorId) || 0;
        this.errorCounts.set(errorId, count + 1);

        const errorInfo = {
            errorId,
            context,
            message: error.message,
            stack: error.stack,
            severity,
            count: count + 1,
            timestamp: new Date().toISOString()
        };

        // Log based on severity
        if (severity === 'error') {
            this.logger.error(errorInfo);
        } else if (severity === 'warn') {
            this.logger.warn(errorInfo);
        } else {
            this.logger.info(errorInfo);
        }

        // Alert on repeated errors
        if (count + 1 >= 5) {
            await this.alertRepeatedError(errorInfo);
        }

        return errorId;
    }

    handleFatalError(type, error) {
        this.logger.error({
            type,
            message: 'Fatal error occurred',
            error: error.message,
            stack: error.stack
        });

        // Attempt graceful shutdown
        setTimeout(() => {
            process.exit(1);
        }, 5000).unref();
    }

    async handleGracefulShutdown(signal) {
        this.logger.info(`Received ${signal}, starting graceful shutdown`);
        
        try {
            // Close database connections
            await this.closeDatabaseConnections();
            
            // Save current state
            await this.saveApplicationState();
            
            // Close servers
            await this.closeServers();
            
            this.logger.info('Graceful shutdown completed');
            process.exit(0);
        } catch (error) {
            this.logger.error('Error during graceful shutdown', error);
            process.exit(1);
        }
    }

    async alertRepeatedError(errorInfo) {
        // Implement alerting logic (email, Slack, etc.)
        console.log(`🚨 REPEATED ERROR ALERT: ${errorInfo.errorId}`);
        console.log(`Context: ${errorInfo.context}`);
        console.log(`Count: ${errorInfo.count}`);
    }
}
