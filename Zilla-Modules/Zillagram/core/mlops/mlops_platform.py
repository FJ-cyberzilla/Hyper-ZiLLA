# core/mlops/mlops_platform.py
import asyncio
from typing import Dict, List, Optional
from datetime import datetime
import pandas as pd
import numpy as np
from sklearn.metrics import accuracy_score, precision_score, recall_score
import mlflow
import bentoml

class MLOpsPlatform:
    """
    End-to-end ML model management with experiment tracking and deployment
    """
    
    def __init__(self, config: Dict):
        self.config = config
        self.model_registry = ModelRegistry(config)
        self.experiment_tracker = ExperimentTracker(config)
        self.model_monitor = ModelPerformanceMonitor(config)
        self.deployment_manager = ModelDeploymentManager(config)
        
        # Initialize MLflow
        mlflow.set_tracking_uri(config['mlops']['mlflow_tracking_uri'])
        mlflow.set_experiment(config['mlops']['experiment_name'])
    
    async def manage_ml_lifecycle(self, model_type: str, training_data: Dict) -> Dict:
        """End-to-end ML model lifecycle management"""
        lifecycle_steps = {
            "experiment_tracking": await self.experiment_tracker.start_experiment(model_type),
            "model_training": await self._train_model_with_validation(model_type, training_data),
            "model_evaluation": await self._evaluate_model_performance(model_type),
            "model_registration": await self.model_registry.register_model(model_type),
            "model_deployment": await self._deploy_model_to_production(model_type),
            "performance_monitoring": await self.model_monitor.start_monitoring(model_type)
        }
        
        results = await asyncio.gather(*lifecycle_steps.values())
        
        return {
            "model_id": results[2].get('model_id'),
            "lifecycle_steps": dict(zip(lifecycle_steps.keys(), results)),
            "overall_success": all(results)
        }
    
    async def _train_model_with_validation(self, model_type: str, training_data: Dict) -> Dict:
        """Train model with comprehensive validation"""
        with mlflow.start_run():
            # Log parameters
            mlflow.log_params(training_data.get('parameters', {}))
            
            # Prepare data
            X_train, X_test, y_train, y_test = await self._prepare_training_data(training_data)
            
            # Train model
            model = await self._train_model(model_type, X_train, y_train)
            
            # Evaluate model
            evaluation_metrics = await self._evaluate_model(model, X_test, y_test)
            
            # Log metrics and model
            mlflow.log_metrics(evaluation_metrics)
            mlflow.sklearn.log_model(model, "model")
            
            return {
                "model": model,
                "metrics": evaluation_metrics,
                "model_uri": mlflow.get_artifact_uri("model")
            }
    
    async def _evaluate_model_performance(self, model_type: str) -> Dict:
        """Comprehensive model performance evaluation"""
        evaluation_metrics = {
            "accuracy": accuracy_score,
            "precision": precision_score,
            "recall": recall_score,
            "f1_score": lambda y_true, y_pred: 2 * (precision_score(y_true, y_pred) * recall_score(y_true, y_pred)) / (precision_score(y_true, y_pred) + recall_score(y_true, y_pred))
        }
        
        # Get test data
        test_data = await self._get_test_dataset(model_type)
        
        evaluation_results = {}
        for metric_name, metric_func in evaluation_metrics.items():
            try:
                score = metric_func(test_data['y_true'], test_data['y_pred'])
                evaluation_results[metric_name] = score
            except Exception as e:
                logger.error(f"Metric {metric_name} calculation failed: {e}")
                evaluation_results[metric_name] = None
        
        # Calculate overall model health score
        evaluation_results['health_score'] = self._calculate_model_health(evaluation_results)
        
        return evaluation_results

class ModelPerformanceMonitor:
    """Continuous model performance monitoring and drift detection"""
    
    async def start_monitoring(self, model_id: str):
        """Start continuous model performance monitoring"""
        while True:
            try:
                # Monitor prediction drift
                drift_analysis = await self._analyze_prediction_drift(model_id)
                if drift_analysis['drift_detected']:
                    await self._handle_model_drift(model_id, drift_analysis)
                
                # Monitor data quality
                data_quality = await self._monitor_data_quality(model_id)
                if not data_quality['acceptable']:
                    await self._handle_data_quality_issues(model_id, data_quality)
                
                # Monitor business metrics
                business_impact = await self._monitor_business_impact(model_id)
                if business_impact['negative_trend']:
                    await self._handle_business_impact(model_id, business_impact)
                
                await asyncio.sleep(3600)  # Check every hour
                
            except Exception as e:
                logger.error(f"Model monitoring error for {model_id}: {e}")
                await asyncio.sleep(300)
    
    async def _analyze_prediction_drift(self, model_id: str) -> Dict:
        """Analyze model prediction drift over time"""
        recent_predictions = await self._get_recent_predictions(model_id)
        historical_baseline = await self._get_historical_baseline(model_id)
        
        drift_metrics = {
            "distribution_drift": await self._calculate_distribution_drift(recent_predictions, historical_baseline),
            "concept_drift": await self._detect_concept_drift(recent_predictions),
            "performance_decay": await self._measure_performance_decay(model_id)
        }
        
        return {
            "drift_detected": any(drift_metrics.values()),
            "drift_metrics": drift_metrics,
            "confidence": max(drift_metrics.values()) if drift_metrics else 0.0
        }
