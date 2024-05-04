// Получаем элементы профиля
const profile = document.getElementById('profile');
const avatar = document.getElementById('avatar');
const username = document.getElementById('username');

// Отправляем запрос на сервер для получения аватарки
fetch('/get-avatar.php'/* Вставить нужное название файла */, {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json'
  }
})
.then((response) => response.json())
.then((data) => {
  if (data.success) {
    // Обновляем аватарку и имя пользователя
    avatar.src = data.avatarUrl;
    username.textContent = data.username;
  } else {
    console.error('Ошибка загрузки аватарки');
  }
})
.catch((error) => {
  console.error(error);
});