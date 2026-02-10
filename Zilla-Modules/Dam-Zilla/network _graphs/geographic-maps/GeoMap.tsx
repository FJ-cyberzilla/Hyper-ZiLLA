import React from 'react';
import { HeatMap } from './HeatMap';
import { LocationTracker } from './LocationTracker';
import { ZillaLocationData } from './types';

interface ZillaGeoMapProps {
    locations: ZillaLocationData[];
    onLocationSelect: (location: ZillaLocationData) => void;
}

export const ZillaGeoMap: React.FC<ZillaGeoMapProps> = ({ locations, onLocationSelect }) => {
    const getLocationColor = (type: string) => {
        switch (type) {
            case 'expected':
                return '#00ff00'; // Green - expected location
            case 'vpn_location':
                return '#ff0000'; // Red - VPN location
            case 'estimated_real':
                return '#ffff00'; // Yellow - estimated real location
            default:
                return '#666666';
        }
    };

    const renderLocationMarkers = () => {
        return locations.map((location, index) => (
            <div
                key={index}
                className="zilla-location-marker"
                style={{
                    position: 'absolute',
                    left: `${location.coordinates.x}%`,
                    top: `${location.coordinates.y}%`,
                    width: location.radius / 2,
                    height: location.radius / 2,
                    backgroundColor: getLocationColor(location.type),
                    borderRadius: '50%',
                    cursor: 'pointer',
                    boxShadow: `0 0 10px ${getLocationColor(location.type)}`,
                    border: '2px solid white',
                    opacity: location.confidence
                }}
                onClick={() => onLocationSelect(location)}
                title={`
                    ${location.type.toUpperCase()}
                    ${location.city}, ${location.country}
                    Confidence: ${Math.round(location.confidence * 100)}%
                    ${location.provider ? `Provider: ${location.provider}` : ''}
                `}
            >
                {location.type === 'expected' && '🏠'}
                {location.type === 'vpn_location' && '🔒'}
                {location.type === 'estimated_real' && '🎯'}
            </div>
        ));
    };

    return (
        <div className="zilla-geo-map">
            <div className="map-container" style={{ position: 'relative', width: '100%', height: '400px' }}>
                {/* Base world map would go here */}
                {renderLocationMarkers()}
                <HeatMap locations={locations} />
                <LocationTracker locations={locations} />
            </div>
            
            <div className="location-legend">
                <div className="legend-item">
                    <span className="color expected">■</span>
                    Expected Location (Carrier)
                </div>
                <div className="legend-item">
                    <span className="color vpn">■</span>
                    VPN/Proxy Location
                </div>
                <div className="legend-item">
                    <span className="color estimated">■</span>
                    Estimated Real Location
                </div>
            </div>
        </div>
    );
};
