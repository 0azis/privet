import json
import aiogram
import redis
from aiogram.types import Message, KeyboardButton, ReplyKeyboardMarkup, ReplyKeyboardRemove, FSInputFile
from aiogram import F, Bot, Dispatcher, Router
from aiogram.fsm.state import State, StatesGroup
from aiogram.fsm.context import FSMContext
from aiogram.filters import CommandStart
from aiogram.enums.content_type import ContentType
import logging
import asyncio
from client import Client 
from client import User
from cities import city_exist

token = "6724580660:AAHp3VQLT-pTBlgnyi1nzoOf1FvK-gKKcxQ"

bot = aiogram.Bot(token)
dp = aiogram.Dispatcher()
client = Client()
cache = redis.Redis()
router = Router()

def startup_keyboard():
  actions = [
    [KeyboardButton(text="Смотреть анкеты")],
    [KeyboardButton(text="Мой профиль")],
    [KeyboardButton(text="Настройки профиля")]
  ]
  keyboard = ReplyKeyboardMarkup(keyboard=actions, resize_keyboard=True)
  return keyboard

class Form(StatesGroup):
  id = State()
  name = State()
  age = State()
  location = State()
  gender = State()
  avatar = State()
  description = State()

class Activities(StatesGroup):
  activities = State() # should have be a list (always)

class Wishes(StatesGroup):
  age = State()
  gender = State()
  activities = State()

class Reaction(StatesGroup):
  user = State()

@router.message(CommandStart())
async def cmd_start(message: Message):
  intro = FSInputFile("../background.png")
  await bot.send_photo(message.chat.id, intro, caption="👋 <b>Привет! Это бот Privet!</b>\nТвой спутник в мире знакомств с настройкой предпочтений и точной рекомендательной системой.\n\n<i>Бот рассчитан на подростков, поэтому просим больших дядь и теть здесь не задерживаться😉</i>", parse_mode='HTML') 
  # await message.answer("👋 <b>Привет! Это бот Privet!</b>\nОдин из самых удобных ботов для знакомств и общения! Бот рассчитан на подростков, поэтому просим больших дядь и теть здесь не задерживаться😉", parse_mode='HTML')
  user_id = message.from_user.id 
  if client.user_exists(user_id):
    await message.answer("Готов смотреть анкеты?", reply_markup=startup_keyboard())
  else:
    create = [[KeyboardButton(text="Создать анкету")]]
    keyboard = ReplyKeyboardMarkup(keyboard=create, resize_keyboard=True)
    await message.answer("Зарегестрируйся и сможешь смотреть анкеты других людей", reply_markup=keyboard)

@router.message(F.text.lower() == "создать анкету")
async def make_profile(message: Message, state: FSMContext):
  await state.update_data(id=message.from_user.id)
  await state.set_state(Form.name)
  await message.answer("<b>Напишите свое имя</b>\n\n-------\n<i>Просим указывать свое настоящее имя, мы бережно храним ваши данные❤️</i>", parse_mode='HTML', reply_markup=ReplyKeyboardRemove())

@router.message(Form.name)
async def set_name(message: Message, state: FSMContext):
  await state.update_data(name=message.text)
  await state.set_state(Form.age)
  await message.answer("<b>Напишите свой возраст</b>\n\n-------\n<i>Поддерживаем современное поколение от 14 до 30 лет🔥</i>", parse_mode='HTML')

@router.message(Form.age)
async def set_age(message: Message, state: FSMContext):
  try:
    age = int(message.text)
    await state.update_data(age=age)
    await state.set_state(Form.location)
    await message.answer("<b>Напишите, где вы живете</b>\n\n-------\n<i>Учитываем все города России💪</i>", parse_mode='HTML')
  except:
    await message.answer("<b>Возраст должен быть числом в диапазоне от 14 до 30 лет</b>", parse_mode='HTML')

@router.message(Form.location)
async def set_location(message: Message, state: FSMContext):
  if message.text.isdigit():
    await message.answer("")
  elif city_exist(message.text) == False:
    await message.answer("<b>К сожалению, мы не нашли такого города в наших базах😭</b>\n\n<i>Попробуйте указать другое значение</i>", parse_mode='HTML')
  else:
    genders = [
      [KeyboardButton(text="Мужской")],
      [KeyboardButton(text="Женский")]
    ]
    keyboard_gender = ReplyKeyboardMarkup(keyboard=genders, resize_keyboard=True)
    await state.update_data(location=message.text)
    await state.set_state(Form.gender)
    await message.answer("<b>Укажите свой пол</b>\n\n-------\n<i>Обещаем, что не будем обращать внимания на несходство внешности и вашего пола😋</i>", parse_mode='HTML', reply_markup=keyboard_gender)

