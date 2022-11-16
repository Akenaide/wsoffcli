import requests


class DeckLog:
    def __init__(self) -> None:
        self.session = requests.Session()
        self.url = "https://decklog-en.bushiroad.com/system/app/api/create/"
        self.token = None

    def get_token(self):
        headers = {
            "origin": "https://decklog-en.bushiroad.com",
            "referer": "https://decklog-en.bushiroad.com/create?c=2",
        }

        response = self.session.post(self.url, headers=headers)
        self.token = response.json()
        return response

    def search(self):
        data = {
            "param": {
                "title_number": "",
                "keyword": "",
                "keyword_type": ["name", "text", "no", "feature"],
                "side": "",
                "card_kind": "0",
                "level_s": "",
                "level_e": "",
                "power_s": "",
                "power_e": "",
                "color": "0",
                "soul_s": "",
                "soul_e": "",
                "cost_s": "",
                "cost_e": "",
                "trigger": "",
                "option_counter": False,
                "option_clock": False,
                "deck_param1": "N",
                "deck_param2": "##TL##",
            },
            "page": 1,
        }
        return self.session.post(
            "https://decklog-en.bushiroad.com/system/app/api/search/2", data,
        )
