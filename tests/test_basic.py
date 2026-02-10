import json

def test_system_status_endpoint(test_client):
    """Test the /api/system-status endpoint."""
    response = test_client.get('/api/system-status')
    assert response.status_code == 200
    data = json.loads(response.data)
    assert data['status'] == 'operational'
    assert data['system_status'] == 'operational'