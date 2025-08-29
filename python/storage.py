# python/storage.py
import logging
import time
import json
from uuid import uuid4
import pymysql
from pymysql.cursors import DictCursor
import redis

from models import User, Order

logger = logging.getLogger(__name__)

class MySQLClient:
    def __init__(self, config):
        self.config = config
        self.db = None
        try:
            self.db = pymysql.connect(
                host=self.config.host,
                port=self.config.port,
                user=self.config.username,
                password=self.config.password,
                database=self.config.database,
                cursorclass=DictCursor,
                connect_timeout=self.config.conn_timeout.total_seconds(),
                read_timeout=self.config.read_timeout.total_seconds(),
                write_timeout=self.config.write_timeout.total_seconds()
            )
            self.init_schema()
            logger.info("Successfully connected to MySQL and initialized schema.")
        except Exception as e:
            logger.warning(f"Failed to connect to MySQL: {e}. Using mock MySQL client.")
            self.db = None # Use mock behavior if connection fails

    def init_schema(self):
        if not self.db:
            return

        with self.db.cursor() as cursor:
            # Create users table
            create_user_table_sql = """
            CREATE TABLE IF NOT EXISTS users (
                id VARCHAR(36) NOT NULL PRIMARY KEY,
                name VARCHAR(255) NOT NULL,
                email VARCHAR(255) NOT NULL UNIQUE
            );
            """
            cursor.execute(create_user_table_sql)
            logger.info("Users table checked/created successfully.")

            # Create orders table
            create_order_table_sql = """
            CREATE TABLE IF NOT EXISTS orders (
                id VARCHAR(36) NOT NULL PRIMARY KEY,
                user_id VARCHAR(36) NOT NULL,
                item VARCHAR(255) NOT NULL,
                amount DECIMAL(10, 2) NOT NULL,
                FOREIGN KEY (user_id) REFERENCES users(id)
            );
            """
            cursor.execute(create_order_table_sql)
            logger.info("Orders table checked/created successfully.")
        self.db.commit()

    def query_row(self, query, *args):
        if not self.db:
            logger.warning("MySQL not implemented for query_row. Returning None.")
            return None
        with self.db.cursor() as cursor:
            cursor.execute(query, args)
            return cursor.fetchone()

    def query(self, query, *args):
        if not self.db:
            logger.warning("MySQL not implemented for query. Returning empty list.")
            return []
        with self.db.cursor() as cursor:
            cursor.execute(query, args)
            return cursor.fetchall()

    def exec_query(self, query, *args):
        if not self.db:
            logger.warning("MySQL not implemented for exec_query. Returning None.")
            return None
        with self.db.cursor() as cursor:
            result = cursor.execute(query, args)
            self.db.commit()
            return result

class RedisClient:
    def __init__(self, config):
        self.config = config
        self.client = None
        try:
            self.client = redis.Redis(
                host=self.config.host,
                port=self.config.port,
                password=self.config.password,
                db=self.config.database,
            )
            self.client.ping()
            logger.info("Successfully connected to Redis.")
        except Exception as e:
            logger.warning(f"Failed to connect to Redis: {e}. Redis client will be None.")
            self.client = None

    def get(self, key):
        if not self.client:
            return None
        return self.client.get(key)

    def set(self, key, value, expiration=None):
        if not self.client:
            return
        self.client.set(key, value, ex=expiration)

    def set_user(self, user: User, expiration=None):
        key = f"user:{user.id}"
        user_data = json.dumps(user.__dict__)
        self.set(key, user_data, expiration)

    def get_user(self, user_id: str):
        user_data = self.get(f"user:{user_id}")
        if user_data:
            return User(**json.loads(user_data))
        return None

    def set_order(self, order: Order, expiration=None):
        key = f"order:{order.id}"
        order_data = json.dumps(order.__dict__)
        self.set(key, order_data, expiration)

    def get_order(self, order_id: str):
        order_data = self.get(f"order:{order_id}")
        if order_data:
            return Order(**json.loads(order_data))
        return None

    def set_user_ids(self, user_ids: list[str], expiration=None):
        self.set("all_user_ids", json.dumps(user_ids), expiration)

    def get_user_ids(self):
        user_ids_data = self.get("all_user_ids")
        if user_ids_data:
            return json.loads(user_ids_data)
        return None

    def start_fault(self, delay: int):
        if not self.client:
            logger.warning("Redis client is not available. Cannot start fault.")
            return
        try:
            self.client.execute_command("FAULT.START", delay)
            logger.info(f"Redis latency fault started with delay: {delay}ms")
        except Exception as e:
            logger.error(f"Failed to send fault command to Redis proxy: {e}")

    def stop_fault(self):
        if not self.client:
            logger.warning("Redis client is not available. Cannot stop fault.")
            return
        try:
            self.client.execute_command("FAULT.STOP")
            logger.info("Redis latency fault stopped.")
        except Exception as e:
            logger.error(f"Failed to send stop fault command to Redis proxy: {e}")

