// Получаем элементы формы
const form = document.getElementById('register-form');
const registerButton = document.getElementById('register-btn');

// Добавляем обработчик события на кнопку регистрации
registerButton.addEventListener('click', (e) => {
  e.preventDefault();
  const nickname = document.getElementById('nickname').value;
  const email = document.getElementById('email').value;
  const password = document.getElementById('password').value;

  const userData = {
    nickname,
    email,
    password
  };

  // Отправляем запрос на сервер
  fetch("https://localhost:443/users", {
    method: 'POST',
    //mode: 'no-cors',
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(userData)
  })
  .then((response) =>
 {
        if (response.ok){
          console.log('ok')
          //window.location.href= "../testsite/index.html"
          location.reload();
        }
        else{
          console.log('error')
        }
        return response.json()
  })
  .then((data) => {
    console.log(data)
  })
/*  .then((data) => {
    document.getElementById('register-response').innerHTML = data;
      data.innerHTML = 'Пользователь зарегистрирован успешно!';
      window.location.href = "../main/main.html"
      data.innerHTML = 'Ошибка регистрации: ' + data.error;*/
  .catch(error => {
    console.error(error);
  })
});
