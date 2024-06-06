// function binEncode(data) {
//   var binArray = []
//   var datEncode = "";

//   for (i=0; i < data.length; i++) {
//       binArray.push(data[i].charCodeAt(0).toString(2)); 
//   } 
//   for (j=0; j < binArray.length; j++) {
//       var pad = padding_left(binArray[j], '0', 8);
//       datEncode += pad + ' '; 
//   }
//   function padding_left(s, c, n) { if (! s || ! c || s.length >= n) {
//       return s;
//   }
//   var max = (n - s.length)/c.length;
//   for (var i = 0; i < max; i++) {
//       s = c + s; } return s;
//   }
//   console.log(binArray);
// }

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
    const src = 'https://www.hdsfoods.co.uk/wp-content/uploads/2014/01/Potato-White.jpg';
    console.log(img);

    const user_id = data.get('user_id');
    const email = data.get('email');
    const nickname = data.get('nickname');
    const user_description = data.get('user_description');
    const cards_count = data.get('cards_count');
    const userTkn = document.getElementById("content-lk");
    // Обновляем аватарку и имя пользователя
    if (data.has('image')){
      src.replace('www.hdsfoods.co.uk/wp-content/uploads/2014/01/Potato-White.jpg','localhost:443/images/pfpimages/${user_id}.png');
    }
    const userLk = document.createElement('div');
      userLk.innerHTML=`
      <div class="wrapper">
          <div class="left">
              <img id="uploaded-image" src=${src} alt="Uploaded Image" width="200">
              <h4>${nickname}</h4>
               <p>Designer</p>
               <input type="file" id="file-input" accept="image/png"/>
               <button id="select-file-btn">Select PNG file</button>
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
                  <h3>Информация</h3>
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
          }
        })
      }
  })
.catch((error) => {
  console.error(error);
});