@router.message(Form.gender)
async def set_gender(message: Message, state: FSMContext): 
  gender = message.text.lower()
  if gender not in ["мужской", "женский"]:
    await message.answer("<b>К сожалению, такого пола не существует</b>\n\n<i>Выберите вариант в меню</i>", parse_mode='HTML')
  else:
    await state.update_data(gender=message.text) 
    await state.set_state(Form.avatar)
    await message.answer("<b>Загрузите фото вашего профиля</b>\n\n-------\n<i>Выберите ваше самое лучшее фото (мы знаем, что у вас все такие😊)</i>", parse_mode='HTML', reply_markup=ReplyKeyboardRemove())

@router.message(Form.avatar)
async def set_avatar(message: Message, state: FSMContext):
  if message.content_type != ContentType.PHOTO:
    await message.answer("<b>Ошибка при добавлении фото</b>\n\n<i>Пожалуйста, убедитесь, что вы отправляете изображение🖼</i>", parse_mode='HTML')
  else:
    await message.bot.download(file=message.photo[-1].file_id, destination=f"./upload/{message.from_user.id}.png")

    await state.set_state(Form.description)
    await message.answer("<b>Расскажите немного о себе</b>\n\n------\n<i>Не стесняйтесь и представьте себя со всех лучших сторон😎</i>", parse_mode='HTML')

@router.message(Form.description)
async def set_description(message: Message, state: FSMContext):
  payload = await state.update_data(description=message.text)
  user = User(id=payload['id'], fullname=payload['name'], age=payload['age'], location=payload['location'], gender=payload['gender'], description=payload['description'])
  if client.create_account(user):
    cache.set(f"{message.from_user.id}", 0)
    await message.answer("Отлично! Вы создали свою анкету. Теперь вы можете смотреть анкеты других людей!", reply_markup=startup_keyboard())
    await state.clear()
  else:
    await message.answer("Ошибка. Попробуйте еще раз")

@router.message(F.text.lower() == "смотреть анкеты")
async def start_profiles(message: Message, state: FSMContext):
  try:
    response = client.get_users(message.from_user.id)
    users = json.loads(response)['data']
    i = int(cache.get(f"{message.from_user.id}"))
    reactions = [
      [KeyboardButton(text="❤️")],
      [KeyboardButton(text="💔")]
    ]
    keyboard_reactions = ReplyKeyboardMarkup(keyboard=reactions, resize_keyboard=True)
    if i == len(users):
      await message.answer("Анкеты по вашим предпочтениям закончились. Я напишу, когда появятся новые.", reply_markup=startup_keyboard())
      await state.clear()
    else:
      if users[i]['id'] == message.from_user.id:
        i += 1 
      user = users[i]
      photo = FSInputFile(f"./upload/{user['id']}.png") 
      await message.answer_photo(photo, caption=f"""⚡️⚡️⚡️\n\n<b>Имя</b>: {user['fullName']}\n<b>Возраст</b>: {user['age']}\n<b>Город</b>: {user['location']}\n<b>Пол</b>: {user['gender']}\n<b>Описание</b>: {user['description']}""", reply_markup=keyboard_reactions, parse_mode='HTML')
      await state.update_data(user=user['id'])
      await state.set_state(Reaction.user)
      cache.set(f"{message.from_user.id}", i+1)
  except:
    await message.answer("Ошибка при получении данных")

@router.message(Reaction.user)
async def set_reaction(message: Message, state: FSMContext):
  data = await state.get_data()
  if message.text == "❤️":
    await bot.send_message(data['user'], text=f"Кажется, кому-то понравилась ваша анкета.\nВот этот человек - @{message.from_user.username}\nПриятного общения!")
  await start_profiles(message, state)

@router.message(F.text.lower() == "настройки профиля")
async def wishes(message: Message):
  pass

@router.message(F.text.lower() == "мой профиль")
async def profile(message: Message):
  profile = client.profile(message.from_user.id)
  profile_json = json.loads(profile)['data']   
  photo = FSInputFile(f"./upload/{profile_json['id']}.png") 
  await message.answer_photo(photo, caption=f"""⚡️⚡️⚡️\n\n<b>Имя</b>: {profile_json['fullName']}\n<b>Возраст</b>: {profile_json['age']}\n<b>Город</b>: {profile_json['location']}\n<b>Пол</b>: {profile_json['gender']}\n<b>Описание</b>: {profile_json['description']}""",  parse_mode='HTML')

@router.message()
async def error_message(message: Message):
  await message.answer("Я тебя не понимаю")

async def main():
  bot = Bot(token=token)
  dp = Dispatcher()
  dp.include_router(router)

  await dp.start_polling(bot)

if __name__ == "__main__":
  asyncio.run(main())

