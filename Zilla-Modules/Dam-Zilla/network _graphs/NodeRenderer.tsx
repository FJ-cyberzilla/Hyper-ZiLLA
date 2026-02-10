import React from 'react';
import { ZillaNode } from './types';

interface ZillaNodeRendererProps {
    node: ZillaNode;
    onNodeClick: (node: ZillaNode) => void;
}

export const ZillaNodeRenderer: React.FC<ZillaNodeRendererProps> = ({ node, onNodeClick }) => {
    const getNodeColor = (type: string, riskLevel?: number) => {
        switch (type) {
            case 'primary':
                return riskLevel > 0.7 ? '#ff4444' : riskLevel > 0.4 ? '#ffaa00' : '#44ff44';
            case 'social_account':
                return '#4285f4';
            case 'vpn_provider':
                return '#ea4335';
            case 'puppet_account':
                return '#fbbc05';
            default:
                return '#666666';
        }
    };

    const getNodeIcon = (type: string) => {
        switch (type) {
            case 'primary':
                return '🎯';
            case 'social_account':
                return '👤';
            case 'vpn_provider':
                return '🔒';
            case 'puppet_account':
                return '🎭';
            default:
                return '●';
        }
    };

    return (
        <div 
            className={`zilla-node ${node.type}`}
            style={{
                backgroundColor: getNodeColor(node.type, node.riskLevel),
                width: node.size * 2,
                height: node.size * 2,
                borderRadius: '50%',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                cursor: 'pointer',
                boxShadow: `0 0 ${node.size}px ${getNodeColor(node.type, node.riskLevel)}`,
                border: '2px solid white',
                fontSize: node.size / 2,
                color: 'white',
                fontWeight: 'bold'
            }}
            onClick={() => onNodeClick(node)}
            title={`
                ${node.label || node.username || node.provider}
                ${node.confidence ? `Confidence: ${Math.round(node.confidence * 100)}%` : ''}
                ${node.riskLevel ? `Risk: ${Math.round(node.riskLevel * 100)}%` : ''}
            `}
        >
            {getNodeIcon(node.type)}
        </div>
    );
};
