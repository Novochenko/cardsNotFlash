// Получаем элементы профиля
const profile = document.getElementById('profile');
const avatar = document.getElementById('avatar');
const username = document.getElementById('username');

// Отправляем запрос на сервер для получения аватарки
fetch('http://127.0.0.1:443/private/whoami/', {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json'
  }
})
.then((response) => {
  if (response.ok){
    console.log('user ok');
    // ...
  }
  else{
    console.log('user error');
  }
  return response.json()})
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