class Store:
    def __init__(self, mysql_client: MySQLClient, redis_client: RedisClient):
        self.mysql = mysql_client
        self.redis = redis_client

    def query_users_from_redis(self):
        if not self.redis or not self.redis.client:
            logger.info("Redis client is not available or not connected. Simulating HTTP operation to fetch users.")
            time.sleep(0.01) # Simulate a network delay for HTTP operation
            users = []
            for i in range(10):
                user = User(
                    id=str(uuid4()),
                    name=f"Mock User HTTP {i+1}",
                    email=f"mock_http{i+1}@example.com"
                )
                users.append(user)
            logger.info("Successfully simulated HTTP operation for 10 users.")
            return users, None

        user_ids = self.redis.get_user_ids()
        if user_ids:
            users = []
            for user_id in user_ids:
                user = self.redis.get_user(user_id)
                if user:
                    users.append(user)
                else:
                    logger.warning(f"Failed to get user {user_id} from Redis. Re-fetching from mock data.")
                    users = [] # Clear incomplete list to re-fetch
                    break
            if users:
                logger.info("All users retrieved from Redis cache by individual IDs.")
                return users, None

        logger.info("Users not in Redis, mocking 10 users and caching them.")
        users = []
        new_user_ids = []
        for i in range(10):
            user = User(
                id=str(uuid4()),
                name=f"Mock User {i+1}",
                email=f"mock{i+1}@example.com"
            )
            users.append(user)
            new_user_ids.append(user.id)

            self.redis.set_user(user)

        self.redis.set_user_ids(new_user_ids)

        logger.info("Mocked 10 users and cached in Redis (individual users and IDs).")
        return users, None

    def query_or_create_users_from_mysql(self):
        # Check if MySQL is initialized
        if not self.mysql or not self.mysql.db:
            logger.info("MySQL is not initialized. Returning 10 mock users.")
            users = []
            for i in range(10):
                user = User(
                    id=str(uuid4()),
                    name=f"Mock User DB Not Init {i+1}",
                    email=f"mock_db_not_init{i+1}@example.com"
                )
                users.append(user)
            return users, None

        # MySQL is initialized, query for users
        logger.info("MySQL is initialized. Querying users from database.")
        users = self.mysql.query("SELECT id, name, email FROM users")

        if users:
            logger.info(f"Found {len(users)} users in database.")
            # Convert database results to User objects
            user_objects = [User(id=user['id'], name=user['name'], email=user['email']) for user in users]
            return user_objects, None
        else:
            logger.info("No users found in database. Creating 10 mock users and writing to database.")
            users = []
            for i in range(10):
                user = User(
                    id=str(uuid4()),
                    name=f"Mock User DB {i+1}",
                    email=f"mock_db{i+1}@example.com"
                )
                users.append(user)

                # Insert user into database
                self.mysql.exec_query(
                    "INSERT INTO users (id, name, email) VALUES (%s, %s, %s)",
                    user.id, user.name, user.email
                )

            logger.info("Successfully inserted 10 mock users into database.")
            return users, None
