// Получаем элементы формы
const form = document.getElementById('register-form');
const nicknameInput = document.getElementById("nickname");
const emailInput = document.getElementById("email");
const passwordInput = document.getElementById("password");
const registerButton = document.getElementById('register-btn');
const responseDiv = document.getElementById('register-response');


// Добавляем обработчик события на кнопку регистрации
registerButton.addEventListener('click', (e) => {
  e.preventDefault();
  const nickname = nicknameInput.value;
  const email = emailInput.value;
  const password = passwordInput.value;
  // Отправляем запрос на сервер
  fetch('/users', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      nickname,
      email,
      password
    })
  })
  .then((response) => response.json())
  .then((data) => {
    if (data.success) {
      responseDiv.innerHTML = 'Пользователь зарегистрирован успешно!';
      window.location.href = "../main/main.html"
    } else {
      responseDiv.innerHTML = 'Ошибка регистрации: ' + data.error;
    }
  })
  .catch((error) => {
    console.error(error);
  });
});