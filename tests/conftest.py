import pytest
from HyperZilla.COMMAND_CENTER.app import create_app
from HyperZilla.database import db

@pytest.fixture(scope='module')
def test_app():
    """Create and configure a new app instance for each test."""
    app = create_app('testing')
    with app.app_context():
        db.create_all()
        yield app
        db.drop_all()

@pytest.fixture(scope='module')
def test_client(test_app):
    """A test client for the app."""
    return test_app.test_client()
