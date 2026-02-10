import { ForceGraph } from './network-graphs/ForceGraph';
import { NodeRenderer } from './network-graphs/NodeRenderer';
import { BehavioralChart } from './network-graphs/behavioral_charts/BehaviorChart';
import { GeoMap } from './network-graphs/geographic-maps/GeoMap';
import { TimeSeriesChart } from './network-graphs/temporal-analysis/TimeSeriesChart';

export class ZillaNetworkVisualizer {
    private forceGraph: ForceGraph;
    private geoMap: GeoMap;
    private behavioralChart: BehavioralChart;

    constructor() {
        this.forceGraph = new ForceGraph();
        this.geoMap = new GeoMap();
        this.behavioralChart = new BehavioralChart();
    }

    async visualizeDigitalFootprint(forensicData: any) {
        // Convert ZILLA-DAM intelligence to visualization data
        const networkData = this.convertToNetworkGraph(forensicData);
        const geographicData = this.convertToGeographicData(forensicData);
        const behavioralData = this.convertToBehavioralData(forensicData);
        const temporalData = this.convertToTemporalData(forensicData);

        // Render all visualizations
        await Promise.all([
            this.renderNetworkGraph(networkData),
            this.renderGeographicMap(geographicData),
            this.renderBehavioralAnalysis(behavioralData),
            this.renderTemporalAnalysis(temporalData)
        ]);

        return this.createDashboard();
    }

    private convertToNetworkGraph(forensicData: any) {
        const nodes = [];
        const links = [];

        // Primary target node
        nodes.push({
            id: 'target',
            type: 'primary',
            label: forensicData.phoneNumber,
            carrier: forensicData.carrierInfo.carrier,
            country: forensicData.carrierInfo.country,
            riskLevel: forensicData.threatAssessment.overallRisk,
            size: 20
        });

        // Social media accounts as nodes
        forensicData.socialMedia.platforms.forEach((platform: any) => {
            if (platform.found) {
                nodes.push({
                    id: platform.username,
                    type: 'social_account',
                    platform: platform.platform,
                    username: platform.username,
                    confidence: platform.confidence,
                    size: 10 + (platform.confidence * 10)
                });

                links.push({
                    source: 'target',
                    target: platform.username,
                    type: 'phone_to_account',
                    strength: platform.confidence
                });
            }
        });

        // VPN/Proxy nodes
        forensicData.vpnAnalysis.vpnProviders.forEach((vpn: any) => {
            nodes.push({
                id: vpn.provider,
                type: 'vpn_provider',
                provider: vpn.provider,
                detectionConfidence: vpn.confidence,
                size: 8
            });

            links.push({
                source: 'target',
                target: vpn.provider,
                type: 'vpn_usage',
                strength: vpn.confidence
            });
        });

        // Puppet network connections
        forensicData.puppetNetwork.accounts.forEach((account: any) => {
            nodes.push({
                id: account.username,
                type: 'puppet_account',
                platform: account.platform,
                coordinationLevel: account.coordinationLevel,
                size: 6
            });

            links.push({
                source: 'target',
                target: account.username,
                type: 'puppet_network',
                strength: account.coordinationLevel === 'HIGH' ? 0.9 : 0.6
            });
        });

        return { nodes, links };
    }

    private convertToGeographicData(forensicData: any) {
        const locations = [];

        // Expected location from carrier
        locations.push({
            type: 'expected',
            coordinates: forensicData.carrierInfo.location.coordinates,
            country: forensicData.carrierInfo.country,
            city: forensicData.carrierInfo.location.city,
            confidence: 0.95,
            radius: 50 // km radius
        });

        // Detected VPN locations
        forensicData.vpnAnalysis.ipLocations.forEach((location: any) => {
            locations.push({
                type: 'vpn_location',
                coordinates: location.coordinates,
                country: location.country,
                city: location.city,
                provider: location.provider,
                confidence: location.confidence,
                radius: 100
            });
        });

        // Estimated real location
        if (forensicData.realLocationEstimate) {
            locations.push({
                type: 'estimated_real',
                coordinates: forensicData.realLocationEstimate.coordinates,
                country: forensicData.realLocationEstimate.country,
                city: forensicData.realLocationEstimate.city,
                confidence: forensicData.realLocationEstimate.confidence,
                radius: 25
            });
        }

        return { locations };
    }

    private convertToBehavioralData(forensicData: any) {
        return {
            activityPatterns: forensicData.behavioralAnalysis.activityHeatmap,
            timezoneAnalysis: forensicData.behavioralAnalysis.primaryTimezone,
            anomalyDetection: forensicData.behavioralAnalysis.anomalies,
            riskIndicators: forensicData.threatAssessment.riskIndicators
        };
    }

    private convertToTemporalData(forensicData: any) {
        return {
            timeline: forensicData.temporalAnalysis.events,
            activityStream: forensicData.behavioralAnalysis.activityStream,
            patternChanges: forensicData.behavioralAnalysis.patternChanges
        };
    }
}
