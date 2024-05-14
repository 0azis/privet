import json
import requests
from dataclasses import dataclass
import dataclasses
@dataclass
class User:
  id: int
  fullname: str
  age: int
  location: str
  gender: str
  description: str

class Client:
  url = 'http://localhost:3000/v1'  
  users = "/users"
  reactions = "/reactions"

  def user_exists(self, id) -> bool:
    resp = requests.get(url=self.url+self.users+"/exist", params={"id": id})
    print(id)
    print(resp)
    return resp.status_code == 200

  def create_account(self, user: User) -> bool:
    body = json.dumps(dataclasses.asdict(user))
    resp = requests.post(url=self.url+self.users, data=body)
    return resp.status_code == 201

  def get_users(self, id):
    resp = requests.get(url=self.url+self.users, params={"id": id})
    return resp.content

  def profile(self, id):
    resp = requests.get(url=self.url+self.users+"/profile", params={"id": id})
    print(resp)
    return resp.content