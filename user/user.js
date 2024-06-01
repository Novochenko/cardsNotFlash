// Отправляем запрос на сервер для получения аватарки
fetch('https://localhost:443/private/lkshow/', {
  method: 'GET',
  credentials: 'include',
  headers: {
    'Content-Type': 'application/json'
  }
})
.then((response) => {
  if (response.ok){
    console.log('user ok');
  }
  else{
    console.log('user error');
  }
return response.json()})
.then((data) => {
    const userTkn = document.getElementById("content-lk");
    // Обновляем аватарку и имя пользователя
    const userLk = document.createElement('div');
      userLk.innerHTML=`
      <div class="wrapper">
          <div class="left">
              <img src="https://www.goodwholefood.com/wp-content/uploads/2016/11/potatoes.jpg" alt="user" width="100">
              <h4>${data.nickname}</h4>
               <p>Designer</p>
          </div>
          <div class="right">
              <div class="info">
                  <h3>Информация</h3>
                  <div class="info_data">
                       <div class="data">
                          <h4>Email</h4>
                          <p>${data.email}</p>
                       </div>
                  </div>
              </div>
    
            <div class="projects">
                  <h3>Информация</h3>
                  <div class="projects_data">
                       <div class="data">
                          <h4>Описание</h4>
                          <p>${data.user_description}</p>
                       </div>
                       <div class="data">
                         <h4>Количество карт</h4>
                          <p>${data.cards_count}</p>
                    </div>
                  </div>
              </div>
          </div>
      </div>
      `;
      userTkn.appendChild(userLk);
  })

.catch((error) => {
  console.error(error);
});