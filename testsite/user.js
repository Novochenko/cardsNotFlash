// Отправляем запрос на сервер для получения аватарки
fetch('https://localhost:443/private/lkshow/', {
  method: 'GET',
  credentials: 'include',
  headers: {
    'Content-Type': 'application/octet-stream'
  }
})
.then((response) => response.formData())
.then((data) => {
    const img = data.has('image');
    let src = 'https://www.hdsfoods.co.uk/wp-content/uploads/2014/01/Potato-White.jpg';
    console.log(img);

    const user_id = data.get('user_id');
    const email = data.get('email');
    const nickname = data.get('nickname');
    const user_description = data.get('user_desription');

    if(!(data.has('user_desription'))){
      user_description = '';
    }
    const cards_count = data.get('cards_count');
    const userTkn = document.getElementById("tab/two");
    // Обновляем аватарку и имя пользователя
    if (data.has('image')){
      let newSrc = src.replace('www.hdsfoods.co.uk/wp-content/uploads/2014/01/Potato-White.jpg','localhost:10443/images/pfpimages/user_id.png');
      newSrc = newSrc.replace('user_id', user_id);
      src = newSrc;
    }
    const headerTake = document.getElementById('header-grid');
    const headerUser = document.createElement('div');
    headerUser.innerHTML=`
    <ul>
      <li id='header-nickname'>${nickname}</li>
      <li><img id="mini-img" src="${src}" alt="Uploaded Image" width="60"></li>
      <button id='header-exit'>Exit</button>
    </ul>
    `;
    headerTake.appendChild(headerUser);
    console.log(src);
        const exitBut = document.getElementById('header-exit');
        exitBut.addEventListener('click', ()=>{
          Swal.fire({
            title: `Вы уверены что хотите выйти?`,
            icon: 'question',
            showCancelButton: true,
            confirmButtonText: 'Выйти',
            cancelButtonText: 'Отмена'
          })
          .then((result) => {
            if(result.isConfirmed){
              fetch("https://localhost:443/private/sessionquit",{
                method:"GET",
                credentials:"include"
              })
              .then(response =>{
                if(response.ok){
                  window.location.href = "../auth/auth.html";
                }
                else{
                  console.log('failed to delete the session');
                }
              })
            }
            else{
              console.log('Выход отменён');
            }
          })
        })
    const userLk = document.createElement('div');
      userLk.innerHTML=`
      <div class="wrapper">
          <div class="left">
              <img id="uploaded-image" src="${src}" alt="Uploaded Image" width="200">
              <h4>${nickname}</h4>
               <p>Designer</p>
               <input type="file" id="file-input" accept="image/png"/>
               <button id="select-file-btn">Select PNG file</button>
               <button id="edit-description">Изменить описание</button>
          </div>
          <div class="right">
              <div class="info">
                  <h3>Информация</h3>
                  <div class="info_data">
                       <div class="data">
                          <h4>Email</h4>
                          <p>${email}</p>
                       </div>
                  </div>
              </div>
    
            <div class="projects">
                  <h3></h3>
                  <div class="projects_data">
                       <div class="data">
                          <h4>Описание</h4>
                          <p>${user_description}</p>
                       </div>
                       <div class="data">
                         <h4>Количество карт</h4>
                          <p>${cards_count}</p>
                    </div>
                  </div>
              </div>
          </div>
      </div>
      `;
      userTkn.appendChild(userLk);
      const editDescription = document.getElementById("edit-description");
      editDescription.addEventListener('click', () =>{
        Swal.fire({
          title: `Введите описание:`,
          html: `
            <label class="description" for="user-description">Описание:</label>
            <input id="user-description" type="text" class="swal2-input">
          `,
          focusConfirm: false,
          showCancelButton: true,
          cancelButtonText: 'Выход',
          confirmButtonText: 'Принять',
          confirmButtonColor: '#3085d6',
          cancelButtonColor: '#d33',
          allowOutsideClick: false,
        
        preConfirm: () => {
          const user_description = document.getElementById('user-description').value;
            // Send data to server or perform other actions
            console.log(`Описание: ${user_description}`);
            const newDescription = {
              "user_description": user_description
            }
            fetch('https://localhost:443/private/lkdescriptionedit',{
              method: "POST",
              credentials: "include",
              headers:{
                "Content-Type": "aplication/json"
              },
              body:JSON.stringify(newDescription)
            })
            .then(response =>{
              if (response.ok){
                location.reload();
                console.log('Описание добавлено');
                
              }
              else{
                console.log('Ошибка добавления описания');
              }
            })
          }
        })
      })
      const uploadedImage = document.getElementById('uploaded-image');
      console.log(uploadedImage.src);
      const uploadBut = document.getElementById('select-file-btn');
      const imageInput = document.getElementById('file-input');
      if (imageInput){
        uploadBut.addEventListener('click', () =>{
          //e.preventDefault();
          console.log('upload');
          const file = imageInput.files[0];
          const reader = new FileReader();
          reader.readAsArrayBuffer(file);
          reader.onload = () => {
            const arrayBuffer = reader.result;
            const uint8Array = new Uint8Array(arrayBuffer);
            const formData = new FormData();
            formData.append('image', new Blob([uint8Array], { type: 'image/png' }));
            req = new Blob([uint8Array], { type: 'image/png' });

            // const blob = new Blob([uint8Array], {type: 'image/png'});
            // const url = URL.createObjectURL(blob);

            fetch('https://localhost:443/private/pfpupload', {
              method: 'POST',
              credentials: 'include',
              headers: {
                  'Content-Type': 'application/octet-stream'
              },
              body: req
              })
              .then((response) => {
                if (response.ok){
                  console.log('image upload successfully');
                  //location.reload();
                }
                else{
                  console.log('image upload error');
                }
              return response;
              })
              .catch((error) => {
                console.error(error);
              });
          }
        })
      }

});
