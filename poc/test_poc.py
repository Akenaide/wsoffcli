import pytest

from poc import DeckLog


def test_get_create_response():
    dl = DeckLog()

    response = dl.get_token()
    assert response.ok
    assert "token" in response.json()


def test_search():
    dl = DeckLog()
    dl.get_token()

    response = dl.search()
    assert response.json()
    pytest.fail()
