# python/models.py
from dataclasses import dataclass

@dataclass
class User:
    id: str
    name: str
    email: str

@dataclass
class Order:
    id: str
    user_id: str
    item: str
    amount: float
