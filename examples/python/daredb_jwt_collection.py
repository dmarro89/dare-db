import json
import logging
import os
import uuid

import urllib3
from dotenv import load_dotenv
from requests import delete, get, post
from requests.auth import HTTPBasicAuth
from requests.exceptions import RequestException

urllib3.disable_warnings()
logging.getLogger("urllib3.connectionpool").setLevel(logging.ERROR)
logging.basicConfig(level=logging.DEBUG, format="%(asctime)s - %(name)s - %(levelname)s - %(message)s")
# logging.basicConfig(level=logging.DEBUG, format="%(name)s - %(levelname)s - %(message)s")

logger = logging.getLogger(__name__)

load_dotenv()


def load_credentials():
    """
    Loads base URL, username, and password from environment variables.
    """
    base_url = os.getenv("BASE_URL") or "https://127.0.0.1:2605"
    username = os.getenv("DB_USERNAME") or "admin"
    password = os.getenv("DB_PASSWORD") or "password"
    return base_url, username, password


class DareDBPyClientBase:

    def __init__(self, username: str, password: str, base_url: str):
        self.username = username
        self.password = password
        self.base_url = base_url
        self.jwt_token = None
        self.get_jwt_token()

    def get_jwt_token(self):
        """Retrieves JWT token by performing a login with username and password."""

        self.auth_url = f"{self.base_url}/login"
        logger.info(f"URL to get JWT token: {self.auth_url}")

        if not self.jwt_token:
            try:
                response = post(self.auth_url, auth=HTTPBasicAuth(self.username, self.password), verify=False)
                response.raise_for_status()
            except RequestException as e:
                logger.error(f"Error getting JWT token: {e}")

                raise
            logger.debug(f"response.content: {response.content}")
            token_data = response.json()
            self.jwt_token = token_data.get("token")

        return self.jwt_token

    def build_headers_with_jwt(self, jwt_token=None):
        """Constructs the header dictionary with authorization token."""
        token = jwt_token or self.jwt_token
        headers = {"Authorization": f"{token}", "Content-Type": "application/json"}
        return headers

    def send_post_with_jwt(self, url: str, data: dict = None):
        """Sends a POST request with provided URL, data, and JWT headers."""
        headers = self.build_headers_with_jwt()
        try:
            response = post(url, headers=headers, data=json.dumps(data), verify=False)
            return response
        except RequestException as e:
            logger.error(f"Error sending POST request: {e}")
            raise

    def send_get_with_jwt(self, url: str):
        """Sends a GET request with provided URL, and JWT headers."""
        headers = self.build_headers_with_jwt()
        try:
            response = get(url, headers=headers, verify=False)
            return response
        except RequestException as e:
            logger.error(f"Error sending POST request: {e}")
            raise

    def send_delete_with_jwt(self, url: str):
        """Sends a DELETE request with provided URL, and JWT headers."""
        headers = self.build_headers_with_jwt()
        try:
            response = delete(url, headers=headers, verify=False)
            return response
        except RequestException as e:
            logger.error(f"Error sending POST request: {e}")
            raise

    def log_response(self, response):
        if response.status_code not in [200, 201]:
            if response.status_code == 404:
                logger.error(f"HTTP Code: {response.status_code}; content: {response.content}, url: {response.url}")
                return
            logger.error(f"HTTP Code: {response.status_code}; content: {response.content}")


class DareDBDataSamplerSimple(DareDBPyClientBase):

    def populate_db_with_sample_data(self):
        """Populates database with sample data entries."""

        url = f"{base_url}/set"
        logger.debug(f"Populate DB with sample data via URL: {url}")
        MAX_REQUESTS = 5
        for i in range(MAX_REQUESTS):
            data = {f"key_{i}_{uuid.uuid4()}": f"value_{i}"}
            response = self.send_post_with_jwt(url, data)
            self.log_response(response)


class DareDBManageCollections(DareDBPyClientBase):
    MAX_COLLECTIONS = 5

    def create(self, name: str = "sample"):
        url = f"{base_url}/collections/{name}"
        response = self.send_post_with_jwt(url)
        self.log_response(response)

    def create_multiple(self):

        collection_name = "sample"
        for i in range(self.MAX_COLLECTIONS):
            url = f"{base_url}/collections/{collection_name}_{i}"
            response = self.send_post_with_jwt(url)
            self.log_response(response)

    def delete_multiple(self):
        collection_name = "sample"
        for i in range(self.MAX_COLLECTIONS):
            url = f"{base_url}/collections/{collection_name}_{i}"
            response = self.send_delete_with_jwt(url)
            self.log_response(response)

    def list(self):
        url = f"{base_url}/collections"
        response = self.send_get_with_jwt(url)
        if response.status_code == 200:
            logger.info(f"\n{json.dumps(response.json(), indent=2)}")
        else:
            self.log_response(response)


class DareDBSamplerForCollections(DareDBManageCollections):

    MAX_REQUESTS = 5

    def populate(self, collection_name: str = "sample"):
        """Populates database with sample data entries."""

        self.create(collection_name)
        url = f"{base_url}/collections/{collection_name}/set"
        logger.debug(f"Populate DB with sample data via URL (collection: {collection_name}): {url}")

        for i in range(self.MAX_REQUESTS):
            data = {f"key_{i}_{uuid.uuid4()}": f"value_{i} in collection: {collection_name}"}
            response = self.send_post_with_jwt(url, data)
            self.log_response(response)

    def get_all_items(self, collection_name: str = "sample"):

        url = f"{base_url}/collections/{collection_name}/items"
        response = self.send_get_with_jwt(url)
        if response.status_code == 200:
            logger.info(f"\n{json.dumps(response.json(), indent=2)}")
        else:
            self.log_response(response)


if __name__ == "__main__":

    base_url, username, password = load_credentials()

    sampler = DareDBDataSamplerSimple(username, password, base_url)
    sampler.populate_db_with_sample_data()

    collections = DareDBManageCollections(username, password, base_url)

    collections.create()
    collections.create_multiple()
    collections.list()
    collections.delete_multiple()

    sampler_collections = DareDBSamplerForCollections(username, password, base_url)
    sampler_collections.populate()
    sampler_collections.get_all_items()
