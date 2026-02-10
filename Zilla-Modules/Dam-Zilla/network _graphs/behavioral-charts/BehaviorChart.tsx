import React from 'react';
import { TrendAnalysis } from './TrendAnalysis';
import { PatternVisualizer } from './PatternVisualizer';
import { ZillaBehavioralData } from './types';

interface ZillaBehaviorChartProps {
    behavioralData: ZillaBehavioralData;
    timezone: string;
}

export const ZillaBehaviorChart: React.FC<ZillaBehaviorChartProps> = ({ 
    behavioralData, 
    timezone 
}) => {
    const renderActivityHeatmap = () => {
        const heatmap = behavioralData.activityPatterns;
        const timezones = Object.keys(heatmap);
        
        return timezones.map(tz => (
            <div key={tz} className="timezone-activity">
                <h4>{tz.toUpperCase()} Activity Pattern</h4>
                <div className="activity-bars">
                    {Array.from({ length: 24 }).map((_, hour) => (
                        <div
                            key={hour}
                            className="activity-bar"
                            style={{
                                height: `${(heatmap[tz][hour] || 0) * 100}%`,
                                backgroundColor: heatmap[tz][hour] > 0.7 ? '#ff4444' : 
                                               heatmap[tz][hour] > 0.4 ? '#ffaa00' : '#44ff44'
                            }}
                            title={`${hour}:00 - Activity: ${Math.round((heatmap[tz][hour] || 0) * 100)}%`}
                        >
                            {hour}
                        </div>
                    ))}
                </div>
            </div>
        ));
    };

    const renderAnomalyDetection = () => {
        return behavioralData.anomalyDetection.map((anomaly, index) => (
            <div key={index} className="anomaly-alert">
                <span className="alert-icon">🚨</span>
                <span className="alert-message">{anomaly.description}</span>
                <span className="alert-confidence">
                    Confidence: {Math.round(anomaly.confidence * 100)}%
                </span>
            </div>
        ));
    };

    return (
        <div className="zilla-behavior-chart">
            <div className="primary-timezone">
                <h3>Primary Timezone Analysis</h3>
                <p>Detected: <strong>{timezone}</strong></p>
            </div>

            <div className="activity-heatmaps">
                <h3>Cross-Timezone Activity Patterns</h3>
                {renderActivityHeatmap()}
            </div>

            <div className="anomaly-detection">
                <h3>Behavioral Anomalies</h3>
                {renderAnomalyDetection()}
            </div>

            <TrendAnalysis data={behavioralData} />
            <PatternVisualizer patterns={behavioralData.activityPatterns} />
        </div>
    );
};
