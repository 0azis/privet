import csv

def city_exist(city: str) -> bool:
  user_city = ''.join(city.split()).upper().replace("-", "")
  with open('./towns.csv') as f:
    arr = list(csv.reader(f))

  low = 0 # minimal index
  high = len(arr) - 1 # max index

  while low < high:
    middle = (low+high) // 2
    suspect = ''.join(arr[middle][0].split()).upper().replace("-", "")

    if suspect == user_city:
      return True

    elif suspect > user_city:
      high = middle - 1

    elif suspect < user_city:
      low = middle + 1 

  return False
