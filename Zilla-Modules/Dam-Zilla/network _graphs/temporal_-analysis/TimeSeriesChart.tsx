import React from 'react';
import { EventStream } from './EventStream';
import { Timeline } from './Timeline';
import { ZillaTemporalData } from './types';

interface ZillaTimeSeriesChartProps {
    temporalData: ZillaTemporalData;
    onEventSelect: (event: any) => void;
}

export const ZillaTimeSeriesChart: React.FC<ZillaTimeSeriesChartProps> = ({ 
    temporalData, 
    onEventSelect 
}) => {
    const renderRiskTimeline = () => {
        return temporalData.timeline.map((event, index) => (
            <div
                key={index}
                className={`timeline-event ${event.type} ${event.riskLevel}`}
                style={{
                    left: `${calculateTimelinePosition(event.timestamp)}%`,
                    backgroundColor: getEventColor(event.type, event.riskLevel)
                }}
                onClick={() => onEventSelect(event)}
            >
                <div className="event-marker">
                    {event.type === 'vpn_detection' && '🔒'}
                    {event.type === 'account_creation' && '👤'}
                    {event.type === 'suspicious_activity' && '🚨'}
                    {event.type === 'pattern_change' && '🔄'}
                </div>
                <div className="event-tooltip">
                    <strong>{event.type.replace('_', ' ').toUpperCase()}</strong>
                    <br />
                    {new Date(event.timestamp).toLocaleString()}
                    <br />
                    Risk: {Math.round(event.riskLevel * 100)}%
                </div>
            </div>
        ));
    };

    const calculateTimelinePosition = (timestamp: string) => {
        // Calculate position based on timestamp
        const startTime = new Date(temporalData.timeline[0].timestamp).getTime();
        const endTime = new Date(temporalData.timeline[temporalData.timeline.length - 1].timestamp).getTime();
        const eventTime = new Date(timestamp).getTime();
        
        return ((eventTime - startTime) / (endTime - startTime)) * 100;
    };

    const getEventColor = (type: string, riskLevel: number) => {
        if (riskLevel > 0.7) return '#ff4444';
        if (riskLevel > 0.4) return '#ffaa00';
        return '#44ff44';
    };

    return (
        <div className="zilla-time-series">
            <div className="timeline-container">
                <h3>Digital Footprint Timeline</h3>
                <div className="timeline-track">
                    {renderRiskTimeline()}
                </div>
            </div>

            <EventStream events={temporalData.activityStream} />
            <Timeline data={temporalData} />
            
            <div className="pattern-changes">
                <h4>Behavioral Pattern Changes</h4>
                {temporalData.patternChanges.map((change, index) => (
                    <div key={index} className="pattern-change">
                        <span className="change-type">{change.type}</span>
                        <span className="change-magnitude">
                            {change.magnitude > 0 ? '+' : ''}{Math.round(change.magnitude * 100)}%
                        </span>
                        <span className="change-time">
                            {new Date(change.timestamp).toLocaleDateString()}
                        </span>
                    </div>
                ))}
            </div>
        </div>
    );
};
